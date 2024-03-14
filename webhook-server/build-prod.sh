#!/bin/bash
docker build -f Dockerfile.multistage -t github-webhook:$1 .
docker tag github-webhook:$1 asia-east1-docker.pkg.dev/github-core/ai/github-webhook:$1
docker push asia-east1-docker.pkg.dev/github-core/ai/github-webhook:$1