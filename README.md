# y2j

A short, simple tool to convert YAML to JSON as long as it's representable.

## Install

```sh
go get github.com/echojc/y2j
```

## Usage

If a file is provided, convert that file:

```sh
$ echo 'abc: 1' > /tmp/a
$ y2j /tmp/a
{"abc":1}
```

Otherwise, read from stdin:

```sh
$ y2j <<.
abc:
  - 1
  - bar: 2
    quux: 3
.
{"abc":[1,{"bar":2,"quux":3}]}
```

## Build

This project depends on [go-yaml](https://github.com/go-yaml/yaml).

```sh
$ go get gopkg.in/yaml.v2
$ go build
```
