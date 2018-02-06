# About
**_peppermint-sparkles_** is a simple CLI client and service backed by a configurable key/value datastore. The current storage options are:

* [bolt datastore](https://github.com/boltdb/bolt)
* [redis](https://redis.io/) 

The purpose of **_peppermint-sparkles_** is to provide a configuration / secrets service for services with encryption in mind. Each _secret_ is stored against a **sha256** key based on an application name and environment. Note, encrytion is **_on_** by default. `-encrypt=false` must be used to _not_ encrypt a value.

Currently, for encrypted data uses [PGP](http://www.pgpi.org/doc/pgpintro/) encryption with a generated, base64 encoded, UUIDv4. The encryption token is **_not_** stored anywhere within the **_peppermint-sparkles_** service or datastore. If the token is not provided by the user, it is generated client-side at the time of encryption and displayed once the configuration has been stored. If the token is lost, it **_can not_** be recovered nor can the data. The tokens must be stored in a safe place and redundancy is recommended to prevent any lost configuration / secret data. It is also _not_ advised to reuse the same token for multiple environments and / or configurations.

**_NOTE_**: By design, there is no option to list out all apps / secrets. A request to the service **_must_** include the app name at a minimum.

---

# Building

The simplest way to build is to use [**_honey-do_**](https://github.com/elliottpolk/honey-do) and **_Docker_**. By default, running **_honey clean build_** will build the binary (for Linux distros).

For manual builds, this has a dependency on the [Go](https://golang.org) toolchain. Ensure this is installed or you have the appropriate _Docker_ image ([golang:latest](https://hub.docker.com/_/golang/)).

```bash
# localhost install of go
$ go build -o $GOPATH/bin/psparkles -ldflags "-X main.version=v2.0.0"

# Docker
$ docker run --rm -it -v $GOPATH:/go -w /go/src/git.platform.io/oa-montreal/peppermint-sparkles golang:latest /bin/bash -c 'go build -o $GOPATH/bin/psparkles -ldflags \"-X main.version=v2.0.0\"'
```

For additional build help and ideas, review the `Honeyfile.yml`
 
---

# Running

**_peppermint-sparkles_** can be run in either _client_ or _server_ mode.

```bash
# client
$ psparkles -h
NAME:
   psparkles - TODO...

USAGE:
   psparkles [global options] command [command options] [arguments...]

VERSION:
   v2.0.0

COMMANDS:
     get, ls, list                  retrieves all or specific secrets
     set, add, create, new, update  adds or updates a secret
     delete, del, rm                deletes the secret for the provided app name and optional environment
     server, serve                  start server
     help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
      
###
# server
$ psparkles serve -h
NAME:
   psparkles server - start server

USAGE:
   psparkles server [command options] [arguments...]

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
$ psparkles serve -dst redis
```

### setting a new secret
There are 3 different ways to add a secret:

* `-s <value>` flag
* `-f <full/path/to/file>` flag
* _"piping"_ the results into the command

```bash
$ cat secret.json | psparkles set --addr http://localhost:8080
INFO[0000] token: NWM5N2E4MGEtYTZmZi00MjhhLWE2OTktNWYwOGIwNGQyOTE2
INFO[0000] secret:
{
   "id": "4dc4ab3350179e0d16c212fcc33240441a4cf7b955aa9c5f2784755bc2ed32f7",
   "app_name": "testing",
   "env": "dev",
   "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSU5QaURGcTJXQ0dKZ2VoVW91OGNCSmhpcEpMS2FJdUkwUHFiUzRBSGtqaVlCaWt3Ry9sYm45NnFECjdJc0V0dUhEWCtDejROUGhHOVhnZCtKTmtzdzc0SDdsbGVydkZlaGxJcG9vSmN1ejV5Sjg2TVM1OW5EUWJWVnQKbWJGaE5wM2d1eWZnKytUbE1DNUhOeWh3WlE3NnRSUEI3VGk0NEFUaktlREhuSUpVZEZUZ29PR2lrK0NBNEdyZwpsdVIxTlZkdnRJUWdkUVNtYUwrdmp3VUs0bk5Fayt6aHNMa0EKPVcyMWIKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}
```

### getting an existing secret

```bash
# encrypted
$ psparkles get -addr http://localhost:8080 -a testing -e dev
INFO[0000]
{
 "id": "4dc4ab3350179e0d16c212fcc33240441a4cf7b955aa9c5f2784755bc2ed32f7",
 "app_name": "testing",
 "env": "dev",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSU5QaURGcTJXQ0dKZ2VoVW91OGNCSmhpcEpMS2FJdUkwUHFiUzRBSGtqaVlCaWt3Ry9sYm45NnFECjdJc0V0dUhEWCtDejROUGhHOVhnZCtKTmtzdzc0SDdsbGVydkZlaGxJcG9vSmN1ejV5Sjg2TVM1OW5EUWJWVnQKbWJGaE5wM2d1eWZnKytUbE1DNUhOeWh3WlE3NnRSUEI3VGk0NEFUaktlREhuSUpVZEZUZ29PR2lrK0NBNEdyZwpsdVIxTlZkdnRJUWdkUVNtYUwrdmp3VUs0bk5Fayt6aHNMa0EKPVcyMWIKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}

# to decrypt
$ psparkles get -addr http://localhost:8080 -a testing -e dev -decrypt -t NWM5N2E4MGEtYTZmZi00MjhhLWE2OTktNWYwOGIwNGQyOTE2
INFO[0000]
{
 "id": "4dc4ab3350179e0d16c212fcc33240441a4cf7b955aa9c5f2784755bc2ed32f7",
 "app_name": "testing",
 "env": "dev",
 "content": "{\"user\": \"some_admin\", \"passwd\": \"some_SUPER.Secret#Value\"}"
}
```

### removing configurations

```bash
# displaying current state just for example
$ psparkles get -addr http://localhost:8080 -a testing 
INFO[0000]
{
 "id": "4dc4ab3350179e0d16c212fcc33240441a4cf7b955aa9c5f2784755bc2ed32f7",
 "app_name": "testing",
 "env": "dev",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSVNDOUhtZjRGTHA5Z3hhOFY3ODFlM3lqaEtuODJ6VlZCR2NUUzRBSGs1QmVDU0dnWEx3VGNkczY4Cnp4MWhwZUU2MitEZTRBamhGOHZnU3VMa0RpdTg0UHpsTVpyV2VEREVuTk96eWFFblVtODQzejlua1JJckd1eUsKdjZ0cXYwSk94dFBnRGVSVFljZGN2bnptR0Rielo0M3c5TDZjNERYalNRYmFIZGRabmpmZ2NlSEhJdURHNEF6ZwpTT1FmQ01oUVlCS3BqSlVYM2YyeEpPNjM0bGNnaW5yaHgzd0EKPVdFai8KLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}

INFO[0000]
{
 "id": "17787a9f81d80650c26f6584232e4faf1ab379b81a8f25f6ff30ef7513030e15",
 "app_name": "testing",
 "env": "test",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSTBMOUtnRjErWTRGZ1BHSTRyNEFKZWF5ZzFZNXpBTFozdzEvUzRBSGtFa2h3NmlzbDVhR3JBY3l6CnBXWEw1K0VxSE9BbzRNamgvTWJnWE9LR2NVK0c0S0hsWUU2ZFBLYUFjOHU5eG9aYVNOWDluOHJJVkdHVFc5ajkKRk9LSUoxbUVEMVRnMmVSQkN5alZoWWxWcGlUK21RUzJDU0JDNEwzakY2b2hHV1pxRkh2Z1BlR1RVK0MvNEtMZwo4T1FwV0FwUithZ2NBRVprT1NDM05aWEM0bVQrK3ZIaDF4b0EKPUhOKzcKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}

# REMOVE
$ psparkles rm -addr http://localhost:8080 -id 4dc4ab3350179e0d16c212fcc33240441a4cf7b955aa9c5f2784755bc2ed32f7

# validate remaining for example
$ psparkles get -addr http://localhost:8080 -a testing
INFO[0000]
{
 "id": "17787a9f81d80650c26f6584232e4faf1ab379b81a8f25f6ff30ef7513030e15",
 "app_name": "testing",
 "env": "test",
 "content": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSTBMOUtnRjErWTRGZ1BHSTRyNEFKZWF5ZzFZNXpBTFozdzEvUzRBSGtFa2h3NmlzbDVhR3JBY3l6CnBXWEw1K0VxSE9BbzRNamgvTWJnWE9LR2NVK0c0S0hsWUU2ZFBLYUFjOHU5eG9aYVNOWDluOHJJVkdHVFc5ajkKRk9LSUoxbUVEMVRnMmVSQkN5alZoWWxWcGlUK21RUzJDU0JDNEwzakY2b2hHV1pxRkh2Z1BlR1RVK0MvNEtMZwo4T1FwV0FwUithZ2NBRVprT1NDM05aWEM0bVQrK3ZIaDF4b0EKPUhOKzcKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ=="
}
```
