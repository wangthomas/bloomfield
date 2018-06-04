# bloomfield

A scalable bloom filter service in Golang.

[![Build Status](https://travis-ci.com/wangthomas/bloomfield.svg?branch=master)](https://travis-ci.com/wangthomas/bloomfield)

## Features

1. Designed for microservice systems.
2. Scalable bloom filters.
3. Client-side hashing
4. Grpc interface.


## Start command

```
./bloomfield -config_file=config.toml
```
Defalut configuration would be used without a configuration file.


## Test Client

A test client is inside tools/testclient. It's using the [gobloomfield](https://github.com/wangthomas/gobloomfield) as the go client.
In examples below it assumes bloomfield is running locally on port 8679. 

### Create a bloom filter "testfilter"

```
./testclient --hostname=localhost:8679 --create testfilter
```

### Add keys to "testfilter"

The return is a boolean slice. Each boolean means if the corresponding key was in the filter or not. If it's first time adding this key
the return should be false.

```
./testclient --hostname=localhost:8679 --add testfilter key1 key2 key3
```

### Check keys in "testfilter"

```
./testclient --hostname=localhost:8679 --has testfilter key1 key2 key3 key4 key5
```

### Remove the bloom filter "testfilter"

```
./testclient --hostname=localhost:8679 --drop testfilter
```

## Clients

Go - [gobloomfield](https://github.com/wangthomas/gobloomfield)


## Why client-side hashing

A traditional bloom filter service would accept a key as a string and perform hashing using multiple hash functions.

Bloomfield is using client-side hashing -- hashing(partial) is done on the client side.

### Benefits of client-side hashing 

1. Avoid overloading.

In a microservice system there are multiple bloomfield clients. Client-side hashing distrubutes the hashing workload accross multiple
entities instead of overloading bloomfield server.

2. Reduce network payload.

As a string, a key would be arbitrarily long. Client-side hashing deterministically converts a key into two uint64 numbers. In most cases 
it reduces the netwrok payload dramatically and makes the payload size consistent. 

(This effect could be amplified if any messaging system is used in the middle. For example Kafka stores the messages for backing up. 
Reducing the playload size would reduce the Kafka storage's footprint.)

3. Security. 

The original keys are encrypted inside the network payload.


## References

https://bitwangtuo.com/2018/01/01/a-scalable-bloom-filter
