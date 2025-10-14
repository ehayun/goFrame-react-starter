# Tzlev Application Design Document

## Overview

Tzlev is a modern web application that combines a **GoFrame (Go)** backend with an **embedded React** frontend. The React application is built using Vite and integrates with the Wowdash Bootstrap admin template.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Client Browser                     │
│  ┌───────────────────────────────────────────────┐  │
│  │  HTML + CSS (Wowdash Template)                │  │
│  │  ├── Bootstrap 5                              │  │
│  │  ├── jQuery & Libraries                       │  │
│  │  └── React Components (SPA)                   │  │
│  └───────────────────────────────────────────────┘  │
│                        │                             │
│                        │ HTTP/API                    │
│                        ▼                             │
└─────────────────────────────────────────────────────┘
                         │
                         │
┌─────────────────────────────────────────────────────┐
│              GoFrame Server (Port 8080)              │
│  ┌───────────────────────────────────────────────┐  │
│  │  Router + CORS Middleware                     │  │
│  │  ├── Static File Server (public/)            │  │
│  │  ├── API Routes (/api/*)                     │  │
│  │  └── SPA Route (/)                           │  │
│  └───────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────┐  │
│  │  Controllers & Services                       │  │
│  │  ├── Business Logic                          │  │
│  │  ├── Authentication (Gmail OAuth)            │  │
│  │  └── Email Service (SMTP)                    │  │
│  └───────────────────────────────────────────────┘  │
│         │                          │                 │
│         ▼                          ▼                 │
│  ┌──────────────┐          ┌──────────────┐        │
│  │ PostgreSQL   │          │    Redis     │        │
│  │   Database   │          │  Sessions &  │        │
│  │   (MCP)      │          │    Cache     │        │
│  └──────────────┘          └──────────────┘        │
└─────────────────────────────────────────────────────┘
```

## Design Principles

### 1. Embedded Frontend Pattern
- **Decision**: React frontend is built and embedded into the Go binary's serving directory
- **Rationale**:
  - Single deployment artifact
  - Simplified deployment process
  - Better performance (no separate frontend server)
- **Implementation**: Vite builds React to `public/` directory, which GoFrame serves

### 2. CORS Configuration
- **Decision**: Enable CORS middleware for all API endpoints
- **Rationale**:
  - Allows API calls from different origins during development
  - Prepares for future external API consumers
  - Required for future mobile apps or separate frontend deployments
- **Implementation**: GoFrame `CORSDefault()` middleware applied globally

### 3. Asset Management
- **Decision**: Copy template assets from `assets/` to `public/` and serve directly from public
- **Rationale**:
  - Separates source template files from served files
  - Allows customization without modifying original template
  - Clean separation between development and production assets
  - GoFrame serves all files from public directory
- **Implementation**:
  - Template assets stored in `assets/` directory (source)
  - Assets copied to `public/` before running
  - All file references point to root paths (css/, js/, images/ served from public)

### 4. Template Integration
- **Decision**: Use existing Wowdash Bootstrap template with React components
- **Rationale**:
  - Leverage existing CSS/JS libraries (jQuery, Bootstrap)
  - Maintain template functionality (theme switching, sidebar, etc.)
  - React handles dynamic content and SPA routing
- **Implementation**:
  - Template assets in `public/` (css/, js/, images/)
  - Template scripts loaded before React initialization
  - React components mirror template structure
  - `blank-page.html` serves as the starting point for creating the main index file

### 5. Progressive Enhancement
- **Decision**: Load template JavaScript (`app.js`) after React mounts
- **Rationale**:
  - Ensures jQuery and dependencies load first
  - Template features (sidebar toggle, theme switching) work correctly
  - React can interact with template features
- **Implementation**: `useEffect` hook in App.jsx dynamically loads app.js

### 6. Internationalization (i18n)
- **Decision**: Use i18next for multilingual support with Hebrew and English
- **Rationale**:
  - Support Hebrew (RTL) and English (LTR) languages
  - Automatic direction switching based on language
  - Scalable for adding more languages
  - Industry-standard i18n solution for React
- **Implementation**:
  - i18next + react-i18next libraries
  - Translation files in `frontend/src/locales/{lang}/translation.json`
  - Language switcher in Navbar
  - Automatic RTL/LTR direction switching

### 7. Client-Side Routing
- **Decision**: Use React Router v6 for SPA navigation
- **Rationale**:
  - Single Page Application experience
  - Dynamic route-based rendering
  - Browser history management
  - Deep linking support
  - SEO-friendly with proper configuration
- **Implementation**:
  - React Router DOM v6
  - Route definitions in App.jsx
  - Dynamic sidebar menu with React Router Links
  - Breadcrumb navigation integration

### 8. Dynamic Menu Configuration
- **Decision**: Centralized menu configuration with i18n keys
- **Rationale**:
  - Easy to add/remove menu items
  - Automatic translation of menu labels
  - Consistent routing and navigation
  - Support for nested menus
- **Implementation**:
  - Menu configuration file with routes, icons, and translation keys
  - Sidebar component renders from configuration
  - Active route highlighting

### 9. Dual-Mode Operation (Web + CLI)
- **Decision**: Single binary supports both web server and CLI modes
- **Rationale**:
  - Administrative tasks via CLI (migrations, seeding, user management)
  - Web server for production runtime
  - Single binary simplifies deployment
  - No need for separate admin tools or scripts
  - Consistent codebase for both interfaces
- **Implementation**:
  - Check for command-line arguments on startup
  - If arguments present → CLI mode (execute command and exit)
  - If no arguments → Web server mode (start HTTP server)
  - CLI uses Cobra for command-line interface and command parsing
  - Shared access to services and business logic

### 10. External Services Integration
- **Decision**: Integrate with external services for persistence, caching, authentication, and notifications
- **Rationale**:
  - **PostgreSQL**: Reliable ACID-compliant database for data persistence
  - **Redis**: Fast in-memory storage for sessions, caching, and temporary data
  - **Gmail OAuth**: Trusted authentication provider with user data
  - **Email Service**: Essential for user communication and notifications
  - Configuration-driven for easy environment management
- **Services**:
  - **Database**: PostgreSQL 16+ with GoFrame gdb
  - **Cache/Sessions**: Redis 7+ for session storage and caching
  - **Authentication**: Gmail OAuth 2.0 for user login
  - **Email**: SMTP (Gmail) for transactional emails
- **Implementation**:
  - Centralized configuration file (`config/config.yaml`)
  - Environment-based configuration overrides
  - Secrets stored in environment variables
  - Connection pooling and health checks
  - Service initialization on startup

## Configuration Management

### Configuration File Structure

The application uses a YAML-based configuration system with environment-specific overrides:

```
config/
├── config.yaml              # Main configuration file (defaults)
├── config.example.yaml      # Example configuration (for documentation)
├── config.development.yaml  # Development environment overrides
└── config.production.yaml   # Production environment overrides
```

### Configuration File Example (`config/config.yaml`)

```yaml
# Application Configuration
app:
  name: "Tzlev"
  environment: "development"  # development, staging, production
  debug: true
  port: 8080

# Database Configuration (GoFrame format)
database:
  default:
    link: "pgsql:${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?sslmode=disable"
    type: "pgsql"
    debug: true
    maxIdle: 5
    maxOpen: 25
    maxLifetime: "5m"

# Redis Configuration
redis:
  host: "localhost"
  port: 6379
  password: "${REDIS_PASSWORD}"  # Environment variable
  database: 0
  poolSize: 10
  sessionPrefix: "tzlev:session:"
  cachePrefix: "tzlev:cache:"
  sessionTTL: "24h"  # Session expiration time

# Gmail OAuth Configuration
oauth:
  google:
    clientID: "${GOOGLE_CLIENT_ID}"       # Environment variable
    clientSecret: "${GOOGLE_CLIENT_SECRET}"  # Environment variable
    redirectURL: "http://localhost:8080/auth/google/callback"
    scopes:
      - "https://www.googleapis.com/auth/userinfo.email"
      - "https://www.googleapis.com/auth/userinfo.profile"

# Email Service Configuration
email:
  smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "${EMAIL_USERNAME}"  # Environment variable (Gmail address)
    password: "${EMAIL_PASSWORD}"  # Environment variable (App Password)
    from: "${EMAIL_FROM}"          # Environment variable
    fromName: "Tzlev Support"
  templates:
    path: "templates/email"

# Session Configuration
session:
  store: "redis"  # redis, memory (for development)
  secretKey: "${SESSION_SECRET}"  # Environment variable
  cookieName: "tzlev_session"
  cookieMaxAge: 86400  # 24 hours in seconds
  cookieSecure: false  # Set true in production with HTTPS
  cookieHttpOnly: true
  cookieSameSite: "lax"  # lax, strict, none

# Security Configuration
security:
  jwtSecret: "${JWT_SECRET}"  # Environment variable
  jwtExpiration: "1h"
  bcryptCost: 12  # Password hashing cost

# Logging Configuration
logging:
  level: "debug"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, file
  filePath: "logs/app.log"
```

### Environment Variables

Create a `.env` file in the project root for sensitive configuration:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tzlev_user
DB_PASSWORD=your_database_password
DB_NAME=tzlev_db

# Redis
REDIS_PASSWORD=your_redis_password

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Email
EMAIL_USERNAME=your_email@gmail.com
EMAIL_PASSWORD=your_gmail_app_password
EMAIL_FROM=noreply@tzlev.com

# Session & Security
SESSION_SECRET=your_random_session_secret_key
JWT_SECRET=your_random_jwt_secret_key
```

### Configuration Loading

The configuration is loaded in the following order (later values override earlier ones):

1. `config/config.yaml` (base configuration)
2. `config/config.{environment}.yaml` (environment-specific overrides)
3. Environment variables (highest priority)

Example implementation in `main.go`:

```go
import (
    "github.com/gogf/gf/v2/os/gcfg"
    "github.com/gogf/gf/v2/frame/g"
)

func loadConfig() error {
    // Load base configuration
    cfg := g.Cfg()

    // Load environment-specific configuration
    env := cfg.MustGet(ctx, "app.environment").String()
    if env != "" {
        envConfigFile := fmt.Sprintf("config/config.%s.yaml", env)
        if gfile.Exists(envConfigFile) {
            cfg.GetAdapter().(*gcfg.AdapterFile).SetFileName(envConfigFile)
        }
    }

    return nil
}
```

### Accessing Configuration

Configuration values are accessed using GoFrame's `g.Cfg()` API:

```go
// In services or controllers
port := g.Cfg().MustGet(ctx, "app.port").Int()
dbHost := g.Cfg().MustGet(ctx, "database.host").String()
redisPassword := g.Cfg().MustGet(ctx, "redis.password").String()
```

## Database Architecture

### Database Stack

- **Database**: PostgreSQL 16+
- **ORM**: GoFrame Database (gdb) - built-in ORM
- **Migrations**: golang-migrate or custom migration system

### Database Connection

Database connection is initialized on application startup using GoFrame's gdb:

```go
// internal/database/database.go
package database

import (
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

// Init initializes database connection using GoFrame's gdb
// Configuration is read from config.yaml automatically
func Init() error {
    ctx := gctx.New()

    // Test database connection
    db := g.DB()
    if err := db.PingMaster(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    g.Log().Info(ctx, "Database connection established")

    return nil
}

// GetDB returns the default database instance
func GetDB() gdb.DB {
    return g.DB()
}

// Close closes the database connection
func Close() error {
    ctx := gctx.New()
    if err := g.DB().Close(ctx); err != nil {
        return err
    }
    return nil
}
```

### Model Layer

Models are defined as simple Go structs in `internal/model/`:

```go
// internal/model/user.go
package model

import (
    "github.com/gogf/gf/v2/os/gtime"
)

type User struct {
    Id          uint        `json:"id"`
    CreatedAt   *gtime.Time `json:"created_at"`
    UpdatedAt   *gtime.Time `json:"updated_at"`
    DeletedAt   *gtime.Time `json:"deleted_at,omitempty"`
    Email       string      `json:"email"`
    Name        string      `json:"name"`
    Picture     string      `json:"picture"`
    GoogleId    string      `json:"google_id"`
    IsActive    bool        `json:"is_active"`
    LastLoginAt *gtime.Time `json:"last_login_at,omitempty"`
}
```

**Note**: GoFrame uses `gtime.Time` for timestamp fields which provides better timezone handling and JSON serialization.

### Repository Layer

Repositories provide data access abstraction in `internal/repository/`:

```go
// internal/repository/user_repository.go
package repository

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
    "tzlev/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
    return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    user.CreatedAt = gtime.Now()
    user.UpdatedAt = gtime.Now()

    _, err := g.DB().Ctx(ctx).Insert("users", user)
    return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := g.DB().Ctx(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        Scan(&user)

    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := g.DB().Ctx(ctx).
        Where("email = ? AND deleted_at IS NULL", email).
        Scan(&user)

    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
    var user model.User
    err := g.DB().Ctx(ctx).
        Where("google_id = ? AND deleted_at IS NULL", googleID).
        Scan(&user)

    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    user.UpdatedAt = gtime.Now()

    _, err := g.DB().Ctx(ctx).
        Where("id = ?", user.Id).
        Update("users", user)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    // Soft delete
    _, err := g.DB().Ctx(ctx).
        Where("id = ?", id).
        Update("users", g.Map{"deleted_at": gtime.Now()})
    return err
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]model.User, error) {
    var users []model.User
    err := g.DB().Ctx(ctx).
        Where("deleted_at IS NULL").
        Offset(offset).
        Limit(limit).
        Scan(&users)

    return users, err
}
```

### Database Migrations

Migrations are managed in `migrations/` directory:

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_sessions_table.up.sql
└── 000002_create_sessions_table.down.sql
```

Example migration file:

```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    picture VARCHAR(500),
    google_id VARCHAR(255) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_email ON users(email);
```

```sql
-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### Running Migrations

Migrations are run via CLI mode:

```bash
# Run all pending migrations
go run main.go migrate up

# Rollback last migration
go run main.go migrate down

# Rollback all migrations
go run main.go migrate reset

# Check migration status
go run main.go migrate status
```

### Database Seeding

Seed data for development:

```go
// internal/cli/seed.go
package cli

import (
    "context"
    "tzlev/internal/database"
    "tzlev/internal/model"
)

func SeedDatabase(ctx context.Context) error {
    // Seed users
    users := []model.User{
        {
            Email:    "admin@tzlev.com",
            Name:     "Admin User",
            GoogleID: "test_google_id_1",
            IsActive: true,
        },
        {
            Email:    "user@tzlev.com",
            Name:     "Regular User",
            GoogleID: "test_google_id_2",
            IsActive: true,
        },
    }

    for _, user := range users {
        if err := database.DB.Create(&user).Error; err != nil {
            return err
        }
    }

    return nil
}
```

Run seeding:

```bash
# Seed the database
go run main.go seed
```

## Redis Architecture

### Redis Stack

- **Redis**: Redis 7+ (in-memory data store)
- **Client**: go-redis/redis (official Go client)
- **Use Cases**: Session storage, caching, rate limiting

### Redis Connection

Redis connection is initialized on application startup:

```go
// internal/redis/redis.go
package redis

import (
    "context"
    "fmt"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
    "github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() error {
    ctx := gctx.New()
    cfg := g.Cfg()

    // Create Redis client
    Client = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d",
            cfg.MustGet(ctx, "redis.host").String(),
            cfg.MustGet(ctx, "redis.port").Int(),
        ),
        Password: cfg.MustGet(ctx, "redis.password").String(),
        DB:       cfg.MustGet(ctx, "redis.database").Int(),
        PoolSize: cfg.MustGet(ctx, "redis.poolSize").Int(),
    })

    // Test connection
    if err := Client.Ping(ctx).Err(); err != nil {
        return fmt.Errorf("failed to connect to Redis: %w", err)
    }

    g.Log().Info(ctx, "Redis connection established")

    return nil
}

func Close() error {
    if Client != nil {
        return Client.Close()
    }
    return nil
}
```

### Session Storage

Redis is used for storing user sessions:

```go
// internal/session/session.go
package session

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
    "tzlev/internal/redis"
)

type Session struct {
    UserID    uint      `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type SessionManager struct {
    prefix string
    ttl    time.Duration
}

func NewSessionManager() *SessionManager {
    ctx := gctx.New()
    cfg := g.Cfg()

    return &SessionManager{
        prefix: cfg.MustGet(ctx, "redis.sessionPrefix").String(),
        ttl:    cfg.MustGet(ctx, "redis.sessionTTL").Duration(),
    }
}

func (sm *SessionManager) key(sessionID string) string {
    return fmt.Sprintf("%s%s", sm.prefix, sessionID)
}

func (sm *SessionManager) Create(ctx context.Context, sessionID string, session *Session) error {
    session.CreatedAt = time.Now()

    data, err := json.Marshal(session)
    if err != nil {
        return fmt.Errorf("failed to marshal session: %w", err)
    }

    key := sm.key(sessionID)
    return redis.Client.Set(ctx, key, data, sm.ttl).Err()
}

func (sm *SessionManager) Get(ctx context.Context, sessionID string) (*Session, error) {
    key := sm.key(sessionID)

    data, err := redis.Client.Get(ctx, key).Bytes()
    if err != nil {
        return nil, err
    }

    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, fmt.Errorf("failed to unmarshal session: %w", err)
    }

    return &session, nil
}

func (sm *SessionManager) Delete(ctx context.Context, sessionID string) error {
    key := sm.key(sessionID)
    return redis.Client.Del(ctx, key).Err()
}

func (sm *SessionManager) Refresh(ctx context.Context, sessionID string) error {
    key := sm.key(sessionID)
    return redis.Client.Expire(ctx, key, sm.ttl).Err()
}
```

### Caching Layer

Redis is used for caching frequently accessed data:

```go
// internal/cache/cache.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
    "tzlev/internal/redis"
)

type CacheManager struct {
    prefix string
}

func NewCacheManager() *CacheManager {
    ctx := gctx.New()
    cfg := g.Cfg()

    return &CacheManager{
        prefix: cfg.MustGet(ctx, "redis.cachePrefix").String(),
    }
}

func (cm *CacheManager) key(cacheKey string) string {
    return fmt.Sprintf("%s%s", cm.prefix, cacheKey)
}

func (cm *CacheManager) Set(ctx context.Context, cacheKey string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal cache value: %w", err)
    }

    key := cm.key(cacheKey)
    return redis.Client.Set(ctx, key, data, ttl).Err()
}

func (cm *CacheManager) Get(ctx context.Context, cacheKey string, dest interface{}) error {
    key := cm.key(cacheKey)

    data, err := redis.Client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }

    if err := json.Unmarshal(data, dest); err != nil {
        return fmt.Errorf("failed to unmarshal cache value: %w", err)
    }

    return nil
}

func (cm *CacheManager) Delete(ctx context.Context, cacheKey string) error {
    key := cm.key(cacheKey)
    return redis.Client.Del(ctx, key).Err()
}

func (cm *CacheManager) DeletePattern(ctx context.Context, pattern string) error {
    key := cm.key(pattern)

    // Find all keys matching the pattern
    keys, err := redis.Client.Keys(ctx, key).Result()
    if err != nil {
        return err
    }

    if len(keys) == 0 {
        return nil
    }

    // Delete all matching keys
    return redis.Client.Del(ctx, keys...).Err()
}

func (cm *CacheManager) Exists(ctx context.Context, cacheKey string) (bool, error) {
    key := cm.key(cacheKey)
    count, err := redis.Client.Exists(ctx, key).Result()
    return count > 0, err
}
```

### Usage in Services

Example of using cache in a service:

```go
// internal/service/user_service.go
package service

import (
    "context"
    "fmt"
    "time"

    "tzlev/internal/cache"
    "tzlev/internal/model"
    "tzlev/internal/repository"
)

type UserService struct {
    userRepo     *repository.UserRepository
    cacheManager *cache.CacheManager
}

func NewUserService() *UserService {
    return &UserService{
        userRepo:     repository.NewUserRepository(),
        cacheManager: cache.NewCacheManager(),
    }
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)

    // Try to get from cache first
    var user model.User
    err := s.cacheManager.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }

    // If not in cache, get from database
    dbUser, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Store in cache for 5 minutes
    _ = s.cacheManager.Set(ctx, cacheKey, dbUser, 5*time.Minute)

    return dbUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
    // Update in database
    if err := s.userRepo.Update(ctx, user); err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    _ = s.cacheManager.Delete(ctx, cacheKey)

    return nil
}
```

### Redis Data Patterns

**Session Keys:**
```
tzlev:session:abc123...  # Session data
```

**Cache Keys:**
```
tzlev:cache:user:1       # User data cache
tzlev:cache:config:*     # Configuration cache
tzlev:cache:list:*       # List data cache
```

## Authentication Architecture

### Authentication Stack

- **OAuth Provider**: Google OAuth 2.0
- **Session Management**: Redis-backed sessions
- **JWT**: JSON Web Tokens for API authentication
- **Library**: golang.org/x/oauth2

### OAuth Configuration

OAuth is configured in the configuration file:

```yaml
oauth:
  google:
    clientID: "${GOOGLE_CLIENT_ID}"
    clientSecret: "${GOOGLE_CLIENT_SECRET}"
    redirectURL: "http://localhost:8080/auth/google/callback"
    scopes:
      - "https://www.googleapis.com/auth/userinfo.email"
      - "https://www.googleapis.com/auth/userinfo.profile"
```

### OAuth Flow

```go
// internal/oauth/google.go
package oauth

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth() error {
    ctx := gctx.New()
    cfg := g.Cfg()

    GoogleOAuthConfig = &oauth2.Config{
        ClientID:     cfg.MustGet(ctx, "oauth.google.clientID").String(),
        ClientSecret: cfg.MustGet(ctx, "oauth.google.clientSecret").String(),
        RedirectURL:  cfg.MustGet(ctx, "oauth.google.redirectURL").String(),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }

    g.Log().Info(ctx, "Google OAuth initialized")

    return nil
}

type GoogleUserInfo struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
    Locale        string `json:"locale"`
}

func GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
    client := GoogleOAuthConfig.Client(ctx, token)

    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }
    defer resp.Body.Close()

    var userInfo GoogleUserInfo
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return nil, fmt.Errorf("failed to decode user info: %w", err)
    }

    return &userInfo, nil
}
```

### Authentication Controller

```go
// internal/controller/auth_controller.go
package controller

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "tzlev/internal/model"
    "tzlev/internal/oauth"
    "tzlev/internal/repository"
    "tzlev/internal/session"
)

type AuthController struct {
    userRepo       *repository.UserRepository
    sessionManager *session.SessionManager
}

func NewAuthController() *AuthController {
    return &AuthController{
        userRepo:       repository.NewUserRepository(),
        sessionManager: session.NewSessionManager(),
    }
}

func (c *AuthController) GoogleLogin(r *ghttp.Request) {
    // Generate random state for CSRF protection
    b := make([]byte, 32)
    rand.Read(b)
    state := base64.URLEncoding.EncodeToString(b)

    // Store state in session
    r.Session.Set("oauth_state", state)

    // Redirect to Google OAuth
    url := oauth.GoogleOAuthConfig.AuthCodeURL(state)
    r.Response.RedirectTo(url)
}

func (c *AuthController) GoogleCallback(r *ghttp.Request) {
    ctx := r.Context()

    // Verify state
    state := r.Get("state").String()
    sessionState := r.Session.Get("oauth_state").String()

    if state == "" || state != sessionState {
        r.Response.WriteJson(g.Map{
            "error": "Invalid state parameter",
        })
        return
    }

    // Exchange code for token
    code := r.Get("code").String()
    token, err := oauth.GoogleOAuthConfig.Exchange(ctx, code)
    if err != nil {
        g.Log().Error(ctx, "Failed to exchange token:", err)
        r.Response.WriteJson(g.Map{
            "error": "Failed to exchange token",
        })
        return
    }

    // Get user info from Google
    googleUser, err := oauth.GetGoogleUserInfo(ctx, token)
    if err != nil {
        g.Log().Error(ctx, "Failed to get user info:", err)
        r.Response.WriteJson(g.Map{
            "error": "Failed to get user info",
        })
        return
    }

    // Find or create user
    user, err := c.userRepo.FindByGoogleID(ctx, googleUser.ID)
    if err != nil {
        // User doesn't exist, create new user
        now := time.Now()
        user = &model.User{
            Email:       googleUser.Email,
            Name:        googleUser.Name,
            Picture:     googleUser.Picture,
            GoogleID:    googleUser.ID,
            IsActive:    true,
            LastLoginAt: &now,
        }

        if err := c.userRepo.Create(ctx, user); err != nil {
            g.Log().Error(ctx, "Failed to create user:", err)
            r.Response.WriteJson(g.Map{
                "error": "Failed to create user",
            })
            return
        }
    } else {
        // Update last login time
        now := time.Now()
        user.LastLoginAt = &now
        if err := c.userRepo.Update(ctx, user); err != nil {
            g.Log().Error(ctx, "Failed to update user:", err)
        }
    }

    // Create session
    sessionID := r.Session.Id()
    sess := &session.Session{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
    }

    if err := c.sessionManager.Create(ctx, sessionID, sess); err != nil {
        g.Log().Error(ctx, "Failed to create session:", err)
        r.Response.WriteJson(g.Map{
            "error": "Failed to create session",
        })
        return
    }

    // Set session cookie
    r.Session.Set("user_id", user.ID)
    r.Session.Set("user_email", user.Email)
    r.Session.Set("user_name", user.Name)

    // Redirect to home page
    r.Response.RedirectTo("/")
}

func (c *AuthController) Logout(r *ghttp.Request) {
    ctx := r.Context()

    // Delete session from Redis
    sessionID := r.Session.Id()
    if err := c.sessionManager.Delete(ctx, sessionID); err != nil {
        g.Log().Error(ctx, "Failed to delete session:", err)
    }

    // Clear session cookie
    r.Session.Clear()

    r.Response.WriteJson(g.Map{
        "status":  "ok",
        "message": "Logged out successfully",
    })
}

func (c *AuthController) GetCurrentUser(r *ghttp.Request) {
    ctx := r.Context()

    // Get session from Redis
    sessionID := r.Session.Id()
    sess, err := c.sessionManager.Get(ctx, sessionID)
    if err != nil {
        r.Response.Status = 401
        r.Response.WriteJson(g.Map{
            "error": "Not authenticated",
        })
        return
    }

    r.Response.WriteJson(g.Map{
        "user": sess,
    })
}
```

### Authentication Middleware

```go
// internal/middleware/auth.go
package middleware

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "tzlev/internal/session"
)

func Auth() func(r *ghttp.Request) {
    sessionManager := session.NewSessionManager()

    return func(r *ghttp.Request) {
        ctx := r.Context()

        // Get session ID from cookie
        sessionID := r.Session.Id()
        if sessionID == "" {
            r.Response.Status = 401
            r.Response.WriteJson(g.Map{
                "error": "Not authenticated",
            })
            return
        }

        // Verify session in Redis
        sess, err := sessionManager.Get(ctx, sessionID)
        if err != nil {
            g.Log().Warning(ctx, "Invalid session:", err)
            r.Response.Status = 401
            r.Response.WriteJson(g.Map{
                "error": "Invalid or expired session",
            })
            return
        }

        // Refresh session TTL
        if err := sessionManager.Refresh(ctx, sessionID); err != nil {
            g.Log().Warning(ctx, "Failed to refresh session:", err)
        }

        // Store user info in request context
        r.SetCtxVar("user_id", sess.UserID)
        r.SetCtxVar("user_email", sess.Email)
        r.SetCtxVar("user_name", sess.Name)

        r.Middleware.Next()
    }
}
```

### Authentication Routes

Add authentication routes to `main.go`:

```go
// In main.go
func setupRoutes(s *ghttp.Server) {
    authCtrl := controller.NewAuthController()

    // Public routes
    s.Group("/auth", func(group *ghttp.RouterGroup) {
        group.GET("/google/login", authCtrl.GoogleLogin)
        group.GET("/google/callback", authCtrl.GoogleCallback)
        group.POST("/logout", authCtrl.Logout)
    })

    // Protected API routes
    s.Group("/api", func(group *ghttp.RouterGroup) {
        // Apply authentication middleware
        group.Middleware(middleware.Auth())

        group.GET("/me", authCtrl.GetCurrentUser)
        // Add other protected routes here
    })
}
```

### JWT Token Generation (Optional)

For API-based authentication (mobile apps, external clients):

```go
// internal/jwt/jwt.go
package jwt

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, name string) (string, error) {
    ctx := gctx.New()
    cfg := g.Cfg()

    secret := cfg.MustGet(ctx, "security.jwtSecret").String()
    expiration := cfg.MustGet(ctx, "security.jwtExpiration").Duration()

    claims := Claims{
        UserID: userID,
        Email:  email,
        Name:   name,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (*Claims, error) {
    ctx := gctx.New()
    cfg := g.Cfg()

    secret := cfg.MustGet(ctx, "security.jwtSecret").String()

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

## Email Service Architecture

### Email Stack

- **SMTP Provider**: Gmail SMTP
- **Template Engine**: Go html/template
- **Library**: net/smtp (standard library) or gomail

### Email Configuration

Email service is configured in the configuration file:

```yaml
email:
  smtp:
    host: "smtp.gmail.com"
    port: 587
    username: "${EMAIL_USERNAME}"
    password: "${EMAIL_PASSWORD}"
    from: "${EMAIL_FROM}"
    fromName: "Tzlev Support"
  templates:
    path: "templates/email"
```

### Email Service Implementation

```go
// internal/email/email.go
package email

import (
    "bytes"
    "fmt"
    "html/template"
    "net/smtp"
    "path/filepath"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

type EmailService struct {
    host         string
    port         int
    username     string
    password     string
    from         string
    fromName     string
    templatePath string
}

func NewEmailService() *EmailService {
    ctx := gctx.New()
    cfg := g.Cfg()

    return &EmailService{
        host:         cfg.MustGet(ctx, "email.smtp.host").String(),
        port:         cfg.MustGet(ctx, "email.smtp.port").Int(),
        username:     cfg.MustGet(ctx, "email.smtp.username").String(),
        password:     cfg.MustGet(ctx, "email.smtp.password").String(),
        from:         cfg.MustGet(ctx, "email.smtp.from").String(),
        fromName:     cfg.MustGet(ctx, "email.smtp.fromName").String(),
        templatePath: cfg.MustGet(ctx, "email.templates.path").String(),
    }
}

type EmailData struct {
    To      []string
    Subject string
    Body    string
}

func (es *EmailService) Send(data *EmailData) error {
    // Authentication
    auth := smtp.PlainAuth("", es.username, es.password, es.host)

    // Build message
    msg := es.buildMessage(data)

    // Send email
    addr := fmt.Sprintf("%s:%d", es.host, es.port)
    if err := smtp.SendMail(addr, auth, es.from, data.To, msg); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}

func (es *EmailService) buildMessage(data *EmailData) []byte {
    msg := fmt.Sprintf("From: %s <%s>\r\n", es.fromName, es.from)
    msg += fmt.Sprintf("To: %s\r\n", data.To[0])
    msg += fmt.Sprintf("Subject: %s\r\n", data.Subject)
    msg += "MIME-Version: 1.0\r\n"
    msg += "Content-Type: text/html; charset=UTF-8\r\n"
    msg += "\r\n"
    msg += data.Body

    return []byte(msg)
}

type WelcomeEmailData struct {
    Name string
    Link string
}

func (es *EmailService) SendWelcomeEmail(to, name string) error {
    // Parse template
    tmplPath := filepath.Join(es.templatePath, "welcome.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    // Execute template
    var body bytes.Buffer
    data := WelcomeEmailData{
        Name: name,
        Link: "https://tzlev.com/dashboard",
    }

    if err := tmpl.Execute(&body, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    // Send email
    return es.Send(&EmailData{
        To:      []string{to},
        Subject: "Welcome to Tzlev!",
        Body:    body.String(),
    })
}

type PasswordResetEmailData struct {
    Name      string
    ResetLink string
}

func (es *EmailService) SendPasswordResetEmail(to, name, resetToken string) error {
    // Parse template
    tmplPath := filepath.Join(es.templatePath, "password_reset.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    // Execute template
    var body bytes.Buffer
    data := PasswordResetEmailData{
        Name:      name,
        ResetLink: fmt.Sprintf("https://tzlev.com/reset-password?token=%s", resetToken),
    }

    if err := tmpl.Execute(&body, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    // Send email
    return es.Send(&EmailData{
        To:      []string{to},
        Subject: "Reset Your Password",
        Body:    body.String(),
    })
}

type NotificationEmailData struct {
    Name    string
    Message string
    Link    string
}

func (es *EmailService) SendNotificationEmail(to, name, message, link string) error {
    // Parse template
    tmplPath := filepath.Join(es.templatePath, "notification.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    // Execute template
    var body bytes.Buffer
    data := NotificationEmailData{
        Name:    name,
        Message: message,
        Link:    link,
    }

    if err := tmpl.Execute(&body, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    // Send email
    return es.Send(&EmailData{
        To:      []string{to},
        Subject: "New Notification",
        Body:    body.String(),
    })
}
```

### Email Templates

Email templates are stored in `templates/email/` directory:

#### Welcome Email Template (`templates/email/welcome.html`)

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #4F46E5;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .content {
            background-color: #f9fafb;
            padding: 30px;
            border-radius: 0 0 5px 5px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #4F46E5;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            color: #6b7280;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Welcome to Tzlev!</h1>
    </div>
    <div class="content">
        <p>Hello {{.Name}},</p>
        <p>Welcome to Tzlev! We're excited to have you on board.</p>
        <p>Get started by exploring your dashboard and discovering all the features we have to offer.</p>
        <a href="{{.Link}}" class="button">Go to Dashboard</a>
    </div>
    <div class="footer">
        <p>&copy; 2024 Tzlev. All rights reserved.</p>
    </div>
</body>
</html>
```

#### Password Reset Email Template (`templates/email/password_reset.html`)

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #DC2626;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .content {
            background-color: #f9fafb;
            padding: 30px;
            border-radius: 0 0 5px 5px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #DC2626;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .warning {
            background-color: #FEF3C7;
            padding: 15px;
            border-radius: 5px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            color: #6b7280;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Reset Your Password</h1>
    </div>
    <div class="content">
        <p>Hello {{.Name}},</p>
        <p>We received a request to reset your password. Click the button below to create a new password:</p>
        <a href="{{.ResetLink}}" class="button">Reset Password</a>
        <div class="warning">
            <strong>Security Notice:</strong> This link will expire in 1 hour. If you didn't request this password reset, please ignore this email.
        </div>
    </div>
    <div class="footer">
        <p>&copy; 2024 Tzlev. All rights reserved.</p>
    </div>
</body>
</html>
```

#### Notification Email Template (`templates/email/notification.html`)

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #059669;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .content {
            background-color: #f9fafb;
            padding: 30px;
            border-radius: 0 0 5px 5px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #059669;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            color: #6b7280;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>New Notification</h1>
    </div>
    <div class="content">
        <p>Hello {{.Name}},</p>
        <p>{{.Message}}</p>
        <a href="{{.Link}}" class="button">View Details</a>
    </div>
    <div class="footer">
        <p>&copy; 2024 Tzlev. All rights reserved.</p>
    </div>
</body>
</html>
```

### Usage in Services

Example of using the email service:

```go
// internal/service/user_service.go
package service

import (
    "context"
    "tzlev/internal/email"
    "tzlev/internal/model"
    "tzlev/internal/repository"
)

type UserService struct {
    userRepo     *repository.UserRepository
    emailService *email.EmailService
}

func NewUserService() *UserService {
    return &UserService{
        userRepo:     repository.NewUserRepository(),
        emailService: email.NewEmailService(),
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
    // Create user in database
    if err := s.userRepo.Create(ctx, user); err != nil {
        return err
    }

    // Send welcome email
    go func() {
        if err := s.emailService.SendWelcomeEmail(user.Email, user.Name); err != nil {
            g.Log().Error(ctx, "Failed to send welcome email:", err)
        }
    }()

    return nil
}
```

### Gmail App Password Setup

To use Gmail SMTP, you need to create an App Password:

1. Go to your Google Account settings
2. Navigate to Security > 2-Step Verification
3. Scroll down and select "App passwords"
4. Generate a new app password for "Mail"
5. Use this password in your `.env` file as `EMAIL_PASSWORD`

## Project Structure

```
tzlev/
├── main.go                      # Application entry point (dual-mode)
│
├── config/                      # Configuration files
│   ├── config.yaml              # Main configuration (defaults)
│   ├── config.example.yaml      # Example configuration
│   ├── config.development.yaml  # Development overrides
│   └── config.production.yaml   # Production overrides
│
├── migrations/                  # Database migrations (SQL)
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_sessions_table.up.sql
│   └── 000002_create_sessions_table.down.sql
│
├── templates/                   # Template files
│   └── email/                   # Email templates (HTML)
│       ├── welcome.html
│       ├── password_reset.html
│       └── notification.html
│
├── internal/                    # Private application code
│   ├── controller/              # HTTP request handlers
│   │   ├── auth_controller.go   # Authentication endpoints
│   │   └── ...                  # Other controllers
│   ├── service/                 # Business logic layer
│   │   ├── user_service.go      # User business logic
│   │   └── ...                  # Other services
│   ├── repository/              # Data access layer
│   │   ├── user_repository.go   # User database operations
│   │   └── ...                  # Other repositories
│   ├── model/                   # Database models (Go structs)
│   │   ├── user.go              # User model
│   │   └── ...                  # Other models
│   ├── database/                # Database connection
│   │   └── database.go          # PostgreSQL + GoFrame gdb setup
│   ├── redis/                   # Redis connection
│   │   └── redis.go             # Redis client setup
│   ├── session/                 # Session management
│   │   └── session.go           # Redis-backed sessions
│   ├── cache/                   # Caching layer
│   │   └── cache.go             # Redis-backed cache
│   ├── oauth/                   # OAuth handlers
│   │   └── google.go            # Google OAuth 2.0
│   ├── jwt/                     # JWT tokens
│   │   └── jwt.go               # Token generation/verification
│   ├── email/                   # Email service
│   │   └── email.go             # SMTP email sender
│   ├── middleware/              # HTTP middleware
│   │   └── auth.go              # Authentication middleware
│   └── cli/                     # CLI command handlers
│       ├── migrate.go           # Database migration commands
│       ├── seed.go              # Database seeding commands
│       ├── user.go              # User management commands
│       └── version.go           # Version information
│
├── assets/                      # Source template assets (Wowdash)
│   ├── css/                     # Template stylesheets
│   ├── js/                      # Template JavaScript libraries
│   └── images/                  # Template images
│   └── ...                      # (Copy to public/ before running)
│
├── frontend/                    # React source code
│   ├── src/
│   │   ├── components/          # React components
│   │   │   ├── Sidebar.jsx      # Left navigation with routing
│   │   │   ├── Navbar.jsx       # Top navigation with language switcher
│   │   │   ├── MainContent.jsx  # Main content area
│   │   │   ├── Footer.jsx       # Footer
│   │   │   └── ThemeCustomization.jsx  # Theme settings
│   │   ├── pages/               # Route pages
│   │   │   ├── Home.jsx         # Home page
│   │   │   ├── Dashboard.jsx    # Dashboard page
│   │   │   └── ...              # Other pages
│   │   ├── locales/             # i18n translation files
│   │   │   ├── en/
│   │   │   │   └── translation.json  # English translations
│   │   │   └── he/
│   │   │       └── translation.json  # Hebrew translations
│   │   ├── config/              # Configuration files
│   │   │   ├── menuConfig.js    # Sidebar menu configuration
│   │   │   └── i18n.js          # i18next configuration
│   │   ├── App.jsx              # Root React component with Router
│   │   └── main.jsx             # React entry point
│   ├── index.html               # HTML template
│   ├── vite.config.js           # Build configuration
│   └── package.json             # Node dependencies
│
├── public/                      # Served by GoFrame (production files)
│   ├── css/                     # Stylesheets (copied from assets/)
│   ├── js/                      # JavaScript libraries (copied from assets/)
│   ├── images/                  # Images (copied from assets/)
│   ├── index.html               # Built HTML (entry point)
│   └── assets/main-*.js         # Built React bundle
│
├── .env                         # Environment variables (secrets)
├── .env.example                 # Example environment variables
├── go.mod                       # Go dependencies
├── go.sum                       # Go dependencies checksums
├── README.md                    # User documentation
└── DESIGN.md                    # This file
```

### Directory Responsibilities

**Root Level**:
- `main.go`: Application entry point with dual-mode support
  - Detects CLI arguments vs web mode
  - Web mode: Starts GoFrame HTTP server with CORS
  - CLI mode: Executes commands and exits

**config/**:
- YAML configuration files
- Base configuration and environment-specific overrides
- Loaded in order: base → environment → env variables

**migrations/**:
- SQL migration files for database schema changes
- Up/down migrations for version control
- Run via CLI: `go run main.go migrate up`

**templates/**:
- HTML template files
- Email templates with Go html/template syntax
- Used by email service for transactional emails

**internal/controller/**:
- HTTP request handlers
- Parse requests, call services, return responses
- Thin layer - business logic in services

**internal/service/**:
- Business logic layer
- Orchestrates repositories, cache, email, etc.
- Transactional operations

**internal/repository/**:
- Data access layer
- CRUD operations on database models
- Uses GoFrame gdb for database queries

**internal/model/**:
- Database models as Go structs
- Struct definitions matching database tables
- JSON serialization tags for API responses

**internal/database/**:
- Database connection initialization
- Connection pooling configuration
- Health checks

**internal/redis/**:
- Redis client initialization
- Connection pooling configuration
- Health checks

**internal/session/**:
- Session management using Redis
- Create, get, delete, refresh sessions
- TTL-based expiration

**internal/cache/**:
- Caching layer using Redis
- Set, get, delete cache entries
- Pattern-based deletion

**internal/oauth/**:
- OAuth 2.0 provider integrations
- Google OAuth configuration
- User info retrieval from OAuth providers

**internal/jwt/**:
- JWT token generation and verification
- Claims structure
- Token signing and validation

**internal/email/**:
- Email service using SMTP
- Template-based email rendering
- Async email sending

**internal/middleware/**:
- HTTP middleware functions
- Authentication, logging, rate limiting
- Applied to route groups

**internal/cli/**:
- CLI command handlers
- Uses Cobra CLI framework
- Shares access to services and business logic
- Each file represents a command group

**assets/**:
- Source Wowdash template files
- Never served directly
- Copied to `public/` during setup

**frontend/**:
- React source code with Vite
- Components, pages, locales, config
- Built to `public/` directory

**public/**:
- All files served by GoFrame (web mode only)
- Contains built React app and copied template files (css/, js/, images/)
- Only directory accessible via HTTP

## Component Architecture

### GoFrame Backend (main.go)

**Dual-Mode Architecture**:

The application operates in two modes based on command-line arguments:

**Mode Detection**:
```go
func main() {
    ctx := gctx.New()

    // Check for CLI arguments
    if len(os.Args) > 1 {
        // CLI Mode - execute command and exit
        runCLI(ctx)
        return
    }

    // Web Mode - start HTTP server
    runWebServer(ctx)
}
```

**Web Mode Responsibilities**:
1. Configure CORS middleware for API access
2. Serve static files from `public/` directory
3. Serve built React application (index.html)
4. Provide RESTful API endpoints
5. Handle business logic through controllers/services

**Web Mode Key Routes**:
- `GET /` → Serves `public/index.html` (React app entry)
- `GET /css/*`, `/js/*`, `/images/*` → Serves static assets from `public/` (copied from assets/)
- `GET /api/*` → API endpoints for backend logic (CORS enabled)

**Web Mode Server Configuration**:
```go
func runWebServer(ctx context.Context) {
    s := g.Server()

    // CORS Middleware - enables cross-origin API calls
    s.Use(func(r *ghttp.Request) {
        r.Response.CORSDefault()
        r.Middleware.Next()
    })

    // Serve static files from public directory
    s.SetServerRoot("public")

    // Main route - serve the index page
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.ServeFile("public/index.html")
    })

    // API routes with CORS enabled
    s.Group("/api", func(group *ghttp.RouterGroup) {
        // API endpoints here
    })

    // Start server
    g.Log().Info(ctx, "Starting Tzlev server on http://localhost:8080")
    s.SetPort(8080)
    s.Run()
}
```

**CLI Mode Responsibilities**:
1. Parse command-line arguments
2. Execute administrative tasks
3. Access services and business logic
4. Return appropriate exit codes
5. No HTTP server started

### React Frontend

**Component Hierarchy**:
```
main.jsx
└── i18next Provider
    └── BrowserRouter (React Router)
        └── App.jsx (Root)
            ├── ThemeCustomization (Floating theme settings panel)
            ├── Sidebar (Left navigation with routing)
            └── Main (Dashboard area)
                ├── Navbar (Top navigation with language switcher)
                ├── Routes (Route-based page rendering)
                │   ├── Home
                │   ├── Dashboard
                │   └── Other Pages
                └── Footer
```

**Key Features**:
1. **App.jsx**:
   - Root component
   - Loads template JavaScript (app.js) via useEffect
   - Manages global layout structure

2. **Sidebar.jsx**:
   - Navigation menu
   - Dropdown menus for sections
   - Logo display

3. **Navbar.jsx**:
   - Search bar
   - Theme toggle
   - User profile dropdown
   - Notification/message dropdowns

4. **MainContent.jsx**:
   - Main content area
   - Breadcrumb navigation
   - API integration example
   - Dynamic content rendering

5. **ThemeCustomization.jsx**:
   - Theme mode selector (light/dark/system)
   - Page direction (LTR/RTL)
   - Color scheme selector

## CLI Mode Architecture

### Overview
The application supports CLI mode for administrative tasks, database operations, and system management without starting the web server.

### CLI Command Structure

**Main CLI Handler** (`main.go - runCLI function`):
```go
func runCLI(ctx context.Context) {
    // Parse commands using GoFrame's gcmd
    parser, err := gcmd.Parse(nil)
    if err != nil {
        g.Log().Fatal(ctx, "Failed to parse CLI args:", err)
    }

    command := parser.GetArg(1).String()

    switch command {
    case "migrate":
        cli.RunMigrate(ctx, parser)
    case "seed":
        cli.RunSeed(ctx, parser)
    case "user:create":
        cli.CreateUser(ctx, parser)
    case "user:list":
        cli.ListUsers(ctx, parser)
    case "version":
        cli.ShowVersion(ctx)
    case "help":
        cli.ShowHelp(ctx)
    default:
        g.Log().Errorf(ctx, "Unknown command: %s", command)
        cli.ShowHelp(ctx)
        os.Exit(1)
    }
}
```

### Available CLI Commands

#### 1. Database Migration
```bash
# Run all pending migrations
./tzlev migrate

# Run migrations with options
go run main.go migrate --up
go run main.go migrate --down
go run main.go migrate --version=20240101
```

**Implementation** (`internal/cli/migrate.go`):
```go
package cli

func RunMigrate(ctx context.Context, parser *gcmd.Parser) {
    action := parser.GetOpt("action", "up").String()

    g.Log().Info(ctx, "Running database migrations...")

    // Migration logic here
    // Access database through services

    g.Log().Info(ctx, "Migrations completed successfully")
    os.Exit(0)
}
```

#### 2. Database Seeding
```bash
# Seed database with test data
./tzlev seed

# Seed specific tables
go run main.go seed --table=users
go run main.go seed --env=development
```

**Implementation** (`internal/cli/seed.go`):
```go
package cli

func RunSeed(ctx context.Context, parser *gcmd.Parser) {
    table := parser.GetOpt("table", "all").String()
    env := parser.GetOpt("env", "development").String()

    g.Log().Infof(ctx, "Seeding database (table: %s, env: %s)...", table, env)

    // Seeding logic here
    // Use internal/service for business logic

    g.Log().Info(ctx, "Database seeded successfully")
    os.Exit(0)
}
```

#### 3. User Management
```bash
# Create new user
./tzlev user:create --email=admin@example.com --role=admin

# List all users
./tzlev user:list

# List with filters
go run main.go user:list --role=admin --active=true
```

**Implementation** (`internal/cli/user.go`):
```go
package cli

func CreateUser(ctx context.Context, parser *gcmd.Parser) {
    email := parser.GetOpt("email").String()
    role := parser.GetOpt("role", "user").String()

    if email == "" {
        g.Log().Error(ctx, "Email is required")
        os.Exit(1)
    }

    g.Log().Infof(ctx, "Creating user: %s (role: %s)", email, role)

    // Use service layer for user creation
    // userService := service.NewUserService()
    // user, err := userService.Create(ctx, email, role)

    g.Log().Info(ctx, "User created successfully")
    os.Exit(0)
}

func ListUsers(ctx context.Context, parser *gcmd.Parser) {
    role := parser.GetOpt("role").String()
    active := parser.GetOpt("active").Bool()

    g.Log().Info(ctx, "Listing users...")

    // Use service layer to fetch users
    // users, err := userService.List(ctx, filters)

    // Display users in table format
    fmt.Println("ID\tEmail\t\tRole\tActive")
    fmt.Println("---\t-----\t\t----\t------")
    // Print user list

    os.Exit(0)
}
```

#### 4. Version Information
```bash
# Show version and build info
./tzlev version
```

**Implementation** (`internal/cli/version.go`):
```go
package cli

const (
    Version = "1.0.0"
    BuildDate = "2024-01-01"
)

func ShowVersion(ctx context.Context) {
    fmt.Printf("Tzlev v%s\n", Version)
    fmt.Printf("Build Date: %s\n", BuildDate)
    fmt.Printf("Go Version: %s\n", runtime.Version())
    os.Exit(0)
}
```

#### 5. Help Command
```bash
# Show available commands
./tzlev help
./tzlev --help
```

**Implementation** (`internal/cli/help.go`):
```go
package cli

func ShowHelp(ctx context.Context) {
    helpText := `
Tzlev - GoFrame + React Application

USAGE:
    ./tzlev [command] [options]

COMMANDS:
    migrate              Run database migrations
    seed                 Seed database with test data
    user:create          Create a new user
    user:list            List all users
    version              Show version information
    help                 Show this help message

WEB MODE:
    ./tzlev              Start web server (no arguments)

EXAMPLES:
    ./tzlev migrate --up
    ./tzlev seed --table=users
    ./tzlev user:create --email=admin@example.com --role=admin
    ./tzlev user:list --role=admin

For web server mode, simply run without arguments:
    ./tzlev
    go run main.go
`
    fmt.Println(helpText)
    os.Exit(0)
}
```

### CLI Command Best Practices

1. **Error Handling**: Always check for required parameters and return appropriate exit codes
2. **Service Layer**: Use existing service layer for business logic
3. **Logging**: Use GoFrame's `g.Log()` for consistent logging
4. **Exit Codes**:
   - `0` for success
   - `1` for errors
   - `2` for invalid arguments
5. **Help Text**: Provide clear help messages for each command

### Shared Services

CLI commands share the same service layer as the web application:

```go
// Example: Using UserService in CLI
import "tzlev/internal/service"

func CreateUser(ctx context.Context, parser *gcmd.Parser) {
    // Initialize service
    userService := service.NewUserService()

    // Use service methods
    user, err := userService.Create(ctx, email, password, role)
    if err != nil {
        g.Log().Error(ctx, "Failed to create user:", err)
        os.Exit(1)
    }

    g.Log().Infof(ctx, "User created: %s (ID: %d)", user.Email, user.ID)
    os.Exit(0)
}
```

## Internationalization (i18n) Architecture

### Overview
The application supports multiple languages with RTL/LTR automatic switching using **i18next** and **react-i18next**.

### Configuration Structure

**i18n Setup** (`frontend/src/config/i18n.js`):
```javascript
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from '../locales/en/translation.json'
import he from '../locales/he/translation.json'

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: en },
      he: { translation: he }
    },
    lng: 'en', // default language
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false
    }
  })

export default i18n
```

### Translation Files

**English** (`frontend/src/locales/en/translation.json`):
```json
{
  "menu": {
    "home": "Home",
    "dashboard": "Dashboard",
    "users": "Users",
    "settings": "Settings"
  },
  "common": {
    "welcome": "Welcome",
    "logout": "Logout",
    "search": "Search"
  }
}
```

**Hebrew** (`frontend/src/locales/he/translation.json`):
```json
{
  "menu": {
    "home": "בית",
    "dashboard": "לוח בקרה",
    "users": "משתמשים",
    "settings": "הגדרות"
  },
  "common": {
    "welcome": "ברוך הבא",
    "logout": "התנתק",
    "search": "חיפוש"
  }
}
```

### Language Switching

**Implementation in Components**:
```javascript
import { useTranslation } from 'react-i18next'

function Navbar() {
  const { t, i18n } = useTranslation()

  const changeLanguage = (lng) => {
    i18n.changeLanguage(lng)
    // Update document direction
    document.documentElement.dir = lng === 'he' ? 'rtl' : 'ltr'
    document.documentElement.setAttribute('data-theme', 'light')
  }

  return (
    <div>
      <span>{t('common.welcome')}</span>
      <button onClick={() => changeLanguage('en')}>English</button>
      <button onClick={() => changeLanguage('he')}>עברית</button>
    </div>
  )
}
```

### RTL/LTR Support

**Automatic Direction Switching**:
- Hebrew (he): RTL (Right-to-Left)
- English (en): LTR (Left-to-Right)
- Updates `<html dir="rtl|ltr">` attribute
- CSS automatically adjusts layout based on direction

**CSS Considerations**:
```css
/* Template CSS already supports RTL via Bootstrap and custom styles */
/* No additional CSS modifications needed */
```

### Usage in Components

```javascript
import { useTranslation } from 'react-i18next'

function MyComponent() {
  const { t } = useTranslation()

  return (
    <div>
      <h1>{t('menu.dashboard')}</h1>
      <p>{t('common.welcome')}</p>
    </div>
  )
}
```

## React Router Architecture

### Overview
React Router v6 provides client-side routing for SPA navigation with deep linking and browser history support.

### Router Setup

**Main Entry** (`frontend/src/main.jsx`):
```javascript
import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import './config/i18n' // Initialize i18n

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
)
```

**Route Configuration** (`frontend/src/App.jsx`):
```javascript
import { Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'
import Users from './pages/Users'

function App() {
  return (
    <>
      <ThemeCustomization />
      <Sidebar />
      <main className="dashboard-main">
        <Navbar />
        <div className="dashboard-main-body">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/users" element={<Users />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </div>
        <Footer />
      </main>
    </>
  )
}
```

### Dynamic Menu Configuration

**Menu Configuration** (`frontend/src/config/menuConfig.js`):
```javascript
export const menuConfig = [
  {
    id: 'home',
    titleKey: 'menu.home',
    icon: 'solar:home-smile-angle-outline',
    path: '/',
  },
  {
    id: 'dashboard',
    titleKey: 'menu.dashboard',
    icon: 'solar:widget-4-outline',
    path: '/dashboard',
  },
  {
    id: 'users',
    titleKey: 'menu.users',
    icon: 'flowbite:users-group-outline',
    path: '/users',
    children: [
      {
        id: 'users-list',
        titleKey: 'menu.usersList',
        path: '/users/list',
      },
      {
        id: 'users-add',
        titleKey: 'menu.addUser',
        path: '/users/add',
      }
    ]
  }
]
```

### Sidebar Integration

**Dynamic Sidebar with Routing** (`frontend/src/components/Sidebar.jsx`):
```javascript
import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { menuConfig } from '../config/menuConfig'

function Sidebar() {
  const { t } = useTranslation()
  const location = useLocation()

  const isActive = (path) => location.pathname === path

  return (
    <aside className="sidebar">
      <div className="sidebar-menu-area">
        <ul className="sidebar-menu">
          {menuConfig.map((item) => (
            <li key={item.id} className={item.children ? 'dropdown' : ''}>
              {item.children ? (
                <>
                  <a href="javascript:void(0)">
                    <iconify-icon icon={item.icon} className="menu-icon" />
                    <span>{t(item.titleKey)}</span>
                  </a>
                  <ul className="sidebar-submenu">
                    {item.children.map((child) => (
                      <li key={child.id}>
                        <Link
                          to={child.path}
                          className={isActive(child.path) ? 'active' : ''}
                        >
                          <i className="ri-circle-fill circle-icon" />
                          {t(child.titleKey)}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </>
              ) : (
                <Link
                  to={item.path}
                  className={isActive(item.path) ? 'active' : ''}
                >
                  <iconify-icon icon={item.icon} className="menu-icon" />
                  <span>{t(item.titleKey)}</span>
                </Link>
              )}
            </li>
          ))}
        </ul>
      </div>
    </aside>
  )
}
```

### Page Components

**Example Page** (`frontend/src/pages/Dashboard.jsx`):
```javascript
import React from 'react'
import { useTranslation } from 'react-i18next'

function Dashboard() {
  const { t } = useTranslation()

  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <h6 className="fw-semibold mb-0">{t('menu.dashboard')}</h6>
      {/* Page content */}
    </div>
  )
}

export default Dashboard
```

### Breadcrumb Navigation

**Dynamic Breadcrumbs with i18n**:
```javascript
import { useLocation, Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

function Breadcrumb() {
  const location = useLocation()
  const { t } = useTranslation()
  const pathSegments = location.pathname.split('/').filter(Boolean)

  return (
    <ul className="d-flex align-items-center gap-2">
      <li>
        <Link to="/">
          <iconify-icon icon="solar:home-smile-angle-outline" />
          {t('menu.home')}
        </Link>
      </li>
      {pathSegments.map((segment, index) => (
        <li key={index}>
          - {t(`menu.${segment}`)}
        </li>
      ))}
    </ul>
  )
}
```

## Build Process

### Development Mode

**Frontend Development Server** (Port 3000):
```bash
cd frontend
npm run dev
```
- Vite dev server with hot reload
- Proxies API calls to GoFrame server (port 8080)
- Fast refresh for React components

**Backend Server** (Port 8080):
```bash
go run main.go
```
- Serves API endpoints with CORS enabled
- Serves static assets from public/ (copied from assets/)

**Workflow**:
1. Developer edits React components
2. Vite hot-reloads changes instantly
3. API calls proxy to GoFrame backend
4. No build step required during development

### Production Build

**Frontend Build**:
```bash
cd frontend
npm run build
```

**Vite Build Process**:
1. Bundles React components
2. Minifies JavaScript
3. Optimizes assets
4. Outputs to `../public/` directory:
   - `index.html` (entry point with injected script tags)
   - `assets/main-*.js` (bundled React app)

**Backend Build**:
```bash
go build -o bin/tzlev main.go
```

**Deployment**:
- Single binary (`bin/tzlev`)
- Requires `public/` directory alongside binary (with css/, js/, images/ subdirectories)
- No separate frontend server needed
- CORS enabled for API endpoints

## Integration Points

### 1. Template JavaScript Integration

**Challenge**: Wowdash template requires jQuery and custom JavaScript for features like:
- Sidebar toggle
- Theme switching
- Dropdown menus
- Mobile navigation

**Solution**:
```javascript
// In App.jsx
useEffect(() => {
  const script = document.createElement('script')
  script.src = '/assets/js/app.js'
  script.async = false
  document.body.appendChild(script)
  return () => document.body.removeChild(script)
}, [])
```

**Why This Works**:
1. jQuery and libraries load from HTML (before React)
2. React initializes
3. App component mounts
4. `useEffect` runs and loads app.js
5. Template features initialize with correct DOM

### 2. API Communication

**Pattern**: React components call GoFrame API

**Example**:
```javascript
// In React component
fetch('/api/health')
  .then(res => res.json())
  .then(data => setHealth(data))
```

**GoFrame Handler**:
```go
group.GET("/health", func(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "status": "ok",
        "message": "Server is running",
    })
})
```

**Benefits**:
- No CORS configuration needed (same origin)
- Simple fetch API usage
- Type-safe Go API handlers

### 3. Static Asset Serving

**Structure**:
```
public/
├── css/          # Template stylesheets (copied from assets/)
├── js/           # Template JavaScript libraries (copied from assets/)
└── images/       # Images, icons, logos (copied from assets/)
```

**GoFrame Configuration**:
```go
s.SetServerRoot("public")  // Serves all files from public directory
```

**Usage in React**:

```jsx
<img src="/assets/images/logo.png" alt="logo"/>
<link rel="stylesheet" href="/assets/css/style.css"/>
```

## Extension Points

### Adding New Features

#### 1. Add a New Page with Route

**Step 1: Create the page component**
```jsx
// Create: frontend/src/pages/Settings.jsx
import React from 'react'
import { useTranslation } from 'react-i18next'

function Settings() {
  const { t } = useTranslation()

  return (
    <div>
      <h6 className="fw-semibold mb-0">{t('menu.settings')}</h6>
      {/* Page content */}
    </div>
  )
}

export default Settings
```

**Step 2: Add route to App.jsx**
```jsx
import Settings from './pages/Settings'

// In Routes component
<Route path="/settings" element={<Settings />} />
```

**Step 3: Add to menu configuration**
```javascript
// In frontend/src/config/menuConfig.js
{
  id: 'settings',
  titleKey: 'menu.settings',
  icon: 'icon-park-outline:setting-two',
  path: '/settings',
}
```

**Step 4: Add translations**
```json
// frontend/src/locales/en/translation.json
{
  "menu": {
    "settings": "Settings"
  }
}

// frontend/src/locales/he/translation.json
{
  "menu": {
    "settings": "הגדרות"
  }
}
```

#### 2. Add a New Menu Item with Submenu

```javascript
// In frontend/src/config/menuConfig.js
{
  id: 'reports',
  titleKey: 'menu.reports',
  icon: 'solar:document-text-outline',
  path: '/reports',
  children: [
    {
      id: 'reports-sales',
      titleKey: 'menu.salesReports',
      path: '/reports/sales',
    },
    {
      id: 'reports-analytics',
      titleKey: 'menu.analyticsReports',
      path: '/reports/analytics',
    }
  ]
}
```

**Create corresponding pages and routes**:
```jsx
// frontend/src/pages/reports/Sales.jsx
// frontend/src/pages/reports/Analytics.jsx

// Add routes in App.jsx
<Route path="/reports/sales" element={<SalesReports />} />
<Route path="/reports/analytics" element={<AnalyticsReports />} />
```

**Add translations**:
```json
{
  "menu": {
    "reports": "Reports",
    "salesReports": "Sales Reports",
    "analyticsReports": "Analytics Reports"
  }
}
```

#### 3. Add New Translations

**Add to English** (`frontend/src/locales/en/translation.json`):
```json
{
  "pages": {
    "users": {
      "title": "User Management",
      "addButton": "Add User",
      "deleteButton": "Delete"
    }
  }
}
```

**Add to Hebrew** (`frontend/src/locales/he/translation.json`):
```json
{
  "pages": {
    "users": {
      "title": "ניהול משתמשים",
      "addButton": "הוסף משתמש",
      "deleteButton": "מחק"
    }
  }
}
```

**Use in component**:
```jsx
const { t } = useTranslation()
<h1>{t('pages.users.title')}</h1>
<button>{t('pages.users.addButton')}</button>
```

#### 4. Add a New Language

**Step 1: Create translation file**
```javascript
// frontend/src/locales/es/translation.json
{
  "menu": {
    "home": "Inicio",
    "dashboard": "Tablero"
  }
}
```

**Step 2: Update i18n config**
```javascript
// frontend/src/config/i18n.js
import es from '../locales/es/translation.json'

i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    he: { translation: he },
    es: { translation: es }  // Add new language
  }
})
```

**Step 3: Add to language switcher**
```jsx
<button onClick={() => changeLanguage('es')}>Español</button>
```

#### 5. New React Component with i18n

```jsx
// Create: frontend/src/components/NewFeature.jsx
import React from 'react'
import { useTranslation } from 'react-i18next'

