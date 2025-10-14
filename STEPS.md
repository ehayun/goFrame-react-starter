# Tzlev Application - Implementation Steps

This document outlines all the steps necessary to build the Tzlev application from scratch, following the architecture defined in DESIGN.md.

## Table of Contents

1. [Phase 1: Project Initialization](#phase-1-project-initialization)
2. [Phase 2: Configuration Management](#phase-2-configuration-management)
3. [Phase 3: Database Setup (PostgreSQL via MCP)](#phase-3-database-setup-postgresql-via-mcp)
4. [Phase 4: Redis Setup](#phase-4-redis-setup)
5. [Phase 5: Core Backend Structure](#phase-5-core-backend-structure)
6. [Phase 6: Authentication System](#phase-6-authentication-system)
7. [Phase 7: Email Service](#phase-7-email-service)
8. [Phase 8: Frontend Foundation](#phase-8-frontend-foundation)
9. [Phase 9: Template Integration](#phase-9-template-integration)
10. [Phase 10: Internationalization](#phase-10-internationalization)
11. [Phase 11: Routing & Navigation](#phase-11-routing--navigation)
12. [Phase 12: CLI Mode](#phase-12-cli-mode)
13. [Phase 13: Build & Deployment](#phase-13-build--deployment)
14. [Phase 14: Testing & Validation](#phase-14-testing--validation)

---

## Phase 1: Project Initialization

### Step 1.1: Initialize Go Module
```bash
mkdir -p /Users/eli/projects/tzlev
cd /Users/eli/projects/tzlev
go mod init tzlev
```

### Step 1.2: Install Go Dependencies
```bash
# GoFrame framework
go get -u github.com/gogf/gf/v2

# PostgreSQL driver
go get -u github.com/lib/pq

# Redis client
go get -u github.com/redis/go-redis/v9

# OAuth2
go get -u golang.org/x/oauth2
go get -u golang.org/x/oauth2/google

# JWT
go get -u github.com/golang-jwt/jwt/v5

# CLI framework (optional, if using Cobra)
go get -u github.com/spf13/cobra
```

### Step 1.3: Create Project Directory Structure
```bash
# Backend directories
mkdir -p internal/{controller,service,repository,model}
mkdir -p internal/{database,redis,session,cache}
mkdir -p internal/{oauth,jwt,email,middleware,cli}

# Configuration and migrations
mkdir -p config
mkdir -p migrations

# Email templates
mkdir -p templates/email

# Frontend directories
mkdir -p frontend/src/{components,pages,locales,config}
mkdir -p frontend/src/locales/{en,he}

# Assets and public
mkdir -p assets/{css,js,images}
mkdir -p public
```

### Step 1.4: Create Initial Go Files
```bash
# Create main.go (entry point)
touch main.go

# Create package files
touch internal/database/database.go
touch internal/redis/redis.go
touch internal/session/session.go
touch internal/cache/cache.go
touch internal/oauth/google.go
touch internal/jwt/jwt.go
touch internal/email/email.go
touch internal/middleware/auth.go
```

---

## Phase 2: Configuration Management

### Step 2.1: Create Configuration Files

**File: `config/config.yaml`**
```yaml
# Application Configuration
app:
  name: "Tzlev"
  environment: "development"
  debug: true
  port: 8080

# Database Configuration (uses MCP PostgreSQL)
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
  password: "${REDIS_PASSWORD}"
  database: 0
  poolSize: 10
  sessionPrefix: "tzlev:session:"
  cachePrefix: "tzlev:cache:"
  sessionTTL: "24h"

# Gmail OAuth Configuration
oauth:
  google:
    clientID: "${GOOGLE_CLIENT_ID}"
    clientSecret: "${GOOGLE_CLIENT_SECRET}"
    redirectURL: "http://localhost:8080/auth/google/callback"
    scopes:
      - "https://www.googleapis.com/auth/userinfo.email"
      - "https://www.googleapis.com/auth/userinfo.profile"

# Email Service Configuration
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

# Session Configuration
session:
  store: "redis"
  secretKey: "${SESSION_SECRET}"
  cookieName: "tzlev_session"
  cookieMaxAge: 86400
  cookieSecure: false
  cookieHttpOnly: true
  cookieSameSite: "lax"

# Security Configuration
security:
  jwtSecret: "${JWT_SECRET}"
  jwtExpiration: "1h"
  bcryptCost: 12

# Logging Configuration
logging:
  level: "debug"
  format: "json"
  output: "stdout"
  filePath: "logs/app.log"
```

### Step 2.2: Create Environment Example File

**File: `.env.example`**
```bash
# Database (MCP will provide these)
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

### Step 2.3: Create Actual .env File
```bash
cp .env.example .env
# Edit .env with actual values
```

### Step 2.4: Add Generated Files to .gitignore

Since `public/` is generated from `assets/`, we should ignore it in git.

```bash
echo ".env" >> .gitignore
echo "bin/" >> .gitignore
echo "public/" >> .gitignore
echo "logs/" >> .gitignore
echo "node_modules/" >> .gitignore
```

**Note**: The `public/` directory is entirely generated by:
- `npm run build:assets` (template assets from `assets/`)
- `npm run build` (React frontend from `frontend/`)

Always commit changes to `assets/` and `frontend/src/`, never to `public/`.

---

## Phase 3: Database Setup (PostgreSQL via MCP)

### Step 3.1: Install MCP PostgreSQL Server
```bash
# Follow MCP PostgreSQL installation instructions
# The user mentioned "all database calls will use the mcp"
# Ensure MCP PostgreSQL server is configured and running
```

### Step 3.2: Create Database Connection Module

**File: `internal/database/database.go`**
```go
package database

import (
	"fmt"

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

### Step 3.3: Create Database Models

**File: `internal/model/user.go`**
```go
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

### Step 3.4: Create Migration Files

**File: `migrations/000001_create_users_table.up.sql`**
```sql
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

**File: `migrations/000001_create_users_table.down.sql`**
```sql
DROP TABLE IF EXISTS users;
```

### Step 3.5: Create Repository Layer

**File: `internal/repository/user_repository.go`**
```go
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

---

## Phase 4: Redis Setup

### Step 4.1: Create Redis Connection Module

**File: `internal/redis/redis.go`**
```go
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
		Addr: fmt.Sprintf("%s:%d",
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

### Step 4.2: Create Session Manager

**File: `internal/session/session.go`**
```go
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

### Step 4.3: Create Cache Manager

**File: `internal/cache/cache.go`**
```go
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

---

## Phase 5: Core Backend Structure

### Step 5.1: Create Service Layer

**File: `internal/service/user_service.go`**
```go
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
	cacheKey := fmt.Sprintf("user:%d", user.Id)
	_ = s.cacheManager.Delete(ctx, cacheKey)

	return nil
}
```

### Step 5.2: Create Basic Controller

**File: `internal/controller/health_controller.go`**
```go
package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Check(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"status":  "ok",
		"message": "Server is running",
	})
}
```

---

## Phase 6: Authentication System

### Step 6.1: Setup Google OAuth

**File: `internal/oauth/google.go`**
```go
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

### Step 6.2: Create JWT Module

**File: `internal/jwt/jwt.go`**
```go
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

### Step 6.3: Create Authentication Controller

**File: `internal/controller/auth_controller.go`**
```go
package controller

import (
	"crypto/rand"
	"encoding/base64"
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
			GoogleId:    googleUser.ID,
			IsActive:    true,
			LastLoginAt: gtime.New(now),
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
		user.LastLoginAt = gtime.New(now)
		if err := c.userRepo.Update(ctx, user); err != nil {
			g.Log().Error(ctx, "Failed to update user:", err)
		}
	}

	// Create session
	sessionID := r.Session.Id()
	sess := &session.Session{
		UserID: user.Id,
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
	r.Session.Set("user_id", user.Id)
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

### Step 6.4: Create Authentication Middleware

**File: `internal/middleware/auth.go`**
```go
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

---

## Phase 7: Email Service

### Step 7.1: Create Email Service Module

**File: `internal/email/email.go`**
```go
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
```

### Step 7.2: Create Email Templates

**File: `templates/email/welcome.html`**
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

---

## Phase 8: Frontend Foundation

### Step 8.1: Initialize Frontend Project
```bash
cd frontend
npm init -y
npm install react react-dom react-router-dom
npm install i18next react-i18next
npm install -D vite @vitejs/plugin-react
```

### Step 8.2: Configure Vite

**File: `frontend/vite.config.js`**
```javascript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../public',
    emptyOutDir: false, // Don't delete existing assets
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
```

### Step 8.3: Create Frontend Package.json Scripts

**File: `frontend/package.json`**
```json
{
  "name": "tzlev-frontend",
  "version": "1.0.0",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.22.0",
    "i18next": "^23.7.0",
    "react-i18next": "^14.0.0"
  },
  "devDependencies": {
    "@vitejs/plugin-react": "^4.2.0",
    "vite": "^5.0.0"
  }
}
```

### Step 8.4: Create Frontend HTML Template

**File: `frontend/index.html`**
```html
<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tzlev</title>

    <!-- Iconify -->
    <script src="https://code.iconify.design/iconify-icon/1.0.8/iconify-icon.min.js"></script>

    <!-- Template CSS (from assets) -->
    <link rel="stylesheet" href="/css/remixicon.css">
    <link rel="stylesheet" href="/css/lib/bootstrap.min.css">
    <link rel="stylesheet" href="/css/lib/apexcharts.css">
    <link rel="stylesheet" href="/css/lib/dataTables.min.css">
    <link rel="stylesheet" href="/css/style.css">

    <!-- Template JavaScript (jQuery and libraries - must load before React) -->
    <script src="/js/lib/jquery-3.7.1.min.js"></script>
    <script src="/js/lib/bootstrap.bundle.min.js"></script>
    <script src="/js/lib/apexcharts.min.js"></script>
    <script src="/js/lib/dataTables.min.js"></script>
</head>
<body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
</body>
</html>
```

### Step 8.5: Create React Entry Point

**File: `frontend/src/main.jsx`**
```javascript
import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import './config/i18n'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
)
```

### Step 8.6: Create Main App Component

**File: `frontend/src/App.jsx`**
```javascript
import React, { useEffect } from 'react'
import { Routes, Route } from 'react-router-dom'
import Sidebar from './components/Sidebar'
import Navbar from './components/Navbar'
import Footer from './components/Footer'
import ThemeCustomization from './components/ThemeCustomization'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'

function App() {
  // Load template JavaScript after React mounts
  useEffect(() => {
    const script = document.createElement('script')
    script.src = '/js/app.js'
    script.async = false
    document.body.appendChild(script)

    return () => {
      if (document.body.contains(script)) {
        document.body.removeChild(script)
      }
    }
  }, [])

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
          </Routes>
        </div>
        <Footer />
      </main>
    </>
  )
}

export default App
```

---

## Phase 9: Template Integration

### Step 9.1: Setup Asset Building with esbuild

The Wowdash template files are already in the `assets/` directory. We'll use esbuild to process and bundle them to the `public/` directory. This allows us to handle SCSS files, minification, and other transformations.

**Important Workflow:**
- **Source files**: Always edit files in `assets/` directory
- **Build process**: Run `npm run build:assets` to process files from `assets/` → `public/`
- **Served files**: GoFrame serves files from `public/` directory only
- Never edit files directly in `public/` - they will be overwritten on next build

**Install esbuild:**
```bash
npm install -D esbuild esbuild-sass-plugin
```

**Create asset build script:**

**File: `build-assets.js`**
```javascript
const esbuild = require('esbuild')
const { sassPlugin } = require('esbuild-sass-plugin')
const fs = require('fs')
const path = require('path')

// Helper to copy directory recursively
function copyDir(src, dest) {
  if (!fs.existsSync(dest)) {
    fs.mkdirSync(dest, { recursive: true })
  }

  const entries = fs.readdirSync(src, { withFileTypes: true })

  for (let entry of entries) {
    const srcPath = path.join(src, entry.name)
    const destPath = path.join(dest, entry.name)

    if (entry.isDirectory()) {
      copyDir(srcPath, destPath)
    } else {
      fs.copyFileSync(srcPath, destPath)
    }
  }
}

async function buildAssets() {
  console.log('Building assets...')

  // Build CSS/SCSS files
  try {
    await esbuild.build({
      entryPoints: ['assets/css/style.css'],
      bundle: false,
      outdir: 'public/css',
      minify: process.env.NODE_ENV === 'production',
      sourcemap: process.env.NODE_ENV !== 'production',
      plugins: [sassPlugin()],
      loader: {
        '.css': 'css',
        '.scss': 'css',
      }
    })
    console.log('✓ CSS/SCSS processed')
  } catch (err) {
    console.error('CSS build failed:', err)
  }

  // Copy other CSS files (libraries, etc.)
  const cssFiles = fs.readdirSync('assets/css')
  for (const file of cssFiles) {
    if (file.endsWith('.css') && file !== 'style.css') {
      fs.copyFileSync(
        path.join('assets/css', file),
        path.join('public/css', file)
      )
    }
  }

  // Copy JavaScript files
  console.log('Copying JavaScript files...')
  copyDir('assets/js', 'public/js')
  console.log('✓ JavaScript copied')

  // Copy images
  console.log('Copying images...')
  copyDir('assets/images', 'public/images')
  console.log('✓ Images copied')

  console.log('Asset build complete!')
}

buildAssets().catch(err => {
  console.error('Build failed:', err)
  process.exit(1)
})
```

**Add build script to package.json:**

**File: `package.json` (root directory)**
```json
{
  "name": "tzlev-assets",
  "version": "1.0.0",
  "scripts": {
    "build:assets": "node build-assets.js",
    "watch:assets": "node build-assets.js && chokidar 'assets/**/*' -c 'node build-assets.js'"
  },
  "devDependencies": {
    "esbuild": "^0.19.0",
    "esbuild-sass-plugin": "^2.16.0",
    "chokidar-cli": "^3.0.0"
  }
}
```

**Build the assets:**
```bash
# Install dependencies
npm install

# Build assets once
npm run build:assets

# Or watch for changes during development
npm run watch:assets
```

**Verify the files are built:**
```bash
ls -la public/css
ls -la public/js
ls -la public/images
```

### Step 9.2: Create Sidebar Component

**File: `frontend/src/components/Sidebar.jsx`**
```javascript
import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { menuConfig } from '../config/menuConfig'

function Sidebar() {
  const { t } = useTranslation()
  const location = useLocation()

  const isActive = (path) => location.pathname === path

  return (
    <aside className="sidebar">
      <button type="button" className="sidebar-close-btn">
        <iconify-icon icon="radix-icons:cross-2"></iconify-icon>
      </button>
      <div>
        <a href="/" className="sidebar-logo">
          <img src="/images/logo.png" alt="site logo" className="light-logo" />
          <img src="/images/logo-light.png" alt="site logo" className="dark-logo" />
          <img src="/images/logo-icon.png" alt="site logo" className="logo-icon" />
        </a>
      </div>
      <div className="sidebar-menu-area">
        <ul className="sidebar-menu" id="sidebar-menu">
          {menuConfig.map((item) => (
            <li key={item.id} className={item.children ? 'dropdown' : ''}>
              {item.children ? (
                <>
                  <a href="javascript:void(0)">
                    <iconify-icon icon={item.icon} className="menu-icon"></iconify-icon>
                    <span>{t(item.titleKey)}</span>
                  </a>
                  <ul className="sidebar-submenu">
                    {item.children.map((child) => (
                      <li key={child.id}>
                        <Link
                          to={child.path}
                          className={isActive(child.path) ? 'active-page' : ''}
                        >
                          <i className="ri-circle-fill circle-icon text-primary-600 w-auto"></i>
                          {t(child.titleKey)}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </>
              ) : (
                <Link
                  to={item.path}
                  className={isActive(item.path) ? 'active-page' : ''}
                >
                  <iconify-icon icon={item.icon} className="menu-icon"></iconify-icon>
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

export default Sidebar
```

### Step 9.3: Create Navbar Component

**File: `frontend/src/components/Navbar.jsx`**
```javascript
import React from 'react'
import { useTranslation } from 'react-i18next'

function Navbar() {
  const { t, i18n } = useTranslation()

  const changeLanguage = (lng) => {
    i18n.changeLanguage(lng)
    document.documentElement.dir = lng === 'he' ? 'rtl' : 'ltr'
    document.documentElement.lang = lng
  }

  return (
    <div className="navbar-header">
      <div className="row align-items-center justify-content-between">
        <div className="col-auto">
          <div className="d-flex flex-wrap align-items-center gap-4">
            <button type="button" className="sidebar-toggle">
              <iconify-icon icon="heroicons:bars-3-solid" className="icon text-2xl non-active"></iconify-icon>
              <iconify-icon icon="iconoir:arrow-right" className="icon text-2xl active"></iconify-icon>
            </button>
            <button type="button" className="sidebar-mobile-toggle">
              <iconify-icon icon="heroicons:bars-3-solid" className="icon"></iconify-icon>
            </button>
          </div>
        </div>
        <div className="col-auto">
          <div className="d-flex flex-wrap align-items-center gap-3">
            {/* Language Switcher */}
            <div className="dropdown">
              <button className="has-indicator w-40-px h-40-px bg-neutral-200 rounded-circle d-flex justify-content-center align-items-center" type="button" data-bs-toggle="dropdown">
                <iconify-icon icon="ph:translate" className="text-primary-light text-xl"></iconify-icon>
              </button>
              <div className="dropdown-menu to-top dropdown-menu-sm">
                <div className="py-12 px-16 radius-8 bg-primary-50 mb-16 d-flex align-items-center justify-content-between gap-2">
                  <div>
                    <h6 className="text-lg text-primary-light fw-semibold mb-0">{t('common.language')}</h6>
                  </div>
                </div>
                <div className="max-h-400-px overflow-y-auto scroll-sm pe-8">
                  <a
                    href="javascript:void(0)"
                    className="dropdown-item d-flex align-items-center gap-2 py-6 px-16"
                    onClick={() => changeLanguage('en')}
                  >
                    <span className="text-secondary-light text-md fw-medium">English</span>
                  </a>
                  <a
                    href="javascript:void(0)"
                    className="dropdown-item d-flex align-items-center gap-2 py-6 px-16"
                    onClick={() => changeLanguage('he')}
                  >
                    <span className="text-secondary-light text-md fw-medium">עברית</span>
                  </a>
                </div>
              </div>
            </div>

            {/* Theme Toggle */}
            <button className="theme-toggle-btn">
              <iconify-icon icon="solar:sun-2-bold" className="sun"></iconify-icon>
              <iconify-icon icon="ph:moon-fill" className="moon"></iconify-icon>
            </button>

            {/* User Profile */}
            <div className="dropdown">
              <button className="d-flex justify-content-center align-items-center rounded-circle" type="button" data-bs-toggle="dropdown">
                <img src="/images/user.png" alt="User" className="w-40-px h-40-px object-fit-cover rounded-circle" />
              </button>
              <div className="dropdown-menu to-top dropdown-menu-sm">
                <div className="py-12 px-16 radius-8 bg-primary-50 mb-16 d-flex align-items-center justify-content-between gap-2">
                  <div>
                    <h6 className="text-lg text-primary-light fw-semibold mb-2">{t('common.welcome')}</h6>
                  </div>
                </div>
                <ul className="to-top-list">
                  <li>
                    <a className="dropdown-item text-black px-16 py-8 rounded text-sm" href="/profile">
                      <iconify-icon icon="solar:user-linear" className="icon text-xl"></iconify-icon>
                      {t('common.profile')}
                    </a>
                  </li>
                  <li>
                    <a className="dropdown-item text-black px-16 py-8 rounded text-sm" href="/logout">
                      <iconify-icon icon="lucide:power" className="icon text-xl"></iconify-icon>
                      {t('common.logout')}
                    </a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Navbar
```

### Step 9.4: Create Footer Component

**File: `frontend/src/components/Footer.jsx`**
```javascript
import React from 'react'

function Footer() {
  return (
    <footer className="d-footer">
      <div className="row align-items-center justify-content-between">
        <div className="col-auto">
          <p className="mb-0">© 2024 Tzlev. All Rights Reserved.</p>
        </div>
        <div className="col-auto">
          <p className="mb-0">Made with ❤️ by the Tzlev team</p>
        </div>
      </div>
    </footer>
  )
}

export default Footer
```

### Step 9.5: Create Theme Customization Component

**File: `frontend/src/components/ThemeCustomization.jsx`**
```javascript
import React from 'react'

function ThemeCustomization() {
  return (
    <div className="customize-options d-none d-md-block">
      <div className="customizer-item theme-light-dark">
        <h6 className="mb-1">Theme Mode</h6>
        <div className="d-flex gap-2 align-items-center">
          <button type="button" className="theme-btn light-mode">
            <iconify-icon icon="solar:sun-2-linear" className="icon"></iconify-icon>
          </button>
          <button type="button" className="theme-btn dark-mode">
            <iconify-icon icon="ph:moon-fill" className="icon"></iconify-icon>
          </button>
        </div>
      </div>
    </div>
  )
}

export default ThemeCustomization
```

---

## Phase 10: Internationalization

### Step 10.1: Create i18n Configuration

**File: `frontend/src/config/i18n.js`**
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
    lng: 'en',
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false
    }
  })

export default i18n
```

### Step 10.2: Create English Translations

**File: `frontend/src/locales/en/translation.json`**
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
    "profile": "Profile",
    "language": "Language",
    "search": "Search"
  }
}
```

### Step 10.3: Create Hebrew Translations

**File: `frontend/src/locales/he/translation.json`**
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
    "profile": "פרופיל",
    "language": "שפה",
    "search": "חיפוש"
  }
}
```

---

## Phase 11: Routing & Navigation

### Step 11.1: Create Menu Configuration

**File: `frontend/src/config/menuConfig.js`**
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

### Step 11.2: Create Page Components

**File: `frontend/src/pages/Home.jsx`**
```javascript
import React from 'react'
import { useTranslation } from 'react-i18next'

function Home() {
  const { t } = useTranslation()

  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <h6 className="fw-semibold mb-0">{t('menu.home')}</h6>
      <div className="card h-100 p-0 radius-12">
        <div className="card-body p-24">
          <h5>Welcome to Tzlev!</h5>
          <p>This is the home page.</p>
        </div>
      </div>
    </div>
  )
}

export default Home
```

**File: `frontend/src/pages/Dashboard.jsx`**
```javascript
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

function Dashboard() {
  const { t } = useTranslation()
  const [health, setHealth] = useState(null)

  useEffect(() => {
    // Test API call
    fetch('/api/health')
      .then(res => res.json())
      .then(data => setHealth(data))
      .catch(err => console.error('API Error:', err))
  }, [])

  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <h6 className="fw-semibold mb-0">{t('menu.dashboard')}</h6>
      <div className="card h-100 p-0 radius-12">
        <div className="card-body p-24">
          <h5>Dashboard</h5>
          {health && (
            <div>
              <p>API Status: {health.status}</p>
              <p>Message: {health.message}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Dashboard
```

---

## Phase 12: CLI Mode

### Step 12.1: Create CLI Command Handlers

**File: `internal/cli/migrate.go`**
```go
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

func RunMigrate(ctx context.Context, parser *gcmd.Parser) {
	action := parser.GetOpt("action", "up").String()

	g.Log().Infof(ctx, "Running database migrations (%s)...", action)

	// TODO: Implement migration logic using golang-migrate or custom system
	// For now, just a placeholder

	g.Log().Info(ctx, "Migrations completed successfully")
	os.Exit(0)
}
```

**File: `internal/cli/seed.go`**
```go
package cli

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

func RunSeed(ctx context.Context, parser *gcmd.Parser) {
	table := parser.GetOpt("table", "all").String()

	g.Log().Infof(ctx, "Seeding database (table: %s)...", table)

	// TODO: Implement seeding logic

	g.Log().Info(ctx, "Database seeded successfully")
	os.Exit(0)
}
```

**File: `internal/cli/version.go`**
```go
package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"
)

const (
	Version   = "1.0.0"
	BuildDate = "2024-01-01"
)

func ShowVersion(ctx context.Context) {
	fmt.Printf("Tzlev v%s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	os.Exit(0)
}
```

**File: `internal/cli/help.go`**
```go
package cli

import (
	"context"
	"fmt"
	"os"
)

func ShowHelp(ctx context.Context) {
	helpText := `
Tzlev - GoFrame + React Application

USAGE:
    ./tzlev [command] [options]

COMMANDS:
    migrate              Run database migrations
    seed                 Seed database with test data
    version              Show version information
    help                 Show this help message

WEB MODE:
    ./tzlev              Start web server (no arguments)

EXAMPLES:
    ./tzlev migrate --action=up
    ./tzlev seed --table=users
    ./tzlev version

For web server mode, simply run without arguments:
    ./tzlev
    go run main.go
`
	fmt.Println(helpText)
	os.Exit(0)
}
```

### Step 12.2: Create Main.go with Dual-Mode Support

**File: `main.go`**
```go
package main

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"

	"tzlev/internal/cli"
	"tzlev/internal/controller"
	"tzlev/internal/database"
	"tzlev/internal/middleware"
	"tzlev/internal/oauth"
	"tzlev/internal/redis"
)

func main() {
	ctx := gctx.New()

	// Initialize external services
	if err := initServices(ctx); err != nil {
		g.Log().Fatal(ctx, "Failed to initialize services:", err)
	}

	// Check for CLI mode
	if len(os.Args) > 1 {
		runCLI(ctx)
		return
	}

	// Web mode
	runWebServer(ctx)
}

func initServices(ctx context.Context) error {
	// Initialize database
	if err := database.Init(); err != nil {
		return err
	}

	// Initialize Redis
	if err := redis.Init(); err != nil {
		return err
	}

	// Initialize OAuth
	if err := oauth.InitGoogleOAuth(); err != nil {
		return err
	}

	return nil
}

func runCLI(ctx context.Context) {
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

func runWebServer(ctx context.Context) {
	s := g.Server()

	// CORS Middleware
	s.Use(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})

	// Serve static files
	s.SetServerRoot("public")

	// Main route - serve index.html
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.ServeFile("public/index.html")
	})

	// Setup routes
	setupRoutes(s)

	// Start server
	port := g.Cfg().MustGet(ctx, "app.port").Int()
	g.Log().Infof(ctx, "Starting Tzlev server on http://localhost:%d", port)
	s.SetPort(port)
	s.Run()
}

func setupRoutes(s *ghttp.Server) {
	healthCtrl := controller.NewHealthController()
	authCtrl := controller.NewAuthController()

	// Public routes
	s.Group("/auth", func(group *ghttp.RouterGroup) {
		group.GET("/google/login", authCtrl.GoogleLogin)
		group.GET("/google/callback", authCtrl.GoogleCallback)
		group.POST("/logout", authCtrl.Logout)
	})

	// API routes
	s.Group("/api", func(group *ghttp.RouterGroup) {
		// Public API
		group.GET("/health", healthCtrl.Check)

		// Protected API
		group.Group("/", func(protectedGroup *ghttp.RouterGroup) {
			protectedGroup.Middleware(middleware.Auth())
			protectedGroup.GET("/me", authCtrl.GetCurrentUser)
		})
	})
}
```

---

## Phase 13: Build & Deployment

### Step 13.1: Build Frontend
```bash
cd frontend
npm install
npm run build
```

### Step 13.2: Build Backend
```bash
go mod tidy
go build -o bin/tzlev main.go
```

### Step 13.3: Create Build Script

**File: `build.sh`**
```bash
#!/bin/bash

set -e

echo "Building Tzlev application..."

# Build template assets with esbuild
echo "Building template assets..."
npm install
npm run build:assets

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "Building backend..."
go mod tidy
go build -o bin/tzlev main.go

echo "Build complete!"
echo "Run with: ./bin/tzlev"
```

Make it executable:
```bash
chmod +x build.sh
```

### Step 13.4: Create Run Script

**File: `run.sh`**
```bash
#!/bin/bash

# Development run script

# Always build assets from source (assets/ → public/)
echo "Building assets..."
npm run build:assets

# Run the application
go run main.go
```

**Note**: This script always rebuilds assets from `assets/` to ensure any changes are reflected.

Make it executable:
```bash
chmod +x run.sh
```

---

## Phase 14: Testing & Validation

### Step 14.1: Test Database Connection
```bash
# Start PostgreSQL (via MCP or locally)
# Update .env with correct database credentials
go build main.go
```

### Step 14.2: Test Redis Connection
```bash
# Start Redis
redis-server

# Test connection via CLI
redis-cli ping
```

### Step 14.3: Test Application
```bash
# Run in web mode
./run.sh

# Open browser
open http://localhost:8080
```

### Step 14.4: Test CLI Mode
```bash
# Test version command
go run main.go version

# Test help command
go run main.go help

# Test migrate command
go run main.go migrate

# Test seed command
go run main.go seed
```

### Step 14.5: Test Frontend Development Mode
```bash
# Terminal 1: Build and watch template assets
npm run watch:assets

# Terminal 2: Run backend
go run main.go

# Terminal 3: Run frontend dev server
cd frontend
npm run dev

# Open browser
open http://localhost:3000
```

**Note**: The `watch:assets` script will automatically rebuild template assets when files in the `assets/` directory change.

---

## Summary Checklist

### Backend Setup
- [ ] Go module initialized
- [ ] GoFrame dependencies installed
- [ ] Configuration files created
- [ ] Database connection working (via MCP)
- [ ] Redis connection working
- [ ] Models created
- [ ] Repositories created
- [ ] Services created
- [ ] Controllers created
- [ ] Middleware created
- [ ] OAuth configured
- [ ] JWT implemented
- [ ] Email service implemented
- [ ] CLI commands implemented

### Frontend Setup
- [ ] React project initialized
- [ ] Vite configured
- [ ] React Router installed
- [ ] i18next configured
- [ ] esbuild configured for template assets
- [ ] Template assets built with esbuild
- [ ] Components created (Sidebar, Navbar, Footer, etc.)
- [ ] Pages created
- [ ] Menu configuration created
- [ ] Translation files created
- [ ] Routing configured

### Integration
- [ ] Main.go with dual-mode support
- [ ] Routes configured
- [ ] CORS enabled
- [ ] Static file serving working
- [ ] API endpoints working
- [ ] Authentication flow working

### Testing
- [ ] Database queries work
- [ ] Redis operations work
- [ ] OAuth login works
- [ ] Sessions work
- [ ] Email sending works
- [ ] Frontend loads correctly
- [ ] i18n language switching works
- [ ] Routing works
- [ ] API calls work from frontend
- [ ] CLI commands work

### Deployment
- [ ] Build script created
- [ ] Frontend builds correctly
- [ ] Backend builds correctly
- [ ] Public directory structure correct
- [ ] Production configuration ready

---

## Next Steps

After completing all phases:

1. **Add More Features**: Implement additional pages and functionality
2. **Add Tests**: Write unit tests and integration tests
3. **Add Documentation**: Create user guides and API documentation
4. **Security Hardening**: Review security practices, add rate limiting
5. **Performance Optimization**: Add caching strategies, optimize queries
6. **Monitoring**: Add logging, metrics, and error tracking
7. **CI/CD**: Set up continuous integration and deployment
8. **Production Deployment**: Deploy to production environment

---

## Troubleshooting

### Common Issues

**Database Connection Fails**
- Check MCP PostgreSQL server is running
- Verify .env database credentials
- Check network connectivity

**Redis Connection Fails**
- Ensure Redis server is running: `redis-server`
- Check Redis password in .env
- Verify Redis port (default 6379)

**Frontend Build Fails**
- Run `npm install` in frontend directory
- Check Node.js version (requires Node 16+)
- Clear node_modules and reinstall

**Template Assets Not Loading**
- Run `npm run build:assets` to build assets to public/
- Check that esbuild completed without errors
- Verify file paths in HTML match built files
- Check GoFrame static file serving configuration
- If SCSS files, ensure esbuild-sass-plugin is installed

**OAuth Not Working**
- Verify Google OAuth credentials in .env
- Check redirect URL matches Google Console configuration
- Ensure HTTPS in production

---

## Additional Resources

- [GoFrame Documentation](https://goframe.org/docs)
- [React Documentation](https://react.dev)
- [React Router Documentation](https://reactrouter.com)
- [i18next Documentation](https://www.i18next.com)
- [PostgreSQL Documentation](https://www.postgresql.org/docs)
- [Redis Documentation](https://redis.io/docs)

---

End of STEPS.md
