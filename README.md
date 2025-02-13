# Banking Service
_A simple RESTful API for managing bank accounts and transactions._

## Limitations
- Uses `float64` for amounts, which can introduce rounding errors.
- No currency handling, assumes all transactions are in the same unit.

## Architecture
This project follows a **hexagonal architecture** to maintain clear separation of concerns:

- **Domain**: Core business logic and entities (`Account`, `Transaction`).
- **Service**: Application logic that orchestrates interactions between domain and adapters.
- **Adapters**:
  - **HTTP**: REST API layer.
  - **Storage**: In-memory repository.
- **Ports**: Defines interfaces to decouple adapters from the core logic.

## Usage
We mention how to quickly get started with running and testing the
server both locally and using `Docker`. For a full overview of all
commands see the `Makefile`.

### Running Locally
To run the service in a development environment

```sh
make run
```

To run tests

```sh
make test
```

### Running with Docker
To build and run the service using Docker.

```sh
make docker-run
```

To stop the Docker container:

```sh
make docker-stop
```

To clean up all Docker resources:

```sh
make docker-clean
```
