package main

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL driver

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
	cfg := g.Cfg()

	// CORS Middleware
	s.Use(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})

	// Setup routes (must be before static file serving)
	setupRoutes(s)

	// Serve static files
	s.SetServerRoot("public")

	// SPA fallback - serve index.html for all non-API routes
	// This allows React Router to handle routing on page refresh
	s.BindHandler("/*", func(r *ghttp.Request) {
		path := r.URL.Path
		// Don't intercept API or auth routes
		if len(path) >= 4 && path[:4] == "/api" {
			r.Response.WriteStatus(404)
			return
		}
		if len(path) >= 5 && path[:5] == "/auth" {
			r.Response.WriteStatus(404)
			return
		}
		// Serve index.html for all other routes
		r.Response.ServeFile("public/index.html")
	})

	// Start server
	port := cfg.MustGet(ctx, "app.port").Int()
	g.Log().Infof(ctx, "Starting Tzlev server on http://localhost:%d", port)
	s.SetPort(port)
	s.Run()
}

func setupRoutes(s *ghttp.Server) {
	healthCtrl := controller.NewHealthController()
	authCtrl := controller.NewAuthController()
	academicYearCtrl := controller.NewAcademicYearController()
	appResourceCtrl := controller.NewAppResourceController()

	// Public routes
	s.Group("/auth", func(group *ghttp.RouterGroup) {
		group.POST("/login", authCtrl.Login)
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
			protectedGroup.GET("/academic-year", academicYearCtrl.GetAcademicYear)
			protectedGroup.POST("/academic-year", academicYearCtrl.SetAcademicYear)
			protectedGroup.GET("/academic-years", academicYearCtrl.GetAcademicYearsList)
			protectedGroup.GET("/app-resources", appResourceCtrl.GetAppResources)
			protectedGroup.GET("/app-resources/{id}", appResourceCtrl.GetAppResource)
			protectedGroup.POST("/app-resources", appResourceCtrl.CreateAppResource)
			protectedGroup.PUT("/app-resources/{id}", appResourceCtrl.UpdateAppResource)
			protectedGroup.DELETE("/app-resources/{id}", appResourceCtrl.DeleteAppResource)
		})
	})
}
