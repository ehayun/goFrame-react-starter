package database

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// getEnvWithDefault returns the value of an environment variable or a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Init initializes database connection using config.yaml and environment variables
func Init() error {
	ctx := gctx.New()
	cfg := g.Cfg()

	// Read sensitive data from environment variables, static config from config.yaml
	config := gdb.ConfigNode{
		Host:             getEnvWithDefault("DB_HOST", "localhost"),
		Port:             getEnvWithDefault("DB_PORT", "5432"),
		User:             getEnvWithDefault("DB_USER", "postgres"),
		Pass:             getEnvWithDefault("DB_PASSWORD", "postgres"),
		Name:             getEnvWithDefault("DB_NAME", "unified"),
		Type:             cfg.MustGet(ctx, "database.default.type").String(),
		Extra:            cfg.MustGet(ctx, "database.default.extra").String(),
		Debug:            cfg.MustGet(ctx, "database.default.debug").Bool(),
		MaxIdleConnCount: cfg.MustGet(ctx, "database.default.maxIdle").Int(),
		MaxOpenConnCount: cfg.MustGet(ctx, "database.default.maxOpen").Int(),
		MaxConnLifeTime:  cfg.MustGet(ctx, "database.default.maxLifetime").Duration(),
	}

	// Set the configuration
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:             config.Host,
				Port:             config.Port,
				User:             config.User,
				Pass:             config.Pass,
				Name:             config.Name,
				Type:             config.Type,
				Extra:            config.Extra,
				Debug:            config.Debug,
				MaxIdleConnCount: config.MaxIdleConnCount,
				MaxOpenConnCount: config.MaxOpenConnCount,
				MaxConnLifeTime:  config.MaxConnLifeTime,
			},
		},
	})

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
