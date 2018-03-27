#!/bin/bash

###
# any cf commands, such as `cf create-service`, should go here

## FIXME - check to see if `find` returns more than 1 result and fail for review
# cd to source directory
cd /tmp/build && cd $(find . -type d -name "peppermint-sparkles" | head -1)

ls -alrt

#current directory should contain the <env>_manifest.yml file
cf push -f pcf/${ENV}_manifest.yml