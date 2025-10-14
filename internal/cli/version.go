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
