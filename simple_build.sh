#!/usr/bin/env sh
#
#
docker buildx build --platform linux/arm/v7,linux/amd64,linux/arm64 --progress plain --pull -t dragonheim/gagent --push -f docker/Dockerfile .
