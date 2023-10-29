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
- Rest API Format
- Using [Snowflake ID](https://en.wikipedia.org/wiki/Snowflake_ID) Generator (Epoch + NodeID + Sequence)
  - Epoch is `2023-10-29 00:00:00`

## TODO
- [ ] Add e2e tests
- [ ] Add Frontend Web UI

# How to use it

## Demo

[WIP: demo site](https://s.m0ai.dev)

### Generate Short URL

```shell
curl -X POST http://localhost:8080/shorten\?url\=https://google.com | jq
  
> {
>   "short_url": "http://localhost:8080/AaecfgMo",
>   "url": "https://google.com"
> }
```

### Get Short URL

```shell
curl -X GET https://localhost:8080/{shortId}

> HTTP/1.1 308 Permanent Redirect
> Content-Type: text/html; charset=utf-8
> Location: https://google.com
> Date: Sun, 29 Oct 2023 08:26:53 GMT
```

# How to deploy it (aws only)

```shell
# clone the repo
make deploy

# if you want to deploy to your own domain
# Domain and ACM is required
DOMAIN=yourdomain.com make deploy
```
