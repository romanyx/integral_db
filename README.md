[![Go Report Card](https://goreportcard.com/badge/github.com/romanyx/integral_db)](https://goreportcard.com/report/github.com/romanyx/integral_db)
[![Build Status](https://travis-ci.org/romanyx/integral_db.svg?branch=master)](https://travis-ci.org/romanyx/integral_db)

``` sh
make

curl -X POST http://localhost:31000/set -d '{"key": "key", "value": "value"}'
curl -X GET http://localhost:31000/get -d '{"key": "key"}'

make stop
```
