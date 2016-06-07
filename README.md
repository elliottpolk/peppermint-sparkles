# confgr
Simple configuration server backed by an embedded key/value [bolt datastore](https://github.com/boltdb/bolt). It is 
intended to be used within a private network.

## usage
```bash
$ confgr
usage: confgr <command> [args]

Available commands:
	server		starts confgr server
	list		retrieves the available app configs
	get		retrieves the available config
	set		adds a new config
	remove		deletes the specified config
```
