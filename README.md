# About
**_confgr_** is a simple configuration service backed by an embedded key/value [bolt datastore](https://github.com/boltdb/bolt) to provide a configuration service with optional encryption for applications. Each configuration is stored against an application name and environment. If no environment is specified, the default is **_default_**.

For encrypted configurations, it uses [PGP](http://www.pgpi.org/doc/pgpintro/) encryption with a base64 encoded UUID. The encryption token is **_not_** stored anywhere within the **_confgr_** application or datastore. The token is generated at the time of encryption (if not provided) and displayed once the configuration has been stored. If the token is lost, it **_can not_** be recovered nor can the data. Please store the tokens is a safe place and ensure redundancy to prevent any lost configuration data. It is also _not_ advised to reuse the same token for multiple environments and / or configurations.

---

# Building

The simplest way to build uses **_make_** and **_Docker_** for building. By default, running **_make_** alone will build the binary and generate the _Docker_ image with a build timestamp at the end. If the build will be the production build to be run, it should be tagged appropriately.

```bash
# assumes the image confgr:v1.0.0-1477364779 exists
> docker tag confgr:v1.0.0-1477364779 confgr:latest
```

To build only the binary, running the **_make build_** results in _build/bin/confgr_ though this will only run on an **alpine linux** environment. To build for a specific OS (**e.g.** macOS), run the following:

```bash
> GOOS=darwin make build
```

---

# Running

**_confgr_** can be run in _server_ or _client_ mode.

## Server

### As a container

```bash
# assumes the image confgr:latest exists
> docker run -d --name confgr -p 9001:8080 confgr server
```

### As a stand-alone binary

```bash
> ./confgr server
```



## Client

If running as a client, it requires a server to exist and the environment variable **_CONFGR_ADDR_** to be set.

```bash
# assumes server container is running on localhost
> export CONFGR_ADDR=localhost:9001
```



### setting a new configuration

```bash
> ./confgr set -app foo -env test -encrypt -config '{"test_key": "test_value"}'
# sample output
token: 			 472501d4-83a2-461e-b89c-886ddff63d30
token as base64: NDcyNTAxZDQtODNhMi00NjFlLWI4OWMtODg2ZGRmZjYzZDMw
stored config:
{"app":"foo","environment":"test","config":"LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSVEyQkFVakgwdTZaZ3lMOTJXK0ZKZ2kyc0RlUlp1a1Y0Y2xUUzRBSGsrYnMyL3cvcHZEODBHb1hECk9WQjVST0ZNVytCRDREUGhXYXZndGVJL2Jxd3E0RGprbjhPWVgvV3JMOVgrY1N6OXBVUW9sZUR4NHpYTythMGEKV1o2bTRIUGhVMGZnaWVUandDOCtJUVd4T3U2OXd4eUZLVnV5NGhJUThGbmg1Q2dBCj00UzlXCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="}
```



### getting an existing configuration

```bash
# encrypted
> ./confgr get -app foo -env test
# sample output
{
   "app": "foo",
   "environment": "test",
   "config": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSVEyQkFVakgwdTZaZ3lMOTJXK0ZKZ2kyc0RlUlp1a1Y0Y2xUUzRBSGsrYnMyL3cvcHZEODBHb1hECk9WQjVST0ZNVytCRDREUGhXYXZndGVJL2Jxd3E0RGprbjhPWVgvV3JMOVgrY1N6OXBVUW9sZUR4NHpYTythMGEKV1o2bTRIUGhVMGZnaWVUandDOCtJUVd4T3U2OXd4eUZLVnV5NGhJUThGbmg1Q2dBCj00UzlXCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="
}

# decrypted
> ./confgr get -app foo -env test -decrypt -token NDcyNTAxZDQtODNhMi00NjFlLWI4OWMtODg2ZGRmZjYzZDMw
# sample output
{
   "app": "foo",
   "environment": "test",
   "config": "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSVEyQkFVakgwdTZaZ3lMOTJXK0ZKZ2kyc0RlUlp1a1Y0Y2xUUzRBSGsrYnMyL3cvcHZEODBHb1hECk9WQjVST0ZNVytCRDREUGhXYXZndGVJL2Jxd3E0RGprbjhPWVgvV3JMOVgrY1N6OXBVUW9sZUR4NHpYTythMGEKV1o2bTRIUGhVMGZnaWVUandDOCtJUVd4T3U2OXd4eUZLVnV5NGhJUThGbmg1Q2dBCj00UzlXCi0tLS0tRU5EIFBHUCBNRVNTQUdFLS0tLS0="
}

decrypted config:
{"test_key": "test_value"}
```



### listing out configurations

```bash
> ./confgr list
# sample output
{
   "foo": [
   	  "dev",
      "test"
   ]
}
```



### removing configurations

```bash
# assuming app 'foo' has configurations 'dev' and 'test'
> ./confgr remove -app foo -env dev
> ./confgr list
{
   "foo": [
      "test"
   ]
}
```

Currently there is no way to remove an entire application without removing each individual environment. The only reason for this limitation, at this time, is to prevent accidental deletion of all environments for a given app.