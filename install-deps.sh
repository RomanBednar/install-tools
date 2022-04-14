#!/bin/bash

OC_CLI_URL=https://mirror.openshift.com/pub/openshift-v4/clients/oc/latest/linux/oc.tar.gz

dnf install -y make golang wget
wget ${OC_CLI_URL} -O /tmp/oc.tar.gz
tar xvzf /tmp/oc.tar.gz --directory /tmp
cp /tmp/oc /bin
