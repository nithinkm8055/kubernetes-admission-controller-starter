.DEFAULT_GOAL := docker-image

IMAGE ?= k8s-namespace-mutator:0.0.1
go-build: 
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o image/ ./cmd/webhook-server

.PHONY: docker-image
docker-image: go-build
	cd image/ && docker build -t $(IMAGE) .
