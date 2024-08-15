# Project Structure Documentation

## Table of Contents
1. [Directory Structure](#directory-structure)
2. [Deployment](#deployment)
3. [API Documentation](#api-documentation)

## Directory Structure

### `/deploy`
Contains all manifest and configuration files related to deployment.
- Dockerfile
- YAML files for image building

### `/docs`
Stores all documentation files, including:
- Swagger YAML for API documentation
- Use case diagrams (if applicable)
- Entity-Relationship Diagrams (ERD)

### `/env`
Houses environment example files for all services and programs.

### `/pkg`
Contains common source code that can be used across multiple domains or contexts:
- Validation helpers
- Common logger implementations
- Shared entities

### `/services`
Holds all source code corresponding to specific services.

### `/cmd`
In Go projects, this directory contains `main.go` files and `main()` functions, which serve as entry points for compilation and execution by the runtime process.

### `/config`
Stores configuration initialization logic and validation.

### `/domain`
Contains domain definitions and interfaces that need to be implemented.

#### `/domain/entity`
Defines all required objects.

#### `/domain/repository`
Specifies interfaces for storing and manipulating data from/to data sources. This layer should focus on data operations without complex logic.

#### `/domain/usecase`
Defines interfaces for business logic. All complex logic should be implemented in this layer or within entity methods.

#### `/domain/service`
Specifies interfaces for 3rd party API calls.

### `/internal`
Houses implementations of interfaces defined in the `/domain` directory.

#### `/internal/delivery`
Contains logic for receiving and responding to requests (often called controllers).
- Prepares entities from requests
- Executes use cases
- Formats responses

Examples:
- HTTP delivery
- WebSocket delivery
- gRPC delivery
- Consumer delivery

#### `/internal/repository`
Implements repository interfaces defined in `/domain/repository`.

Example: For a "create user" interface in the domain, this layer might contain MySQL or MongoDB implementations.

#### `/internal/usecase`
Implements use case interfaces defined in `/domain/usecase`.

#### `/internal/service`
Implements service interfaces for 3rd party API calls defined in `/domain/service`.

## Deployment

To deploy the project using Docker Compose:

```bash
docker compose up -d --build
```

## API Documentation

After deploying, view the Swagger documentation at:

[http://localhost:8088/swagger](http://localhost:8088/swagger)