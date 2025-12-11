# Go Web Template

A production-ready Go web service template with best practices and modern technologies.

## 技术栈 (Tech Stack)

- **Gin** - High-performance HTTP web framework
- **GORM** - ORM library for MySQL
- **go-redis/v9** - Redis client
- **Zap** - Structured, leveled logging
- **Lumberjack** - Log rotation with automatic file management
- **Viper** - Configuration management with environment variable support
- **Validator** - Request validation (gin's built-in validator/v10)
- **Wire** - Compile-time dependency injection
- **Swagger** - API documentation (swaggo/swag)

## 项目结构 (Project Structure)

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration structures
│   ├── handler/
│   │   └── user.go           # HTTP handlers
│   ├── middleware/
│   │   ├── cors.go           # CORS middleware
│   │   ├── logger.go         # Logging middleware
│   │   └── recovery.go       # Panic recovery middleware
│   ├── model/
│   │   └── user.go           # Data models
│   ├── repository/
│   │   └── user.go           # Data access layer
│   ├── service/
│   │   └── user.go           # Business logic layer
│   └── wire/
│       ├── wire.go           # Wire dependency injection definitions
│       └── wire_gen.go       # Wire generated code
├── pkg/
│   ├── database/
│   │   └── mysql.go          # MySQL connection
│   ├── redis/
│   │   └── redis.go          # Redis connection
│   ├── logger/
│   │   └── logger.go         # Zap logger wrapper
│   └── response/
│       └── response.go       # Unified response format
├── docs/
│   ├── docs.go               # Swagger documentation
│   ├── swagger.json
│   └── swagger.yaml
├── config.yaml               # Configuration file
├── go.mod
├── go.sum
├── Makefile                  # Build commands
├── Dockerfile                # Docker support
├── docker-compose.yaml       # Docker Compose (MySQL & Redis)
└── README.md                 # Project documentation
```

## 快速开始 (Quick Start)

### 前置要求 (Prerequisites)

- Go 1.21 or higher
- Docker and Docker Compose (optional, for local development)
- MySQL 8.0+ (if not using Docker)
- Redis 7.0+ (if not using Docker)

### 安装依赖 (Install Dependencies)

```bash
# Install Go dependencies
go mod download

# Install development tools
make install-tools
```

### 配置 (Configuration)

Edit `config.yaml` to configure your application:

```yaml
server:
  port: 8080
  mode: debug # debug, release, test

database:
  host: localhost
  port: 3306
  user: root
  password: root
  database: go_web_template
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600 # seconds

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10

logger:
  level: debug # debug, info, warn, error
  format: json # json, console
  output: stdout # stdout, file
  file_path: logs/app.log
  max_size: 100       # 单个日志文件最大尺寸（MB）
  max_backups: 3      # 保留的旧日志文件最大数量
  max_age: 28         # 保留旧日志文件的最大天数
  compress: true      # 是否压缩旧日志文件
```

You can also use environment variables to override configuration:
- `SERVER_PORT`
- `DATABASE_HOST`
- `DATABASE_PASSWORD`
- etc.

### 使用 Docker Compose 运行 (Run with Docker Compose)

The easiest way to get started is using Docker Compose:

```bash
# Start MySQL and Redis
make docker-up

# Wait a few seconds for services to be ready, then start the application
make run
```

### 手动运行 (Manual Run)

If you have MySQL and Redis already running:

```bash
# Generate Wire code
make wire

# Generate Swagger documentation
make swagger

# Build the application
make build

# Run the application
make run
# or
./bin/server
```

### 访问应用 (Access Application)

- **API Base URL**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Swagger UI**: http://localhost:8080/swagger/index.html

## API 文档 (API Documentation)

### Health Check

```bash
GET /health
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  }
}
```

### User Management APIs

#### Create User

```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30
}
```

#### Get User

```bash
GET /api/v1/users/:id
```

#### List Users

```bash
GET /api/v1/users?page=1&page_size=10
```

#### Update User

```bash
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "age": 25
}
```

#### Delete User

```bash
DELETE /api/v1/users/:id
```

### Response Format

All API responses follow this format:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code`: 0 for success, non-zero for errors
- `message`: Human-readable message
- `data`: Response payload (optional)

