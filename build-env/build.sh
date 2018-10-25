#!/bin/bash
set -e
docker build -t dbeliakov/gitlab-ci-go-npm-env:latest .
docker push dbeliakov/gitlab-ci-go-npm-env:latest
