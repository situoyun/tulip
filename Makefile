SHELL := /bin/bash
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
# VERSION=git-$(subst /,-,$(BRANCH))-$(shell git describe --tag --dirty)
VERSION=init
IMAGE_TAG=$(VERSION)
IAMGE_REPO=d.autops.xyz
LOCAL_REPO=d.local.autops.xyz

all:

dist:
	env GOOS=linux GOARCH=amd64 go build -v -o ./tulip

distosx:
	env GOOS=darwin GOARCH=amd64 go build -v -o ./tulip

generate:
	go generate ./core/store/...

push:
	docker build -t ${IAMGE_REPO}/tulip:${IMAGE_TAG} .
	docker push ${IAMGE_REPO}/tulip:${IMAGE_TAG}
	docker rmi ${IAMGE_REPO}/tulip:${IMAGE_TAG}

ship:push

push_aarch64:
	docker build -f Dockerfile.aarch64 -t ${IAMGE_REPO}/aarch64/tulip:${IMAGE_TAG} .
	docker push ${IAMGE_REPO}/aarch64/tulip:${IMAGE_TAG}
	docker rmi ${IAMGE_REPO}/aarch64/tulip:${IMAGE_TAG}

protoc:
	protoc --go_out=plugins=grpc:. library/proto/*.proto
	go install ./library/proto
	goimports -l -w ./library/proto

push_dev:
	docker build  -t ${LOCAL_REPO}/tulip:${IMAGE_TAG} .
	docker push ${LOCAL_REPO}/tulip:${IMAGE_TAG}
	docker rmi ${LOCAL_REPO}/tulip:${IMAGE_TAG}

save:
	docker build -t ${IAMGE_REPO}/tulip:${IMAGE_TAG} .
	docker save -o tulip.tar ${IAMGE_REPO}/tulip:${IMAGE_TAG}
	tar -zcvf tulip.tar.gz tulip.tar
	rm -rf tulip.tar