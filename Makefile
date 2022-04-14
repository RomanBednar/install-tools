IMAGE_TAG ?= install-tools:test
ENGINE ?= podman
SSH_PUB_KEY ?= ~/.ssh/id_rsa.pub
PODMAN_AUTH_FILE ?= ${XDG_RUNTIME_DIR}/containers/auth.json

ENTER_CMD = ${ENGINE} run --rm -it --privileged localhost/${IMAGE_TAG} /bin/bash
DEVEL_CMD = ${ENGINE} run --rm -it -v ./:/app:z --privileged localhost/${IMAGE_TAG} /bin/bash

get-secrets:
	mkdir -p ./secrets || .
	cp ${SSH_PUB_KEY} ./secrets
	cp ${PODMAN_AUTH_FILE} ./secrets

build:
	go build -o "$(abspath ./bin/)/install-tool"

image: get-secrets
	${ENGINE} build . -t ${IMAGE_TAG}

enter:
	${ENTER_CMD} || make image && ${ENTER_CMD}

devel: get-secrets
	${DEVEL_CMD}

clean:
	rm -rf ./secrets