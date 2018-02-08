#!/usr/bin/env bash

# list rbd pools
ceph osd lspools

# create rbd pool
ceph osd pool create kube 64 64

# remove rbd pool
ceph osd pool delete hongtao hongtao --yes-i-really-really-mean-it

# create rbd image
rbd create mlcloud-mysql --size 1024 --pool kube

# disable rbd feature
rbd feature disable kube/mlcloud-mysql deep-flatten
rbd feature disable kube/mlcloud-mysql fast-diff
rbd feature disable kube/mlcloud-mysql object-map
rbd feature disable kube/mlcloud-mysql exclusive-lock
