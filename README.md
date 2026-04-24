# CoreStack

CoreStack — a scalable e-commerce backend API built with Go, featuring authentication, product management, shopping cart, order processing, and Stripe integration.

## Features

- **User Authentication**: Registration, login, email verification, and token refresh
- **Product Management**: Create, read, update, and delete products
- **Product Categories**: Organize products by categories
- **Shopping Cart**: Add, update, and remove items from cart
- **Order Management**: Process and track orders
- **Product Reviews**: Users can review and manage product reviews
- **Payment Processing**: Stripe integration for secure payments with webhook support
- **Admin Features**: Category management and admin-only operations
- **Email Verification**: Verify user email addresses during registration
- **JWT Authentication**: Secure token-based authentication with access and refresh tokens
- **CORS Support**: Cross-origin request support for frontend integration
- **Middleware Security**: Role-based access control (user, verified, admin)

## Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Payment**: Stripe API
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/                  # Command-line application entry points and main server initialization
├── internal/             # Private application code (handlers, models, services, repositories)
│   ├── config/          # Configuration management and environment variable handling
│   ├── database/        # Database connection setup and migration runner
│   ├── handlers/        # HTTP request handlers for all API endpoints
│   ├── middleware/      # Authentication and authorization middleware
│   ├── models/          # Data models and database entity definitions
│   ├── repositories/    # Data access layer for database operations
│   ├── services/        # Business logic and external service integrations
│   └── utils/           # Helper utilities (JWT, encryption, etc.)
├── supabase/            # Database migrations and schema definitions
├── Dockerfile           # Docker image configuration for containerization
├── docker-compose.yml   # Multi-container orchestration setup
├── go.mod              # Go module file with dependencies
└── Makefile            # Build and development task automation
```

## Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/ulvinamazow/CoreStack.git
   cd CoreStack
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```
3. **Install development tools**
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```
4. **Start PostgreSQL**
   ```bash
   # Using Docker (recommended)
   docker run --name corestack -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:15
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

### Docker Deployment

1. **Using Docker Compose (recommended)**
   ```bash
   docker compose up -d
   ```
   This will start both the PostgreSQL database and the API server.

2. **Using Docker only**
   ```bash
   # Build the image
   docker build -t corestack .
   
   # Run the container
   docker run -p 5000:5000 --env-file .env corestack
   ```

## API Endpoints

### Authentication
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `GET /api/verify-email` - Verify email address
- `POST /api/logout` - User logout (requires auth)
- `POST /api/resend-verification` - Resend verification email (requires auth)

### User Profile
- `GET /api/profile` - Get user profile (requires auth)
- `PUT /api/profile` - Update user profile (requires auth)

### Products
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product details
- `POST /api/products` - Create product (requires email verification)
- `PUT /api/products/:id` - Update product (requires email verification)
- `DELETE /api/products/:id` - Delete product (requires email verification)

### Categories
- `GET /api/categories` - List all categories
- `POST /api/categories` - Create category (admin only)

### Shopping Cart
- `GET /api/cart` - Get cart items (requires email verification)
- `POST /api/cart` - Add item to cart (requires email verification)
- `PUT /api/cart/:item_id` - Update cart item (requires email verification)
- `DELETE /api/cart/:item_id` - Remove item from cart (requires email verification)
- `POST /api/cart/checkout` - Checkout and create order (requires email verification)

### Orders
- `GET /api/orders` - Get user orders (requires email verification)
- `GET /api/orders/:id` - Get order details (requires email verification)

### Reviews
- `GET /api/products/:id/reviews` - Get product reviews
- `POST /api/products/:id/reviews` - Create review (requires email verification)
- `PUT /api/reviews/:id` - Update review (requires auth)
- `DELETE /api/reviews/:id` - Delete review (requires auth)

### Webhooks
- `POST /api/webhooks/stripe` - Stripe webhook handler

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | PostgreSQL user | postgres |
| `DB_PASSWORD` | PostgreSQL password | (required) |
| `DB_NAME` | Database name | corestack |
| `DATABASE_URL` | PostgreSQL connection URL | (optional) |
| `JWT_SECRET` | JWT signing secret | secret |
| `JWT_ACCESS_HOURS` | Access token validity hours | 24 |
| `JWT_REMEMBER_DAYS` | Remember me duration days | 30 |
| `JWT_REFRESH_DAYS` | Refresh token validity days | 60 |
| `SMTP_HOST` | SMTP server host | smtp.gmail.com |
| `SMTP_PORT` | SMTP server port | 587 |
| `SMTP_USER` | SMTP username | (required) |
| `SMTP_PASSWORD` | SMTP password | (required) |
| `APP_URL` | Application base URL | http://localhost:5000 |
| `VERIFICATION_TOKEN_EXPIRY_HOURS` | Email verification token expiry | 72 |
| `STRIPE_SECRET_KEY` | Stripe API secret key | (required) |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | (required) |
| `STRIPE_CURRENCY` | Default payment currency | usd |
| `PORT` | Server port | 5000 |
| `ADMIN_EMAIL` | Admin email address | example@gmail.com |

## Building and Running

### Build the application
Run main.go to load environment variables from the .env file. 
```bash
go build -o corestack cmd/server/main.go
```

### Run the application
```bash
./corestack
```

## Testing
```bash
go test ./...
```

## Contributing

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add some amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request
