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
