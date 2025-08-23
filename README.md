# CarZone - Car Management System

![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

CarZone is a comprehensive car management REST API built with Go and PostgreSQL. It provides a complete system for managing cars and their engines with full CRUD operations, following clean architecture principles.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Usage Examples](#usage-examples)
- [Docker Support](#docker-support)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

### Car Management

- Create, read, update, and delete cars
- Get cars by ID or brand
- Comprehensive car information including engine specifications
- Price management and tracking
- Timestamp tracking for creation and updates

### Engine Management

- Full CRUD operations for engines
- Engine specifications management (displacement, cylinders, range)
- Standalone engine operations
- Integration with car records

### Technical Features

- Clean Architecture implementation
- RESTful API design
- PostgreSQL database with proper relations
- UUID-based entity identification
- Comprehensive input validation
- Error handling and logging
- Environment-based configuration
- Docker containerization support

## 🏗️ Architecture

CarZone follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Handlers     │    │    Services     │    │     Stores      │
│   (HTTP Layer)  │───▶│ (Business Logic)│───▶│ (Data Access)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Models      │    │   Interfaces    │    │    Database     │
│ (Data Structures)│    │  (Contracts)    │    │  (PostgreSQL)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Layer Responsibilities

1. **Handlers**: Handle HTTP requests/responses, routing, and input validation
2. **Services**: Implement business logic and coordinate between layers
3. **Stores**: Handle data persistence and database operations
4. **Models**: Define data structures and validation logic
5. **Interfaces**: Define contracts between layers
6. **Driver**: Manage database connections and configurations

## 📁 Project Structure

```
CarZone/
├── 📄 main.go                    # Application entry point and server setup
├── 📄 go.mod                     # Go module dependencies
├── 📄 go.sum                     # Dependency checksums
├── 📄 docker-compose.yml         # Docker Compose configuration
├── 📄 Dockerfile                 # Docker container configuration
├── 📄 README.md                  # Project documentation
├── 📄 .env                       # Environment variables (not in repo)
│
├── 📁 driver/                    # Database connection management
│   └── 📄 postgres.go           # PostgreSQL driver and connection pool
│
├── 📁 models/                    # Data models and validation
│   ├── 📄 car.go                # Car entity and validation logic
│   └── 📄 engine.go             # Engine entity and validation logic
│
├── 📁 store/                     # Data access layer
│   ├── 📄 interface.go          # Store interfaces/contracts
│   ├── 📄 schema.sql            # Database schema and sample data
│   ├── 📁 car/
│   │   └── 📄 car.go            # Car data access operations
│   └── 📁 engine/
│       └── 📄 engine.go         # Engine data access operations
│
├── 📁 service/                   # Business logic layer
│   ├── 📄 interface.go          # Service interfaces/contracts
│   ├── 📁 car/
│   │   └── 📄 car.go            # Car business logic
│   └── 📁 engine/
│       └── 📄 engine.go         # Engine business logic
│
└── 📁 handler/                   # HTTP handlers layer
    ├── 📁 car/
    │   └── 📄 car.go            # Car HTTP request handlers
    └── 📁 engine/
        └── 📄 engine.go         # Engine HTTP request handlers
```

## 🛠️ Technologies Used

### Backend

- **Go 1.24.3**: Main programming language
- **Gorilla Mux**: HTTP router and URL matcher
- **lib/pq**: PostgreSQL driver for Go
- **Google UUID**: UUID generation and parsing

### Database

- **PostgreSQL**: Primary database for data persistence
- **SQL**: Database schema and query language

### Development Tools

- **godotenv**: Environment variable management
- **Docker**: Containerization platform
- **Docker Compose**: Multi-container application management

### Architecture Patterns

- **Clean Architecture**: Separation of concerns and dependency inversion
- **Repository Pattern**: Data access abstraction
- **Service Layer Pattern**: Business logic encapsulation
- **Dependency Injection**: Loose coupling between components

## 📋 Prerequisites

Before running CarZone, ensure you have the following installed:

- **Go 1.24.3 or higher**
- **PostgreSQL 13+ or Docker**
- **Git** (for cloning the repository)

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/PrateekKumar15/CarZone.git
cd CarZone
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Set Up Environment Variables

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=carzone
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
```

### 4. Database Setup

#### Option A: Using Local PostgreSQL

```sql
-- Connect to PostgreSQL and create database
CREATE DATABASE carzone;
-- Run the schema.sql file to set up tables and sample data
\i store/schema.sql
```

#### Option B: Using Docker

```bash
# Start PostgreSQL with Docker
docker run --name carzone-postgres \
  -e POSTGRES_DB=carzone \
  -e POSTGRES_USER=your_username \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 \
  -d postgres:13

# Copy and execute schema
docker cp store/schema.sql carzone-postgres:/schema.sql
docker exec -it carzone-postgres psql -U your_username -d carzone -f /schema.sql
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## ⚙️ Configuration

CarZone uses environment variables for configuration. Create a `.env` file:

```env
# Database Settings
DB_HOST=localhost          # Database host
DB_PORT=5432              # Database port
DB_USER=postgres          # Database username
DB_PASSWORD=password      # Database password
DB_NAME=carzone          # Database name
DB_SSLMODE=disable       # SSL mode for database connection

# Server Settings
SERVER_PORT=8080         # HTTP server port
```

## 📡 API Endpoints

### Car Endpoints

| Method | Endpoint                                        | Description       | Request Body |
| ------ | ----------------------------------------------- | ----------------- | ------------ |
| `GET`  | `/cars/{id}`                                    | Get car by ID     | None         |
| `GET`  | `/cars/brand?brand={brand}&engine={true/false}` | Get cars by brand | None         |
| `POST` | `/cars`                                         | Create new car    | Car JSON     |

### Engine Endpoints

| Method | Endpoint                       | Description          | Request Body |
| ------ | ------------------------------ | -------------------- | ------------ |
| `GET`  | `/engines/{id}`                | Get engine by ID     | None         |
| `GET`  | `/engines/brand?brand={brand}` | Get engines by brand | None         |
| `POST` | `/engines`                     | Create new engine    | Engine JSON  |

### Request/Response Examples

#### Create Car

```bash
curl -X POST http://localhost:8080/cars \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Model S",
    "year": 2024,
    "brand": "Tesla",
    "fuel_type": "Electric",
    "engine": {
      "displacement": 0,
      "no_of_cylinders": 0,
      "car_range": 400
    },
    "price": 79999.99
  }'
```

#### Create Engine

```bash
curl -X POST http://localhost:8080/engines \
  -H "Content-Type: application/json" \
  -d '{
    "displacement": 2000,
    "no_of_cylinders": 4,
    "car_range": 600
  }'
```

#### Get Car by ID

```bash
curl http://localhost:8080/cars/550e8400-e29b-41d4-a716-446655440000
```

## 🗄️ Database Schema

### Cars Table

```sql
CREATE TABLE car (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year VARCHAR(4) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    fuel_type VARCHAR(50) NOT NULL,
    engine_id UUID NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (engine_id) REFERENCES engine(id)
);
```

### Engines Table

```sql
CREATE TABLE engine (
    id UUID PRIMARY KEY,
    displacement INT NOT NULL,
    no_of_cylinders INT NOT NULL,
    car_range INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Entity Relationships

- **One-to-One**: Each car has exactly one engine
- **Foreign Key**: `car.engine_id` references `engine.id`
- **Cascade Delete**: Deleting an engine removes associated cars

## 🐳 Docker Support

### Using Docker Compose (Recommended)

```yaml
version: "3.8"
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=carzone
      - DB_PASSWORD=password
      - DB_NAME=carzone

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=carzone
      - POSTGRES_USER=carzone
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./store/schema.sql:/docker-entrypoint-initdb.d/schema.sql

volumes:
  postgres_data:
```

### Run with Docker Compose

```bash
docker-compose up --build
```

## 🔧 Development

### Code Organization Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Dependency Injection**: Dependencies are injected, not hardcoded
3. **Interface Segregation**: Small, focused interfaces
4. **Error Handling**: Comprehensive error handling at each layer
5. **Validation**: Input validation at multiple levels

### Adding New Features

1. **Define Models**: Add new entities in `models/`
2. **Create Interfaces**: Define contracts in `store/interface.go` and `service/interface.go`
3. **Implement Store**: Add data access logic in `store/{entity}/`
4. **Implement Service**: Add business logic in `service/{entity}/`
5. **Create Handlers**: Add HTTP handlers in `handler/{entity}/`
6. **Register Routes**: Add routes in `main.go`

### Code Style Guidelines

- Use meaningful variable and function names
- Add comprehensive comments
- Follow Go naming conventions
- Handle errors appropriately
- Use dependency injection
- Keep functions small and focused

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./service/car/
```

### Test Structure

```
tests/
├── unit/           # Unit tests for individual components
├── integration/    # Integration tests for API endpoints
└── fixtures/       # Test data and fixtures
```

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Commit changes**: `git commit -am 'Add some feature'`
4. **Push to branch**: `git push origin feature/your-feature-name`
5. **Submit a Pull Request**

### Development Workflow

1. Ensure all tests pass
2. Follow code style guidelines
3. Add tests for new features
4. Update documentation as needed
5. Create descriptive commit messages

## 🔍 Troubleshooting

### Common Issues

1. **Database Connection Issues**

   - Check environment variables
   - Verify PostgreSQL is running
   - Check firewall settings

2. **Port Already in Use**

   - Change `SERVER_PORT` in `.env`
   - Kill existing processes on port 8080

3. **Module Import Issues**
   - Run `go mod tidy`
   - Check `go.mod` file

## 📈 Performance Considerations

- Database connection pooling
- Prepared statements for queries
- Proper indexing on frequently queried fields
- Transaction management for data consistency
- Error logging and monitoring

## 🔐 Security Features

- Input validation and sanitization
- SQL injection prevention with prepared statements
- Environment variable configuration
- Error message sanitization

## 🚀 Future Enhancements

- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Full-text search capabilities
- [ ] Audit logging
- [ ] Metrics and monitoring
- [ ] Unit and integration tests
- [ ] CI/CD pipeline
- [ ] API versioning
- [ ] Swagger/OpenAPI documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Prateek Kumar** - [PrateekKumar15](https://github.com/PrateekKumar15)

## 🙏 Acknowledgments

- Go community for excellent libraries
- PostgreSQL for robust database capabilities
- Gorilla Mux for HTTP routing
- Docker for containerization support

## 📞 Support

If you have any questions or need help with CarZone, please:

1. Check the [Issues](https://github.com/PrateekKumar15/CarZone/issues) page
2. Create a new issue with detailed description
3. Contact the maintainers

---

**Happy Coding! 🚗💨**