function NewFeature() {
  const { t } = useTranslation()

  return (
    <div>
      <h2>{t('features.newFeature.title')}</h2>
      <p>{t('features.newFeature.description')}</p>
    </div>
  )
}

export default NewFeature
```

#### 6. New API Endpoint
```go
// In main.go or internal/controller/
s.Group("/api", func(group *ghttp.RouterGroup) {
    group.GET("/users", func(r *ghttp.Request) {
        r.Response.WriteJson(g.Map{
            "users": []string{"user1", "user2"},
        })
    })
})
// Note: CORS is already enabled globally
```

```javascript
// In React component
fetch('/api/users')
  .then(res => res.json())
  .then(data => setUsers(data.users))
```

#### 7. Add New CLI Command

**Step 1: Create command handler**
```go
// Create: internal/cli/export.go
package cli

import (
    "context"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gcmd"
    "os"
)

func ExportData(ctx context.Context, parser *gcmd.Parser) {
    format := parser.GetOpt("format", "csv").String()
    output := parser.GetOpt("output", "data.csv").String()

    g.Log().Infof(ctx, "Exporting data to %s (format: %s)", output, format)

    // Export logic here
    // Use service layer for data access

    g.Log().Info(ctx, "Export completed successfully")
    os.Exit(0)
}
```

**Step 2: Register command in main.go**
```go
func runCLI(ctx context.Context) {
    parser, err := gcmd.Parse(nil)
    if err != nil {
        g.Log().Fatal(ctx, "Failed to parse CLI args:", err)
    }

    command := parser.GetArg(1).String()

    switch command {
    case "export":
        cli.ExportData(ctx, parser)  // Add new command
    case "migrate":
        cli.RunMigrate(ctx, parser)
    // ... other commands
    }
}
```

**Step 3: Update help text**
```go
// In internal/cli/help.go
func ShowHelp(ctx context.Context) {
    helpText := `
COMMANDS:
    export               Export data to file
    migrate              Run database migrations
    ...
EXAMPLES:
    ./tzlev export --format=csv --output=users.csv
`
    fmt.Println(helpText)
}
```

**Step 4: Test the command**
```bash
go run main.go export --format=csv --output=data.csv
./tzlev export --help
```

#### 8. New Page/Route

**Option A: Client-side routing with React Router**
```bash
cd frontend
npm install react-router-dom
```

```jsx
// In App.jsx
import { BrowserRouter, Routes, Route } from 'react-router-dom'

