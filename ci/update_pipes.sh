#!/bin/bash

fly -t dev set-pipeline -p peppermint-sparkles -c dev.pipeline.yml -n -l concourse-credentials.yml && \
    fly -t dev unpause-pipeline -p peppermint-sparkles
