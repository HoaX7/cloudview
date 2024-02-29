#!/bin/bash

source .env.prod

if [ -z $PROJECT_ID ]; then
    echo "The variable 'PROJECT_ID' is not set in '.env.deploy' file"
    exit 1
fi

build="docker build -t webapp ."
tag="docker tag webapp us-central1-docker.pkg.dev/${PROJECT_ID}/webapp/app:v1"
push="docker push us-central1-docker.pkg.dev/${PROJECT_ID}/webapp/app:v1"

# If on mac os, need to build docker for linux platform
if [[ "$OSTYPE" =~ ^darwin ]]; then
    build="docker build --platform=linux/amd64 -t webapp ."
    echo "building docker image for linux platform on mac os, this might take a while..."
else
    echo "building default docker image.."
fi

# build image
$build
# adds a line break
echo -e "\n"

# tag image
echo "tagging docker image: ${tag}"
$tag
echo -e "\n"

# push image to gcp registry
echo "pushing image to gcp 'webapp' repository"
$push