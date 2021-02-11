#!/usr/bin/env bash

TAG=$(echo $(date) | shasum | awk '{print $1}')

docker build . -t sc-local:$TAG
kind load docker-image sc-local:$TAG

sed 's/SOMEVERSION/'"$TAG"'/' test-deploy.yml | kubectl apply -f -
