package main

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/alecthomas/kong"
	"github.com/cyclimse/fediverse-blahaj/internal/config"
)

type Context struct {
	Debug   bool
	Config  config.Config
	Ctx     context.Context
	Version string
}

var cli struct {
	Debug  bool          `help:"Enable debug mode."`
	Config config.Config `embed:""`

	API   APICmd   `cmd:"" help:"Start the API." default:"1"`
	Crawl CrawlCmd `cmd:"" help:"Start the crawler."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.

	info, ok := debug.ReadBuildInfo()
	if !ok {
		ctx.FatalIfErrorf(errors.New("failed to read build info"))
	}

	cmdContext := &Context{
		Debug:   cli.Debug,
		Config:  cli.Config,
		Ctx:     context.Background(),
		Version: info.Main.Version,
	}

	err := ctx.Run(cmdContext)
	ctx.FatalIfErrorf(err)
}
