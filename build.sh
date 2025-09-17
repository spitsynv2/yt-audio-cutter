#!/bin/sh
VERSION="latest"

docker build --platform=linux/amd64 --no-cache -f deploy/DockerFile.server -t yt-audio-cutter-server:$VERSION .
