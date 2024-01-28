package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	namespaceResource     = metav1.GroupVersionResource{Version: "v1", Resource: "namespaces"}
	metadataLabelToRemove = map[string]string{"labelName": "labelValue"}
)

func applyNamespaceDefaults(req *admission.AdmissionRequest) ([]patchOperation, error) {
	// This handler should only get called on Namespace objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != namespaceResource {
		log.Printf("expect resource to be %s", namespaceResource)
		return nil, nil
	}

	// Parse the Namespace object.
	raw := req.Object.Raw
	namespaceObj := corev1.Namespace{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &namespaceObj); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// Create patch operations to apply sensible defaults, if those options are not set explicitly.
	var patches []patchOperation
	labels := namespaceObj.GetLabels()
	for keyToRemove, valueOfKeyToRemove := range metadataLabelToRemove {
		if _, ok := labels[keyToRemove]; ok {
			if labels[keyToRemove] == valueOfKeyToRemove {
				patches = append(patches, patchOperation{
					Op:   "remove",
					Path: "/metadata/labels/" + keyToRemove,
				})
			}
		}
	}

	return patches, nil
}

func main() {

	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applyNamespaceDefaults))

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
