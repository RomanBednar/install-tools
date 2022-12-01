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
	GO111MODULE=on go build -o "$(abspath ./bin/)/install-tool"

image: get-secrets
	${ENGINE} build . -t ${IMAGE_TAG}

enter: get-secrets
	${ENTER_CMD} || echo "You need to build the app image first with 'make image'"

devel: get-secrets
	${DEVEL_CMD}

clean:
	rm -rf ./secrets