#!/usr/bin/env bash

mkdir -p /var/log/containers
mkdir -p /var/log/pods/kube-system_kube-apiserver-kind-control-plane_67c12f125236243ab0e88b8f61755eb7/kube-apiserver
echo "hi" > /var/log/pods/kube-system_kube-apiserver-kind-control-plane_67c12f125236243ab0e88b8f61755eb7/kube-apiserver/0.log
ln -s /var/log/pods/kube-system_kube-apiserver-kind-control-plane_67c12f125236243ab0e88b8f61755eb7/kube-apiserver/0.log /var/log/containers/kube-apiserver-kind-control-plane_kube-system_kube-apiserver-ae02df5195171ea4ce0a44fd6b0db90b179cbbeef98414bc3484b1137a82e902.log
