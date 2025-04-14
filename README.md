# gokvstore

A go library which implements a key value store.

## Table of contents

1. [Introduction](#introduction)
1. [Getting started](#getting-started)
    1. [Prerequisites](#prerequisites)
    1. [Installation](#installation)
1. [How to use](#how-to-use)
1. [License](/LICENSE)

## Introduction

I created this library to learn how generics work in golang. 

This library might be extended in the future. Currently it uses no concurrency for kvstore manipulation. 

## Getting started

### Prerequisites

- Golang installed
- Operatingsystem: Linux or Windows, not tested on mac
- IDE like VS Code, if you want to contribute or change the code

### Installation

The recommended way to use this module is using the go cli.

    go get github.com/IT-Administrators/gokvstore

## How to use

To create a key value store the ```NewKVStore``` function should be used.

```Go
// Create empty key value storage.
var kvs = NewKVStore[string, any]()
```

The key value store satisfies the ```Storer``` interface with the functions:

- Get
- Put
- Update
- Delete

To add a value to the store the ```Put``` function can be used.

```Go
// Insert key value to storage.
// Put(key, value)
kvs.Put(1, "Message1")
```

The key value storage also implements two functions ```Save``` and ```Load``` which enable the user to persist the storage and use it in other sessions or applications. In the background these functions use the ```gob``` module to save the storage as a binary file.

The ```Load``` method can also overwrite keys which where changed after the store was saved.

```Go
// Save store to file.
kvs.Save("store.gob")
```

## License

[MIT](./LICENSE)