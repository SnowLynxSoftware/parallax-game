## Parallax Game 🎮

The `parallax-game` project is the Monolith application for the Parallax Game, written in Go--for the 2025 PBBG.com Game Jam. It provides APIs to manage the application and a client for user interactions.

## Quick Start 🚀

### Prerequisites

- Go 1.25 or later
- Docker and Docker Compose
- Make (for using the Makefile commands)

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/snowlynxsoftware/parallax-game.git
   cd parallax-game
   ```

2. Set up environment variables:

   ```bash
   cp .env.example .env
   ```

   Update the `.env` file with your local configuration.

3. Install dependencies:

   ```bash
   make deps
   ```

4. Start the development environment:
   ```bash
   make run
   ```

This will start all services (app, PostgreSQL, Redis) in Docker containers and make them available at:

- **Application**: http://localhost:3000
- **PostgreSQL**: localhost:5432

## Development Commands 🛠️

The project uses a comprehensive Makefile for development tasks. Run `make help` to see all available commands:

### Essential Commands

```bash
make help              # Show all available commands
make build             # Build production binary
make test              # Run all tests
make compose-up        # Start Docker Compose in detached mode
make compose-down      # Stop Docker Compose services
```

### Building & Running

```bash
make build             # Build for production (Linux/amd64)
make run               # Run application locally
make migrate           # Run database migrations
```

### Testing

```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make test-race         # Run tests with race detection
make benchmark         # Run benchmark tests
```

### Code Quality

```bash
make check             # Run all checks (format, vet, lint, test)
make fmt               # Format Go code
make lint              # Run linter
make vet               # Run go vet
```

### Docker Commands

```bash
make docker-build      # Build Docker image for production
make docker-push       # Build and push to registry
make compose-up        # Start services with docker-compose
make compose-down      # Stop docker-compose services
make compose-logs      # Show docker-compose logs
```

### Utilities

```bash
make clean             # Clean build artifacts
make clean-docker      # Clean Docker containers and images
make status            # Show project status
```

## Project Structure 📁

```
├── cmd/                    # CLI commands and handlers
├── config/                 # Application configuration
├── docs/                   # Markdown files with documentation
├── migrations/             # Database migrations
├── infrastructure/         # Configuration for Infrastructure with OpenTofu and Ansible
├── scripts/                # Any scripts that we might need to run as part of CI, etc.
├── server/                 # Main server logic
│   ├── controllers/        # HTTP controllers
│   ├── database/           # Database layer
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models
│   ├── services/           # Business logic
│   └── util/               # Utilities
├── static/                 # Static assets
├── templates/              # HTML templates
├── docker-compose.yml      # Development environment
├── Dockerfile              # Production container
└── Makefile                # Development automation
```

## Production Deployment 🚀

### Building for Production

```bash
make docker-build         # Build production Docker image
make docker-push          # Build and push to registry
```

### Container Registry Authentication

When deploying, you may need to authenticate with the Container Registry:

```bash
docker login ghcr.io -u <YOUR_GITHUB_USERNAME> -p <YOUR_PAT>
```

### Environment Variables

The application requires the following environment variables:

- `CLOUD_ENV`: Environment identifier
- `DEBUG_MODE`: Enable debug logging
- `DB_CONNECTION_STRING`: PostgreSQL connection string
- `AUTH_HASH_PEPPER`: Password hashing pepper
- `JWT_SECRET_KEY`: JWT signing secret
- `SENDGRID_API_KEY`: SendGrid API key for emails
- `CORS_ALLOWED_ORIGIN`: CORS allowed origin
- `COOKIE_DOMAIN`: Cookie domain for sessions

## Additional Resources 📚

- [Installing NGINX](https://ubuntu.com/tutorials/install-and-configure-nginx#2-installing-nginx)
- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)

## Contributing 🤝

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and run tests: `make check`
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Open a Pull Request

## License 📄

This project is free and open source software owned by Snow Lynx Software, LLC. It is licensed under MIT. All Rights Reserved.
