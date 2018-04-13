#!/bin/bash
#
#any other cf commands such as 'cf create-service' can go here
#

# redis service
cf create-service p-redis ${REDIS_PLAN} ${REDIS_NAME}

#current script path is /source/concourse/shared/scripts/push.sh
#cd back to the root and then into the 'target' folder created in assemble.sh where the build folder and manifest.ymls were copied
cd /tmp/build/*/target/..

#current directory should contain the manifest.yml file
cf push -f ${MANIFEST}
