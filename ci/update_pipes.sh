#!/bin/bash

fly -t dev set-pipeline -p secrets-dev -c dev.pipeline.yml -n -l concourse-credentials.yml && \
    fly -t dev unpause-pipeline -p secrets-dev
