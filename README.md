
# Rider Assignment System

The Rider Assignment System is a Go-based microservice that simulates a ride-sharing platform's matching service. It manages riders, drivers, and trips, matching riders with nearby available drivers based on their locations using geohashing techniques.

## Features

- Register riders and drivers
- Update driver status and location
- Request rides and match riders with the nearest available drivers
- Retrieve driver and trip details
- Complete trips and update driver availability
- Caching with Redis for efficient driver lookups

## Prerequisites

- Go 1.17 or higher
- Docker (for containerization)
- PostgreSQL (for database)
- Redis (for caching)

## Project Structure

```
rider-assignment-system/
├── api/
│   ├── handlers.go
│   ├── routes.go
├── cache/
│   ├── redis.go
├── config/
│   ├── config.go
│   ├── config.yaml
├── database/
│   ├── db.go
│   ├── migrations/
│       ├── 000001_create_tables.up.sql
│       ├── 000001_create_tables.down.sql
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
├── geohash/
│   ├── geohash.go
├── matching/
│   ├── matcher.go
├── models/
│   ├── driver.go
│   ├── rider.go
│   ├── trip.go
├── postman_collection.json
├── main.go
├── go.mod
├── go.sum
```

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/rider-assignment-system.git
cd rider-assignment-system
```

### 2. Configure the Application

Edit the `config/config.yaml` file to set up your database and Redis configurations.

### 3. Run Database Migrations

Install `golang-migrate` and run the migrations:

```bash
go get -u -d github.com/golang-migrate/migrate/cmd/migrate
migrate -path=./database/migrations -database "postgres://postgres:postgres@localhost:5432/matcha?sslmode=disable" up
```

### 4. Build and Run the Application

#### Using Docker

Run the application using Docker Compose:

```bash
cd docker
docker-compose up
```

#### Without Docker

Ensure PostgreSQL and Redis are running locally, then:

```bash
go build -o rider-assignment-system main.go
./rider-assignment-system
```

### 5. Import the Postman Collection

1. Open Postman.
2. Import the `postman_collection.json` file.
3. Test the endpoints.

## API Endpoints

### Rider Routes
- `POST /riders`: Register a new rider.

### Driver Routes
- `POST /drivers`: Register a new driver.
- `GET /drivers/{driver_id}`: Get driver details by ID.
- `PUT /drivers/{driver_id}/status`: Update driver's status.
- `PUT /drivers/{driver_id}/location`: Update driver's location.

### Trip Routes
- `POST /trips`: Rider requests a ride.
- `GET /trips/{trip_id}`: Get trip details by ID.
- `PUT /trips/{trip_id}/complete`: Mark a trip as completed.

## Environment Configuration

The configuration file `config/config.yaml` contains the following:

```yaml
db:
  user: postgres
  password: postgres
  dbname: matcha
  sslmode: disable
  host: db
  port: "5432"

redis:
  addr: redis:6379
  password: ""
  db: 0
```

Update the settings according to your environment.

## Running Tests

You can add unit tests for handlers and other functionalities using Go's testing framework. Create test files in the corresponding directories, and run:

```bash
go test ./...
```

## Security Considerations

- Implement authentication and authorization for production use.
- Sanitize user inputs to prevent SQL injection.
- Use HTTPS for secure communications.

