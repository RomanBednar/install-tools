include ./config/conf.env
CURRENT_DIR = $(shell pwd)

ENTER_CMD = ${engine} run --platform linux/x86_64 --rm -it --workdir /code --cap-add SYS_ADMIN --privileged -v ${CURRENT_DIR}:/code -v ${homeDir}:/root -v /var/run/docker.sock:/var/run/docker.sock ${imageRepo}/${imageName}:${imageTag} /bin/bash

build:
	GO111MODULE=on go build -o "$(abspath ./bin/)/ocp-install-tool"

image:
	${engine} build -f Dockerfile.backend -t ${imageRepo}/${imageName}:${imageTag} .

enter:
	${ENTER_CMD} || echo "You need to build the app image first with 'make image'"

clean:
	rm -rf ./secrets

show-config:
	echo "Secret file is: ${pullSecretFile}"
	echo "SSH public key file is: ${sshPublicKeyFile}"

start:
	podman compose up --build || docker compose up --build

stop:
	podman compose down || docker compose down

prune-images:
	podman compose down --rmi local || docker compose down --rmi local

