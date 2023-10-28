# go-url-short

This is a simple url shortener written in golang.
It uses [pulumi](https://www.pulumi.com/) to deploy to AWS Lambda with API Gateway.

# Pre-requisite

- [docker](https://docs.docker.com/get-docker/)
- [pulumi](https://www.pulumi.com/docs/get-started/install/)
- [go>=1.20](https://golang.org/doc/install)
    - [air](https://github.com/cosmtrek/air)
- aws (only for deploy)

## Feature

- Support for multiple database storage options (in-memory, PostgreSQL).
- Deploy to AWS Lambda using Pulumi
- Integration of custom domains via ACM.
- Encode IDs using a base-62

## TODO

- [ ] Snowflake ID Generator
- [ ] Add tests
- [ ] Support REST API Format
- [ ] Add more database store (redis, mysql, etc)
- [ ] Add Frontend Web UI

# How to use it

## Demo

[WIP: demo site](https://s.m0ai.dev)

### Create Short URL

```shell
curl -X POST \
  http://localhost:8080/ \
  -H 'Content-Type: application/json' \
  -d '{ "url": "https://google.com" }'
```

### Get Short URL

```shell
curl -X GET \
  https://localhost:8080/{shortId} \
  -H 'Content-Type: application/json'
```

# How to deploy it (aws only)

```shell
# clone the repo
make deploy

# if you want to deploy to your own domain
# Domain and ACM is required
DOMAIN=yourdomain.com make deploy
```
