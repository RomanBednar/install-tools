FROM docker:latest

ARG GOLANG_VERSION=1.22.0

COPY artifacts/linux-ca-vcenter/* /usr/local/share/ca-certificates/

#we need the go version installed from apk to bootstrap the custom version built from source
RUN apk update && apk upgrade --available && apk add go gcc bash musl-dev openssl-dev ca-certificates && update-ca-certificates

RUN wget https://dl.google.com/go/go$GOLANG_VERSION.src.tar.gz && tar -C /usr/local -xzf go$GOLANG_VERSION.src.tar.gz

RUN cd /usr/local/go/src && ./make.bash

ENV PATH=$PATH:/usr/local/go/bin

RUN rm go$GOLANG_VERSION.src.tar.gz

#we delete the apk installed version to avoid conflict
RUN apk del go

RUN go version


###
ENV OC_CLI_URL=https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/latest/openshift-client-linux-arm64.tar.gz

RUN wget ${OC_CLI_URL} -O /tmp/oc.tar.gz
RUN tar xvzf /tmp/oc.tar.gz --directory /tmp
RUN cp /tmp/oc /bin

