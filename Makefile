SHELL := /bin/sh

all: build run

build:
	docker build \
		-t integral/db:0.0.1 \
		-f docker/Dockerfile \
		.

run:
	kubectl apply -f kubernetes/services/db.yaml
	kubectl apply -f kubernetes/deployments/db.yaml

stop:
	kubectl delete services db
	kubectl delete deployments db

test:
	go test -v --race ./...

bench:
	go test --bench=. ./...