## 开发 (Development)

### Makefile Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Generate Swagger documentation
make swagger

# Generate Wire dependency injection code
make wire

# Start Docker services (MySQL & Redis)
make docker-up

# Stop Docker services
make docker-down

# Install development tools
make install-tools

# Format code
make fmt

# Tidy dependencies
make tidy
```

### Project Architecture

This project follows a clean architecture pattern:

1. **Handler Layer** (`internal/handler`): HTTP request handling, parameter validation
2. **Service Layer** (`internal/service`): Business logic, orchestration
3. **Repository Layer** (`internal/repository`): Data access, database operations
4. **Model Layer** (`internal/model`): Data structures and domain models
5. **Middleware Layer** (`internal/middleware`): Cross-cutting concerns (logging, CORS, recovery)
6. **Package Layer** (`pkg`): Reusable components (database, redis, logger, response)

### Dependency Injection with Wire

This project uses [Wire](https://github.com/google/wire) for compile-time dependency injection. 

To regenerate Wire code after changes:

```bash
make wire
```

### Adding New Features

1. Define your model in `internal/model`
2. Create repository interface and implementation in `internal/repository`
3. Implement business logic in `internal/service`
4. Create HTTP handlers in `internal/handler`
5. Add routes in `cmd/server/main.go`
6. Update Wire configuration in `internal/wire/wire.go`
7. Regenerate Wire code: `make wire`
8. Add Swagger annotations to handlers
9. Regenerate Swagger docs: `make swagger`

## 测试 (Testing)

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...
```

## 部署 (Deployment)

### Using Docker

```bash
# Build Docker image
docker build -t go-web-template .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_HOST=your-mysql-host \
  -e DATABASE_PASSWORD=your-password \
  -e REDIS_HOST=your-redis-host \
  go-web-template
```

### Binary Deployment

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/main.go

# Copy binary and config.yaml to server
# Run the binary
./bin/server
```

## 特性 (Features)

- ✅ Clean architecture with clear separation of concerns
- ✅ Dependency injection with Wire
- ✅ Structured logging with Zap
- ✅ Log rotation with automatic file management (size-based and time-based rotation, compression, auto-cleanup)
- ✅ Configuration management with Viper
- ✅ Request validation with Gin's validator
- ✅ Redis caching support
- ✅ GORM for database operations
- ✅ Automatic database migrations
- ✅ RESTful API design
- ✅ Swagger API documentation
- ✅ CORS middleware
- ✅ Request logging middleware
- ✅ Panic recovery middleware
- ✅ Unified response format
- ✅ Docker and Docker Compose support
- ✅ Makefile for common tasks
- ✅ Secure password hashing with bcrypt
- ✅ Proper HTTP status codes for errors

## 生产环境建议 (Production Recommendations)

This template provides a solid foundation, but for production use, consider:

- **Error Handling**: Implement custom error types to distinguish between different error scenarios (validation errors, not found errors, internal errors) for more precise HTTP status code mapping
- **Cache Strategy**: For high-traffic applications, consider using cache versioning or cache tags instead of key scanning for better performance
- **Authentication**: Add JWT or OAuth2 authentication middleware
- **Rate Limiting**: Implement API rate limiting to prevent abuse
- **Monitoring**: Add metrics collection (Prometheus) and tracing (OpenTelemetry)
- **Testing**: Add comprehensive unit tests and integration tests
- **Database Migrations**: Use a migration tool like golang-migrate for better version control
- **Graceful Shutdown**: Implement graceful shutdown to handle in-flight requests
- **Environment-specific Configs**: Separate configs for dev, staging, and production
- **API Versioning**: Consider API versioning strategy for future changes

## 许可证 (License)

MIT License

## 贡献 (Contributing)

Contributions are welcome! Please feel free to submit a Pull Request.

## 支持 (Support)

If you have any questions or issues, please open an issue on GitHub.
