# Kosha Redis Connector

Redis is an open source, in-memory data store used as a database, cache, streaming engine, and message broker.

The connector APIs allow you to perform 'RESTful' operations such as reading, modifying, adding or deleting data (keys, lists, hashes, sets etc) from your redis database. 

The APIs also support Cross-Origin Resource Sharing (CORS).


![Twitter](images/redis.png)

This Connector API exposes REST API endpoints to perform any operations on Redis in a simple, quick and intuitive fashion.

It describes various API operations, related request and response structures, and error codes.

## Build

To build the project binary, run
```
    go build -o main .

```

## Run locally

To run the project, simply provide env variables to supply the API key and Freshdesk domain name.


```bash
go build -o main .
USERNAME=<USERNAME> PASSWORD=<PASSWORD>  REDIS_HOST=<REDIS_HOST> ./main
```

This will start a worker and expose the API on port `8012` on the host machine

Swagger docs is available at `https://localhost:8012/docs`

## Generating Swagger Documentation

To generate `swagger.json` and `swagger.yaml` files based on the API documentation, simple run -

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go --parseDependency --parseInternal
```

To generate OpenAPISpec version 3 from Swagger 2.0 specification, run -

```bash
npm i api-spec-converter
npx api-spec-converter --from=swagger_2 --to=openapi_3 --syntax=json ./docs/swagger.json > openapi.json
```
