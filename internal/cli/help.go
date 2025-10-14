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
