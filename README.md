``` sh
make

curl -X POST http://localhost:31000/set -d '{"key": "key", "value": "value"}'
curl -X GET http://localhost:31000/get -d '{"key": "key"}'

make stop
```
