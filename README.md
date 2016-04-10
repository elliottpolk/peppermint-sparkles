# confgr
Simple configuration server backed by an embedded key/value [bolt datastore](https://github.com/boltdb/bolt). It is 
intended to be used within a private network.

## usage
```bash
$ confgr
usage: confgr command [arguments]

commands:
    server  starts the confgr server
    list    lists out the available app configs
    get retrieves the config for the provided app name
    set sets the config for the specified app
    remove  removes the config and app for the provided app
```
