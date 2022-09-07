# Backend Rest API template

[![Test application](https://github.com/abdivasiyev/go_project_template/actions/workflows/test.yml/badge.svg?event=watch)](https://github.com/abdivasiyev/go_project_template/actions/workflows/test.yml)

Rest API template
===

Project structure:
===
- cmd - this folder contains main entry point of project
- config - contains configurations and global constants
- migrations - contains database `DDL` commands for using with `go-migrate` cli tool
- nginx - contains nginx configurations
- scripts - contains docker container entrypoint scripts
- templates - contains static html templates
- pkg
  - faker_util - contains fake data generators
  - helpers - contains commonly used functions
  - logger - logging utility
  - mailer - email sender
  - response - returns json response with checking errors in response
  - router - contains gin router initialisation
  - security - contains security functions such as generating hash, jwt token validation
  - sentry - contains `sentry` tracer implementation
  - storage - contains query interfaces to database and database connections
  - translator - translates validation errors to human-readable texts
  - validator - validates json requests
- internal
  - handler - contains versioned http handlers, you can also add another type of handlers there, such as event handlers or callback handlers
  - job - contains periodically runnable tasks in background
  - middleware - contains checking permissions, logging requests
  - models - contains request-response DTO s
  - repository - contains database queries, you can implement interfaces to your specific databases there
  - server - contains http server initialisation
  - services - contains business logic of project and all of your code should be here, services can use each other as dependency
  - types - type aliases for using long types 

> Currently, there are app versioning, authorization, swagger documentation, file uploading and rendering, profiling, user and role handlers.

> All packages joined together with `DI` pattern and used `uber.fx` for this. Every package are testable and mockable for every type of testing.

Used 3rd party libraries:
===
- [Go SqlMock](https://github.com/DATA-DOG/go-sqlmock)
- [Sentry SDK](https://github.com/getsentry/sentry-go)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [Go Validator](https://github.com/go-playground/validator)
- [Redis SDK](https://github.com/go-redis/redis)
- [Go Jwt](https://github.com/golang-jwt/jwt)
- [Google UUID](https://github.com/google/uuid)
- [Jmoiron Sqlx](https://github.com/jmoiron/sqlx)
- [Go `pq` lib](https://github.com/lib/pq)
- [Go `errors`](https://github.com/pkg/errors)
- [Viper](https://github.com/spf13/viper)
- [Uber FX](https://github.com/uber-go/fx)
- [Uber Zap](https://github.com/uber-go/zap)
- [Go crypto](https://cs.opensource.google/go/x/crypto)
- [Lumberjack](https://github.com/natefinch/lumberjack)