# (Golang) Backend service for estate tree mapping and optimal drone survey path planning. 

Submitted by **Rachmat Maulana Saleh**

# Provided Requirements

To run this project you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.21
2. [GNU Make](https://www.gnu.org/software/make/)
3. [oapi-codegen](https://github.com/deepmap/oapi-codegen)

    Install the latest version with:
    ```
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
    ```
4. [mock](https://github.com/uber-go/mock)

    Install the latest version with:
    ```
    go install go.uber.org/mock/mockgen@latest
    ```

5. [Docker](https://docs.docker.com/get-docker/) version 20
   
   We will use this for testing your API.

6. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

7. [Node](https://nodejs.org/en) v20

   We will use this for testing your API

8. [NPM](https://www.npmjs.com/) v10

    We will use this for testing your API.

## Initiate The Project

To start working, execute

```
make init
```

## Running

You should be able to run using the script `run.sh`:

```bash
./run.sh
```

You may see some errors since you have not created the API yet.

However for testing, you can use Docker run the project, run the following command:

```
docke -compose up --build
```

You should be able to access the API at http://localhost:8080

If you change `database.sql` file, you need to reinitate the database by running:

```
docker compose down --volumes
```

## Testing

To run test, run the following command:

```
make test
```


/////////////////////
---

# My Implementation
 
## Architecture Decisions
 
### 1. Code Generation from OpenAPI Spec
API contract is defined first in `api.yml`. Server interfaces and types are generated with `oapi-codegen`, keeping the implementation in sync with the spec.
 
### 2. Repository Pattern
All database interactions sit behind `EstateRepository` interface, making handlers fully testable without a real database.
 
### 3. Service Layer for Business Logic
Pure business logic lives in `service/` — drone flight path calculation and tree statistics — independently testable with no HTTP or database dependency.
 
### 4. Mock Generation
Mocks are generated from the repository interface using `mockgen`, keeping unit tests fast and isolated.
 
### 5. OpenAPI Middleware for Validation
Request validation (required fields, value ranges) is handled by `oapi-codegen/echo-middleware` before requests reach the handler.
 
### 6. Boustrophedon Drone Path
The drone surveys every plot in a snake pattern — odd rows left to right, even rows right to left. Height adjusts to `treeHeight + 1` at each plot. When `max-distance` is set, the drone lands at the nearest plot before exceeding the limit.
 
### 7. Multi-stage Docker Build
Builder stage compiles the binary with `golang:1.25-alpine` — bumped from 1.21 to support
Air hot-reload during local development (`air` requires Go 1.22+). Production stage runs
with `alpine:latest` for a minimal image size with no Go toolchain included.
 
---
 
## Modules
 
| Module | Purpose |
|---|---|
| `github.com/labstack/echo/v4` | HTTP framework |
| `github.com/oapi-codegen/echo-middleware` | OpenAPI request validation |
| `github.com/golang/mock` | Mock generation for unit tests |
| `github.com/google/uuid` | UUID generation for IDs |
| `github.com/lib/pq` | PostgreSQL driver |
| `github.com/stretchr/testify` | Test assertions |


## Go Version

The project was developed using **Go 1.21** as specified in `go.mod`, however the
Dockerfile uses `golang:1.25-alpine` to enable hot-reload support via
[Air](https://github.com/air-verse/air) during local development.

Air requires Go 1.22+ to install properly. Bumping the Docker base image to 1.25
was the minimal change needed to get `go install github.com/air-verse/air@latest`
working inside the container without touching the module itself.

The production binary is still compiled and compatible with Go 1.21.