<BrowserRouter>
  <Routes>
    <Route path="/" element={<MainContent />} />
    <Route path="/dashboard" element={<Dashboard />} />
  </Routes>
</BrowserRouter>
```

**Option B: Server-side routing**
```go
// In main.go
s.BindHandler("/admin", func(r *ghttp.Request) {
    r.Response.ServeFile("public/admin.html")
})
```

### Recommended Patterns

#### State Management
- Start with React hooks (useState, useContext)
- Add Redux/Zustand for complex state
- Keep API calls in custom hooks

#### API Layer
- Create dedicated API service module:
```javascript
// frontend/src/services/api.js
export const api = {
  health: () => fetch('/api/health').then(r => r.json()),
  getUsers: () => fetch('/api/users').then(r => r.json()),
}
```

#### Error Handling
- Use error boundaries in React
- Return consistent error responses from Go
```go
r.Response.WriteJson(g.Map{
    "error": true,
    "message": "Something went wrong",
})
```

#### 8. Add New Database Model

**Step 1: Create model**
```go
// Create: internal/model/product.go
package model

import (
    "github.com/gogf/gf/v2/os/gtime"
)

type Product struct {
    Id          uint        `json:"id"`
    CreatedAt   *gtime.Time `json:"created_at"`
    UpdatedAt   *gtime.Time `json:"updated_at"`
    DeletedAt   *gtime.Time `json:"deleted_at,omitempty"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Price       float64     `json:"price"`
    Stock       int         `json:"stock"`
}
```

**Step 2: Create migration**
```sql
-- Create: migrations/000003_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER DEFAULT 0
);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);
```

```sql
-- Create: migrations/000003_create_products_table.down.sql
DROP TABLE IF EXISTS products;
```

**Step 3: Create repository**
```go
// Create: internal/repository/product_repository.go
package repository

import (
    "context"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
    "tzlev/internal/model"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
    return &ProductRepository{}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
    product.CreatedAt = gtime.Now()
    product.UpdatedAt = gtime.Now()

    _, err := g.DB().Ctx(ctx).Insert("products", product)
    return err
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
    var product model.Product
    err := g.DB().Ctx(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        Scan(&product)

    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *ProductRepository) List(ctx context.Context, offset, limit int) ([]model.Product, error) {
    var products []model.Product
    err := g.DB().Ctx(ctx).
        Where("deleted_at IS NULL").
        Offset(offset).
        Limit(limit).
        Scan(&products)

    return products, err
}

func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
    product.UpdatedAt = gtime.Now()

    _, err := g.DB().Ctx(ctx).
        Where("id = ?", product.Id).
        Update("products", product)
    return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
    // Soft delete
    _, err := g.DB().Ctx(ctx).
        Where("id = ?", id).
        Update("products", g.Map{"deleted_at": gtime.Now()})
    return err
}
```

**Step 4: Create service**
```go
// Create: internal/service/product_service.go
package service

import (
    "context"
    "fmt"
    "time"
    "tzlev/internal/cache"
    "tzlev/internal/model"
    "tzlev/internal/repository"
)

type ProductService struct {
    productRepo  *repository.ProductRepository
    cacheManager *cache.CacheManager
}

func NewProductService() *ProductService {
    return &ProductService{
        productRepo:  repository.NewProductRepository(),
        cacheManager: cache.NewCacheManager(),
    }
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
    cacheKey := fmt.Sprintf("product:%d", id)

    // Try cache first
    var product model.Product
    err := s.cacheManager.Get(ctx, cacheKey, &product)
    if err == nil {
        return &product, nil
    }

    // Get from database
    dbProduct, err := s.productRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cache for 10 minutes
    _ = s.cacheManager.Set(ctx, cacheKey, dbProduct, 10*time.Minute)

    return dbProduct, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) error {
    return s.productRepo.Create(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) error {
    if err := s.productRepo.Update(ctx, product); err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := fmt.Sprintf("product:%d", product.ID)
    _ = s.cacheManager.Delete(ctx, cacheKey)

    return nil
}
```

**Step 5: Create controller**
```go
// Create: internal/controller/product_controller.go
package controller

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "tzlev/internal/model"
    "tzlev/internal/service"
)

type ProductController struct {
    productService *service.ProductService
}

func NewProductController() *ProductController {
    return &ProductController{
        productService: service.NewProductService(),
    }
}

func (c *ProductController) GetProduct(r *ghttp.Request) {
    ctx := r.Context()
    id := r.Get("id").Uint()

    product, err := c.productService.GetProductByID(ctx, id)
    if err != nil {
        r.Response.Status = 404
        r.Response.WriteJson(g.Map{"error": "Product not found"})
        return
    }

    r.Response.WriteJson(g.Map{"product": product})
}

func (c *ProductController) CreateProduct(r *ghttp.Request) {
    ctx := r.Context()

    var product model.Product
    if err := r.Parse(&product); err != nil {
        r.Response.Status = 400
        r.Response.WriteJson(g.Map{"error": "Invalid request"})
        return
    }

    if err := c.productService.CreateProduct(ctx, &product); err != nil {
        r.Response.Status = 500
        r.Response.WriteJson(g.Map{"error": "Failed to create product"})
        return
    }

    r.Response.WriteJson(g.Map{"product": product})
}
```

**Step 6: Register routes in main.go**
```go
// In main.go
func setupRoutes(s *ghttp.Server) {
    productCtrl := controller.NewProductController()

    s.Group("/api", func(group *ghttp.RouterGroup) {
        group.Middleware(middleware.Auth())

        group.GET("/products/:id", productCtrl.GetProduct)
        group.POST("/products", productCtrl.CreateProduct)
    })
}
```

**Step 7: Run migration**
```bash
go run main.go migrate up
```

#### 9. Add New Email Template

**Step 1: Create HTML template**
```html
<!-- Create: templates/email/order_confirmation.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .header { background-color: #4F46E5; color: white; padding: 20px; }
        .content { padding: 30px; }
        .footer { text-align: center; color: #6b7280; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Order Confirmation</h1>
    </div>
    <div class="content">
        <p>Hello {{.Name}},</p>
        <p>Your order #{{.OrderID}} has been confirmed.</p>
        <p>Total: ${{.Total}}</p>
    </div>
    <div class="footer">
        <p>&copy; 2024 Tzlev. All rights reserved.</p>
    </div>
</body>
</html>
```

**Step 2: Add email method to service**
```go
// Add to: internal/email/email.go
type OrderConfirmationData struct {
    Name    string
    OrderID string
    Total   float64
}

func (es *EmailService) SendOrderConfirmation(to, name, orderID string, total float64) error {
    tmplPath := filepath.Join(es.templatePath, "order_confirmation.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    var body bytes.Buffer
    data := OrderConfirmationData{
        Name:    name,
        OrderID: orderID,
        Total:   total,
    }

    if err := tmpl.Execute(&body, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    return es.Send(&EmailData{
        To:      []string{to},
        Subject: "Order Confirmation",
        Body:    body.String(),
    })
}
```

**Step 3: Use in service**
```go
// In order service
func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) error {
    if err := s.orderRepo.Create(ctx, order); err != nil {
        return err
    }

    // Send confirmation email
    go func() {
        user, _ := s.userRepo.FindByID(ctx, order.UserID)
        _ = s.emailService.SendOrderConfirmation(
            user.Email,
            user.Name,
            order.ID,
            order.Total,
        )
    }()

    return nil
}
```

#### 10. Add Configuration Option

**Step 1: Add to config.yaml**
```yaml
# Add new feature configuration
features:
  enableNotifications: true
  maxUploadSize: "10MB"
  allowedFileTypes:
    - "image/png"
    - "image/jpeg"
    - "application/pdf"
```

**Step 2: Access in code**
```go
enableNotifications := g.Cfg().MustGet(ctx, "features.enableNotifications").Bool()
maxUploadSize := g.Cfg().MustGet(ctx, "features.maxUploadSize").String()
allowedTypes := g.Cfg().MustGet(ctx, "features.allowedFileTypes").Strings()
```

## Technology Stack

### Backend
- **GoFrame v2.9.4**: Web framework with CORS support and built-in ORM
- **Go 1.24.2**: Programming language
- **Cobra**: CLI framework for command-line interface

### Database & ORM
- **PostgreSQL 16+**: Primary database (ACID-compliant, relational)
- **GoFrame gdb**: Built-in ORM for database operations
- **golang-migrate**: Database migration management (optional)

### Caching & Sessions
- **Redis 7+**: In-memory data store for sessions and caching
- **go-redis/redis**: Official Redis client for Go

### Authentication & Security
- **Google OAuth 2.0**: User authentication via Google
- **golang.org/x/oauth2**: OAuth 2.0 client library
- **JWT (golang-jwt/jwt)**: JSON Web Tokens for API authentication
- **bcrypt**: Password hashing (if needed)

### Email
- **SMTP (Gmail)**: Email delivery service
- **net/smtp**: Go standard library for SMTP
- **html/template**: Go template engine for email rendering

### Frontend Core
- **React 19**: UI library
- **Vite 7**: Build tool and dev server
- **React Router v6**: Client-side routing and navigation

### Internationalization
- **i18next**: Internationalization framework
- **react-i18next**: React bindings for i18next
- **Supported Languages**:
  - English (en) - LTR
  - Hebrew (he) - RTL

### UI & Styling
- **Bootstrap 5**: CSS framework with RTL support
- **jQuery 3.7.1**: Template compatibility
- **Iconify**: Icon library

### Template
- **Wowdash**: Bootstrap 5 Admin Dashboard Template

### Development Tools
- **Vite HMR**: Hot module replacement for React development
- **Cobra CLI**: Command-line interface framework for admin tasks

## Development Workflow

### Initial Setup

**Prerequisites**:
- Go 1.24.2 or higher
- Node.js 16+ and npm
- PostgreSQL 16+ (installed and running)
- Redis 7+ (installed and running)
- Google OAuth credentials (Client ID and Secret)
- Gmail account with App Password

**Setup Steps**:

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd tzlev
   ```

2. **Install PostgreSQL** (if not installed)
   ```bash
   # macOS
   brew install postgresql@16
   brew services start postgresql@16

   # Ubuntu/Debian
   sudo apt install postgresql-16
   sudo systemctl start postgresql

   # Create database and user
   psql postgres
   CREATE DATABASE tzlev_db;
   CREATE USER tzlev_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE tzlev_db TO tzlev_user;
   \q
   ```

3. **Install Redis** (if not installed)
   ```bash
   # macOS
   brew install redis
   brew services start redis

   # Ubuntu/Debian
   sudo apt install redis-server
   sudo systemctl start redis
   ```

4. **Set up Google OAuth**
   - Go to [Google Cloud Console](https://console.cloud.google.com)
   - Create a new project or select existing one
   - Enable Google+ API
   - Create OAuth 2.0 credentials (Web application)
   - Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
   - Copy Client ID and Client Secret

5. **Set up Gmail App Password**
   - Go to your Google Account settings
   - Navigate to Security > 2-Step Verification
   - Scroll down and select "App passwords"
   - Generate a new app password for "Mail"
   - Copy the generated password

6. **Configure environment variables**
   ```bash
   # Copy example env file
   cp .env.example .env

   # Edit .env with your values
   nano .env
   ```

   Required environment variables:
   ```bash
   # Database
   DB_PASSWORD=your_database_password

   # Redis
   REDIS_PASSWORD=                # Leave empty if no password

   # Google OAuth
   GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=your_google_client_secret

   # Email
   EMAIL_USERNAME=your_email@gmail.com
   EMAIL_PASSWORD=your_gmail_app_password
   EMAIL_FROM=noreply@tzlev.com

   # Session & Security
   SESSION_SECRET=your_random_session_secret_key
   JWT_SECRET=your_random_jwt_secret_key
   ```

7. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

   This will install:
   - GoFrame v2 (includes gdb ORM)
   - Cobra CLI framework
   - go-redis
   - oauth2
   - JWT library

8. **Install Node dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

   This installs:
   - React, react-dom, react-router-dom
   - i18next, react-i18next
   - Vite

9. **Copy template assets**
   ```bash
   cp -r assets/* public/
   ```

   This copies css/, js/, images/ directories directly into public/

10. **Run database migrations**
    ```bash
    go run main.go migrate up
    ```

11. **Seed database (optional)**
    ```bash
    go run main.go seed
    ```

12. **Verify setup**
    - Check PostgreSQL: `psql -U tzlev_user -d tzlev_db`
    - Check Redis: `redis-cli ping` (should return PONG)
    - Check config: `cat config/config.yaml`

### Daily Development

**Web Mode (Frontend Development)**:
1. Start GoFrame server: `go run main.go` (no arguments → web mode)
2. Start Vite dev server: `cd frontend && npm run dev`
3. Visit http://localhost:3000
4. Edit React components with hot reload
5. Edit Go code and restart server

**CLI Mode (Testing CLI Commands)**:
```bash
# Test CLI commands during development
go run main.go help
go run main.go version
go run main.go user:list

# Test with flags
go run main.go migrate --up
go run main.go seed --table=users
```

### Before Deployment
1. Build frontend: `cd frontend && npm run build`
2. Test production build (web mode): `go run main.go` → http://localhost:8080
3. Test CLI commands: `go run main.go help`
4. Build binary: `go build -o bin/tzlev main.go`
5. Deploy binary + public/ directory

### Running in Production

**Web Mode** (Start HTTP Server):
```bash
./bin/tzlev          # No arguments → web server mode
# or
./tzlev
```

**CLI Mode** (Administrative Tasks):
```bash
./tzlev migrate                               # Run migrations
./tzlev seed                                  # Seed database
./tzlev user:create --email=admin@test.com   # Create user
./tzlev user:list                            # List users
./tzlev version                              # Show version
./tzlev help                                 # Show help
```

## Configuration

### Vite Configuration (frontend/vite.config.js)
```javascript
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../public',           // Build to Go's public directory
    emptyOutDir: false,            // Don't delete existing assets
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // Proxy API to Go server
        changeOrigin: true
      }
    }
  }
})
```

### GoFrame Configuration
- Port: 8080 (configurable)
- Static root: `public/` (serves all files from public directory)
- Template assets: Copied from `assets/` to `public/` (css/, js/, images/)
- API prefix: `/api`
- CORS: Enabled globally via `CORSDefault()` middleware

## Security Considerations

### Backend
- **CORS**: Enabled for API endpoints - configure allowed origins in production
- Input validation on all API endpoints
- Use GoFrame's built-in security features
- Implement authentication/authorization as needed
- Rate limiting for API endpoints
- In production, restrict CORS to specific domains

### Frontend
- Sanitize user input
- Use HTTPS in production
- Implement CSRF protection
- Secure API communication

## Performance Optimization

### Frontend
- Code splitting with React lazy loading
- Image optimization
- Bundle size monitoring
- Tree shaking (automatic with Vite)

### Backend
- Static file caching
- Gzip compression
- Database connection pooling
- CDN for static assets (optional)

## Testing Strategy

### Frontend Testing
```bash
npm install -D vitest @testing-library/react
```
- Unit tests for components
- Integration tests for API calls
- E2E tests with Playwright/Cypress

### Backend Testing
```go
go test ./...
```
- Unit tests for handlers
- Integration tests for API endpoints
- Load testing with tools like k6

## Future Enhancements

### Potential Additions
1. **Authentication**: JWT-based auth system with i18n support
2. **Database**: Add PostgreSQL/MySQL with MCP
3. **WebSockets**: Real-time features with GoFrame WebSocket support
4. **Docker**: Containerization for easy deployment
5. **CI/CD**: GitHub Actions for automated builds
6. **API Documentation**: Swagger/OpenAPI integration
7. **Monitoring**: Prometheus metrics, logging
8. **Additional Languages**: Extend i18n to support Arabic, French, Spanish, etc.
9. **PWA Support**: Progressive Web App capabilities
10. **Theme Persistence**: Save user's theme and language preferences

### Scalability
- Horizontal scaling with load balancer
- Redis for session management
- Message queue for async tasks
- Microservices architecture (if needed)

## Troubleshooting

### Common Issues

**Issue**: React not loading
- Check: `public/index.html` exists after build
- Check: GoFrame server is serving correct path

**Issue**: API calls fail
- Check: GoFrame server is running on port 8080
- Check: API routes are registered correctly
- Check: CORS middleware is enabled (required for external API calls)

**Issue**: Assets not found (CSS/JS/Images)
- Check: Assets copied from `assets/` to `public/`
- Run: `cp -r assets/* public/`
- Check: GoFrame serving from `public/` directory
- Verify files exist: `public/css/`, `public/js/`, `public/images/`

**Issue**: Template features not working
- Check: jQuery loaded before React
- Check: app.js loaded after React mounts
- Check: Assets in correct public/ directory structure

## Conclusion

This architecture provides a comprehensive, production-ready foundation for building a modern, multilingual, full-stack web application with Go and React. The **dual-mode operation** (web + CLI) provides flexibility for both production runtime and administrative tasks from a single binary. The embedded frontend pattern simplifies deployment while maintaining development flexibility. The integration with i18next enables seamless Hebrew and English support with automatic RTL/LTR switching. React Router provides a smooth SPA experience with dynamic, configurable navigation. The Wowdash template integration provides a professional UI out of the box that works seamlessly with both LTR and RTL layouts.

### External Services Integration

The architecture integrates with industry-standard external services to provide enterprise-grade capabilities:

- **PostgreSQL**: Reliable, ACID-compliant relational database with GoFrame gdb ORM for type-safe data access
- **Redis**: High-performance in-memory store for sessions and caching, reducing database load
- **Google OAuth 2.0**: Secure user authentication leveraging Google's identity platform
- **Gmail SMTP**: Transactional email delivery with HTML template support

### Configuration Management

- **Environment-based configuration**: YAML-based config with environment-specific overrides
- **Secrets management**: Sensitive data stored in environment variables
- **Flexible deployment**: Same codebase works in development, staging, and production

### Data Layer

- **Repository pattern**: Clean abstraction over database operations
- **Service layer**: Business logic with caching, email, and transaction support
- **Migration management**: Version-controlled database schema changes via CLI
- **Seeding support**: Database initialization for development and testing

### Authentication & Security

- **OAuth 2.0 flow**: Industry-standard authentication with Google
- **Redis-backed sessions**: Distributed session storage with TTL expiration
- **JWT support**: Token-based authentication for API clients
- **Middleware architecture**: Reusable authentication and authorization logic

### Communication

- **Template-based emails**: HTML email templates with Go html/template
- **Async email sending**: Non-blocking email delivery
- **Multiple templates**: Welcome, password reset, notifications, and custom templates

### Key Advantages

**Architecture**:
- **Dual-mode operation**: Single binary for web server and CLI commands
- **Administrative CLI**: Database migrations, seeding, user management without web interface
- Single deployment artifact (binary + public/ directory)
- Fast development cycle with hot reload
- Type-safe backend with Go and GoFrame gdb
- Modern, reactive frontend with React 19

**Internationalization**:
- **Multilingual support** with Hebrew (RTL) and English (LTR)
- Automatic RTL/LTR switching based on language
- Centralized translation management
- Easy to add new languages

**Frontend**:
- **Client-side routing** for SPA experience
- **Dynamic menu configuration** with i18n
- Professional UI with Wowdash Bootstrap 5 template
- Component-based architecture with React

**Backend**:
- CORS enabled for API flexibility
- RESTful API with consistent error handling
- Database connection pooling
- Redis caching for performance
- Health checks for all external services

**Data & Storage**:
- PostgreSQL for persistent data with ACID guarantees
- Redis for sessions, caching, and temporary data
- Repository pattern for data access abstraction
- Migration-based schema management

**Authentication & Security**:
- Google OAuth integration for user login
- Redis-backed session management
- JWT for API authentication
- Middleware-based access control

**Communication**:
- Email service with HTML templates
- Async email delivery
- Template-based notifications

**Development**:
- Clear separation of concerns (controller → service → repository)
- Shared service layer between web and CLI modes
- Configuration-driven architecture
- Easy to extend and maintain
- Scalable for additional features

**Deployment**:
- Environment-based configuration
- Secrets management via environment variables
- Single binary deployment
- Health checks and monitoring-ready

This architecture is production-ready and provides all the necessary components for building a modern, scalable, multilingual web application with authentication, data persistence, caching, and email capabilities.
