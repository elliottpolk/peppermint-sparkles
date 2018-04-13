#!/bin/bash

fly -t manulife-ci set-pipeline -p peppermint-sparkles -c pipeline.yml -n -l concourse-credentials.yml && \
    fly -t manulife-ci unpause-pipeline -p peppermint-sparkles
