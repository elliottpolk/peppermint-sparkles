# Peppermint Sparkles

## About

**_Peppermint Sparkles_** (_**sparkles** for short_) is a simple, zero-knowledge service and CLI client backed by a configurable key/value datastore. The current storage options are:

* [bolt datastore](https://github.com/boltdb/bolt)
* [redis](https://redis.io/)

The purpose of **_Peppermint Sparkles_** is to provide a configuration / secrets service, _for_ services, with encryption in mind. Each _secret_ is stored against a **sha256** key based on an application name and environment. Note, encrytion is **_on_** by default. `-encrypt=false` must be used if there is a desire to _not_ encrypt a value.

Currently, encrypted data uses [PGP](http://www.pgpi.org/doc/pgpintro/) encryption. If a token/passphrase is not provided by the user, a base64 encoded UUIDv4 is generated client-side at the time of encryption and displayed **_only_** once the configuration has been successfully stored. To encrypt with said token, we use an AES256 cypher with 15360-bit keys. The speed of encryption will vary depending on the content being encrypted.

The encryption token itself is **_not_** stored anywhere within the **_Peppermint Sparkles_** client, service, or datastore. If the token is lost, it **_can not_** be recovered nor can the data encrypted with said token. The tokens must be stored in a safe, secure place and redundancy is recommended to prevent any lost configuration / secret data. It is also _not_ advised to reuse the same token for multiple environments and / or configurations.

This is a fork and extension of the open-source project [confgr](https://github.com/elliottpolk/confgr). The original project was created under the MIT lincense and this repo _should_ continue that as a result.

**_NOTE_**: By design, there is no option to list out all apps / secrets. A request to the service **_must_** include the secret ID, the app name, and environment.
---

## Peppermint Sparkles Helper

A Docker image has been created for helping with integration patterns. Usage and examples are available for the following platforms:

* [Concourse](ci/README.md)

---

## Builds
The service currently uses the **(optional)** tool ```make``` for builds. The _golang Docker image_ has this pre-installed. To build:

```bash
# ensure the below commands will be run from the project dir
$ cd $GOPATH/src/git.platform.manulife.io/oa-montreal/peppermint-sparkles

# the default targets are 'clean build'
$ make

# to build and package for deployment to PCF or Artifactory
$ make clean build package

# using the Go Docker image
$ docker run --rm -it -v "${GOPATH}":/go -w /go/src/git.platform.manulife.io/oa-montreal/${PWD##*/} golang:latest /bin/sh -c 'make'

# to build and package for deployment to PCF or Artifactory
$ docker run --rm -it -v "${GOPATH}":/go -w /go/src/git.platform.manulife.io/oa-montreal/${PWD##*/} golang:latest /bin/sh -c 'make clean build package'

```

**NOTE:** the default build is for **_Linux_** environments. The ```GOOS=<environment>``` must be specified for OS-specific builds.

```bash

# macOS
$ GOOS=darwin make
# OR
$ docker run --rm -it -v "${GOPATH}":/go -w /go/src/git.platform.manulife.io/oa-montreal/${PWD##*/} golang:latest /bin/sh -c 'GOOS=darwin make'

# Windows
$ GOOS=windows make
# OR
$ docker run --rm -it -v "${GOPATH}":/go -w /go/src/git.platform.manulife.io/oa-montreal/${PWD##*/} golang:latest /bin/sh -c 'GOOS=windows make'
```

---

### Testing

There are 2 types of tests, unit and integration. The **unit tests** can be run on any platform, but the **integration tests** require a machine with **_Docker_** installed as it spins up containers for some dependencies.

```bash

# unit tests
$ make clean unit-tests

# integration tests
# this make take some time to run as there are intentional delays to allow for containers to boot, though longer than 10 mins would be too long a run time
$ make clean integration-tests

# all
$ make clean all-tests
```

---

## Running

**_Peppermint Sparkles_** can be run in either _client_ or _server_ mode.

```bash
# client
$ sparkles -h
NAME:
   sparkles - Server and client for managing super special secrets 🦄

USAGE:
   sparkles [global options] command [command options] [arguments...]

VERSION:
   v3.2.0

COMMANDS:
     get, ls, list                  retrieves secrets
     set, add, create, new, update  adds or updates a secret
     delete, del, rm                deletes a secret
     server, serve                  start the server
     help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   Copyright © 2018
      
###
# server
$ sparkles serve -h
NAME:
   sparkles server - start the server

USAGE:
   sparkles server [command options] [arguments...]

OPTIONS:
   --port value, -p value               HTTP port to listen on (default: "8080") [$PSPARKLES_HTTP_PORT]
   --tls-port value                     HTTPS port to listen on (default: "8443") [$PSPARKLES_HTTPS_PORT]
   --tls-cert value                     TLS certificate file for HTTPS [$PSPARKLES_TLS_CERT]
   --tls-key value                      TLS key file for HTTPS [$PSPARKLES_TLS_KEY]
   --datastore-addr value, --dsa value  address for the remote datastore (default: "localhost:6379") [$PSPARKLES_DS_ADDR]
   --datastore-file value, --dsf value  name / location of file for storing secrets (default: "/var/lib/peppermint-sparkles/psparkles.db") [$PSPARKLES_DS_FILE]
   --datastore-type value, --dst value  backend type to be used for storage (default: "file") [$PSPARKLES_DS_TYPE]
   --help, -h                           show help (default: false)

# assumes a redis instance is running on localhost:6379
$ sparkles serve -dst redis
```

### setting a new secret
There are 3 different ways to add a secret:

* `-s <value>` flag
* `-f <full/path/to/file>` flag
* _"piping"_ the results into the command (currently broken on **_non-macOS_** systems)

```bash
$ cat secret.json | sparkles set --addr http://localhost:8080
INFO[0000] token: OTUzMmE1N2QtZjU5MS00N2Y2LWIxZmEtMzBlYzllZjNlYzNj
INFO[0000] secret:
{
   "id": "50711b9b-4fb3-4192-affe-73c735174ad8",
   "app_name": "testing",
   "env": "dev",
   "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSUVNU1ZrSjBqMlZ4Z3NINzI0U01ZekE4OUdLbGVUMDMzaGZyUzRBSGtQNkZhcjJYbWEvbnYzWnlNCkVKbmNyT0drcitBcTRPZmhxbC9nYytJYlJRV3k0S2JsOEhSRjVSdUhyb1prN0dPMlcvcTJ4U3FELzNEZWxLZ0wKeEJ6V1hDWjVKSWpnU2VUQTcwNEE3eTNFbVhrWXNLWXlhUUJDNEtMajhCekZMN1Y1a2NIZ24rRTdEdUNnNEhuZwpST1NLVFBIU3NiUXpYeWRYeUxwWU9vWFc0cG0wM1IzaE1UWUEKPUZLVUwKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}
```

### getting an existing secret

```bash
# encrypted
$ sparkles get -addr http://localhost:8080 -a testing -e dev --id 50711b9b-4fb3-4192-affe-73c735174ad8
INFO[0000]
{
 "id": "50711b9b-4fb3-4192-affe-73c735174ad8",
 "app_name": "testing",
 "env": "dev",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSUVNU1ZrSjBqMlZ4Z3NINzI0U01ZekE4OUdLbGVUMDMzaGZyUzRBSGtQNkZhcjJYbWEvbnYzWnlNCkVKbmNyT0drcitBcTRPZmhxbC9nYytJYlJRV3k0S2JsOEhSRjVSdUhyb1prN0dPMlcvcTJ4U3FELzNEZWxLZ0wKeEJ6V1hDWjVKSWpnU2VUQTcwNEE3eTNFbVhrWXNLWXlhUUJDNEtMajhCekZMN1Y1a2NIZ24rRTdEdUNnNEhuZwpST1NLVFBIU3NiUXpYeWRYeUxwWU9vWFc0cG0wM1IzaE1UWUEKPUZLVUwKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}

# to decrypt
$ sparkles get -addr http://localhost:8080 -a testing -e dev --id 50711b9b-4fb3-4192-affe-73c735174ad8 --decrypt -t OTUzMmE1N2QtZjU5MS00N2Y2LWIxZmEtMzBlYzllZjNlYzNj
INFO[0000]
{
 "id": "50711b9b-4fb3-4192-affe-73c735174ad8",
 "app_name": "testing",
 "env": "dev",
 "content": "{\"user\": \"some_admin\", \"passwd\": \"some_SUPER.Secret#Value\"}"
}
```

### removing configurations

```bash
# displaying current state just for example
$ sparkles get -addr http://localhost:8080 -a testing -e dev --id 50711b9b-4fb3-4192-affe-73c735174ad8
INFO[0000]
{
 "id": "50711b9b-4fb3-4192-affe-73c735174ad8",
 "app_name": "testing",
 "env": "dev",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSUVNU1ZrSjBqMlZ4Z3NINzI0U01ZekE4OUdLbGVUMDMzaGZyUzRBSGtQNkZhcjJYbWEvbnYzWnlNCkVKbmNyT0drcitBcTRPZmhxbC9nYytJYlJRV3k0S2JsOEhSRjVSdUhyb1prN0dPMlcvcTJ4U3FELzNEZWxLZ0wKeEJ6V1hDWjVKSWpnU2VUQTcwNEE3eTNFbVhrWXNLWXlhUUJDNEtMajhCekZMN1Y1a2NIZ24rRTdEdUNnNEhuZwpST1NLVFBIU3NiUXpYeWRYeUxwWU9vWFc0cG0wM1IzaE1UWUEKPUZLVUwKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}

# REMOVE
$ sparkles rm -addr http://localhost:8080 --id 50711b9b-4fb3-4192-affe-73c735174ad8

# validate for example
$ sparkles get -addr http://localhost:8080 -a testing -e dev --id 50711b9b-4fb3-4192-affe-73c735174ad8
$

```

---

## TODO

- [ ] Audit Tool
    - [ ] CLI
    - [ ] WebUI
- [ ] Hardware key integration
- [ ] `fly` / _Concourse_ integration
