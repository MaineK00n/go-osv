# Test Script For go-osv
Documentation on testing for developers

## Getting Started
```terminal
$ pip install -r requirements.txt
```

## Run test
Use `127.0.0.1:1325` and `127.0.0.1:1326` to diff the server mode between the latest tag and your working branch.

If you have prepared the two addresses yourself, you can use the following Python script.
```terminal
$ python diff_server_mode.py --help
usage: diff_server_mode.py [-h] [--debug | --no-debug]
                           {id,package} {all,crates.io,dwf,go,linux,oss-fuzz,pypi}

positional arguments:
  {id,package}          Specify the mode to test.
  {all,crates.io,dwf,go,linux,oss-fuzz,pypi}
                        Specify the Fetch Type to be started in server mode when testing.

optional arguments:
  -h, --help            show this help message and exit
  --debug, --no-debug   print debug message
```

[GNUmakefile](../GNUmakefile) has some tasks for testing.  
Please run it in the top directory of the go-osv repository.

- build-integration: create the go-osv binaries needed for testing
- clean-integration: delete the go-osv process, binary, and docker container used in the test
- fetch-rdb: fetch data for RDB for testing
- fetch-redis: fetch data for Redis for testing
- diff-cveid: Run tests for ID in server mode
- diff-package: Run tests for Package in server mode
- diff-server-rdb: take the result difference of server mode using RDB
- diff-server-redis: take the result difference of server mode using Redis
- diff-server-rdb-redis: take the difference in server mode results between RDB and Redis

## About the ID and Packages list used for testing
Duplicates are removed from the latest fetched data and prepared.  
For example, for sqlite3, you can get it as follows.  

```terminal
$ sqlite3 go-osv.sqlite3
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
# ID
sqlite> .output integration/id/all.txt
sqlite> SELECT DISTINCT alias FROM osv_aliases;

# Packages
sqlite> .output integration/package/all.txt
sqlite> SELECT DISTINCT name FROM osv_packages;
```
