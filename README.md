## Golang monorepo with REST, gRPC, CQRS, opentelemetry, gorillaMUX, swag, mailhog and rabbitMQ 
#### based on docker containers, docker compose, Makefile 

[![Go Report Card](https://goreportcard.com/badge/github.com/RafalSalwa/auth-api)](https://goreportcard.com/report/github.com/RafalSalwa/auth-api)
[![Run Gosec](https://github.com/RafalSalwa/interview-srv-go/actions/workflows/gosec.yml/badge.svg)](https://github.com/RafalSalwa/interview-srv-go/actions/workflows/gosec.yml)

[![codecov](https://codecov.io/gh/RafalSalwa/interview-srv-go/graph/badge.svg?token=T0DZIOYDR8)](https://codecov.io/gh/RafalSalwa/interview-srv-go)

Codacy:
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/c0054f5a3f1343029e2a3acb76931ebc)](https://app.codacy.com/gh/RafalSalwa/auth-api/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

Code climate:
[![Maintainability](https://api.codeclimate.com/v1/badges/a2df28f0afa241c0d07b/maintainability)](https://codeclimate.com/github/RafalSalwa/auth-api/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/a2df28f0afa241c0d07b/test_coverage)](https://codeclimate.com/github/RafalSalwa/auth-api/test_coverage)


SonarQube:
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=RafalSalwa_auth-api&metric=coverage)](https://sonarcloud.io/summary/new_code?id=RafalSalwa_auth-api)


## Topics Covered
- REST, gRPC, CQRS, Docker, distributed tracing (jaeger, otel, prometheus,grafana, NewRelics), testify
- Onion, Clean architecture, Monorepo,
- Swagger, Postman docs ([/docs](docs) directory)
- JWT , auth and logging middleware
- Env with Viper
- hot reload with cosmtrek/air
- golangcilint, pre-commit-hooks
- MySQL, database/sql, gorm
- Redis
- MongoDB
- Jenkins, GitHub actions
  Plans:
- .gitlabci, buildspec
- mockery for testify, more tests


## Services:
![arch](docs/go_arch.png)
- ### API Gateway
  - Backend For Frontend approach, gateway takes HTTP requests and decide which service over gRPC should be called
    - security & Auth methods 
      - Basic Auth
      - Bearer Token
      - API key
      - JWT
    - HTTP Handlers based on GorillaMux
    - Router with middlewares
      - content_type
      - correlation_id
      - CORS
      - JWT Token decode
      - request_log
    - CQRS
      - separated commands and queries that connect to specific services
        - example commands: signup, change_password
        - example queries: sign_in, user_details, get_verification_code
- ### Auth service
  - gRPC server for commands:
    - SignIn
    - SignUp
    - Verify  account code
  - Service layer to manage data flow
    - rpc -> service -> repositories -> clients
  - Redis cache to prevent eventual consistency
  - MongoDB for users & logs storage
  - MySQL as main RDBMS for users
  - RabbitMQ direct queue with dead letter exchange for email sendout
- ### User service
  - gRPC server for commands 
    - get user(id)
    - get user details
    - change password
  - Mysql repository based on gORM
  - Redis cache for cacheable users
- ### Consumer service
  - AMQP consumer that read events from rabbitMQ, send emails and store logs in mongoDB

- ### Tester service
  - workers pool service that constantly creates users from registration, activation and signIn for JWT Tokens
  - optional daisy chain pattern for concurrency control
  

## Build
make build

make up

## Run tests
make test_unit

make test_integration

### gRPC
- make proto

## Credentials
#### services
- interview@interview.com:VeryG00dPass!


screenshots:
![db](docs/db_design.png)
![jaeger](docs/jaeger.png)
![postman](docs/postman.png)