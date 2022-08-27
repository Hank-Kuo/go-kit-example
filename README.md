# go-kit-example
## Full list what has been used 
* Go
* Jaeger
* Prometheus
* Grafana
* Zipkin
* Postgres
* Docker and docker-compose


## Rrerequisite
- go & grpc
```bash
brew install go
brew install protobuf
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
```
- go vendor
```bash
make tidy
```
## Starting 
run local 
```bash
make local
``` 

run service
```bash
make run APP={xxx}
```

call service
```bash
make call APP={xxx}
```

## API

- product API
request
```json
{
    "user": "hank.kuo",
    "product": "Bike",
    "price": 100,
    "fee": 5,
    "currency": "USD"
}
```
response
```json
{
    "status": "success",
    "message": "user buy a product",
    "data": {
        "cost": 3150
    },
}
```

- price API request
request
```json
{
    "price": 100,
    "fee": 5,
    "currency": "USD"
}
```
response
```json
{
    "status": "success",
    "message": "total cost",
    "data": {
        "cost": 105
    },
}
```

request
```json
{
    "cost": 105,
    "currency": "USD"
}
```

response
```json
{
    "status": "success",
    "message": "total cost",
    "data": {
        "cost": 105
    },
}
```

support currencies: 
- "USD": 30
- "GBP": 35
- "JYP": 0.22
- "TWD": 1


## Jaeger UI:
http://localhost:16686

## Prometheus UI:
http://localhost:9090

## Grafana UI:
http://localhost:3000


## References
- https://github.com/go-kit/examples
- https://github.com/cage1016/ms-demo
- https://github.com/hwholiday/learning_tools/tree/master/hconfig
- https://github.com/sagikazarmark/modern-go-application/blob/main/internal/app/todocli/configure.go