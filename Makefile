include ./config/conf.env
CURRENT_DIR = $(shell pwd)

#ENTER_CMD = ${engine} run --rm -it --privileged localhost/${imageTag} /bin/bash
DEVEL_CMD = ${engine} run --rm -it -v ./:/app:z --privileged localhost/${IMAGE_TAG} /bin/bash
ENTER_CMD = ${engine} run --workdir /code --cap-add SYS_ADMIN --privileged -it -v ${CURRENT_DIR}:/code -v ${homeDir}:/root -v /var/run/docker.sock:/var/run/docker.sock localhost/${imageTag} /bin/bash

get-secrets:
	mkdir -p ./secrets || .
	cp ${sshPublicKeyFile} ./secrets
	jq 'del(.credsStore, .currentContext)' ${pullSecretFile} | tr -d '[:space:]' > ./secrets/config.json

build:
	GO111MODULE=on go build -o "$(abspath ./bin/)/install-tool"

image: #get-secrets
	${engine} build -f Dockerfile.devel -t localhost/${imageTag} .

enter: #get-secrets
	${ENTER_CMD} || echo "You need to build the app image first with 'make image'"

devel: get-secrets
	${DEVEL_CMD}

clean:
	rm -rf ./secrets

show-config:
	echo "Secret file is: ${pullSecretFile}"
	echo "SSH public key file is: ${sshPublicKeyFile}"

