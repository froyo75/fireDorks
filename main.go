package main

import (
	"fireDorks/libs"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Context struct {
	Debug bool
}
type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

var cli struct {
	ProxyUrl       string      `short:"x" name:"proxyurl" required:"" default:"https://www.google.com" help:"The proxy URL to use (default: https://www.google.com)."`
	Dork           string      `short:"k" name:"dork" required:"" help:"The dork query to provide. (example: site:microsoft.com)."`
	MaxPages       int         `short:"p" name:"maxpages" optional:"" default:"1000" help:"Specify the maximum number of Google search pages (default:1000)."`
	ResultsPerPage int         `short:"r" name:"maxresults" optional:"" default:"100" help:"Specify the number of results per page to return (default:100)."`
	MaxWorkers     int         `short:"w" name:"maxworkers" optional:"" default:"5" help:"Specify the number of workers to process queries (default:5)."`
	Pattern        string      `short:"s" name:"regex" optional:"" help:"Specify a pattern to search using the POSIX standard regex."`
	RawOnly        bool        `short:"e" name:"rawonly" optional:"" help:"Specify whether values should be extracted from raw HTTP response data using a provided pattern."`
	LinksOnly      bool        `short:"l" name:"linksonly" optional:"" default:"false" help:"Specify whether to return only links in the output results (default:false)."`
	Format         string      `short:"f" name:"format" optional:"" default:"txt" help:"Specify the output format (txt/csv/json) (default:txt)."`
	Delay          int         `short:"d" name:"delay" optional:"" default:"2" help:"Specify the delay (in seconds) to avoid being fingerprinting (default:2)."`
	OutConsole     bool        `short:"c" name:"outconsole" optional:"" default:"false" help:"Specify whether to output all results to the console (default:false)."`
	OutPath        string      `short:"o" name:"outpath" help:"Output results to a specific file."`
	Version        VersionFlag `name:"version" short:"v" help:"Print version information and quit."`
}

func main() {
	start := time.Now()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "02 Jan 2006 15:04:05"})
	ctx := kong.Parse(&cli, kong.Name("FireDorks"), kong.Description("A tool for scraping Google Search results by @froyo75"), kong.UsageOnError(), kong.Vars{"version": "1.0"})
	ctx.Command()
	urls := libs.GenQueries(cli.Dork, cli.MaxPages, cli.ResultsPerPage, cli.ProxyUrl)
	dataResults := libs.ProcessQuery(cli.Pattern, cli.RawOnly, urls, cli.Delay, cli.MaxWorkers, cli.Format, cli.LinksOnly)
	results := libs.OutResults(dataResults, cli.Format, cli.LinksOnly, cli.RawOnly)
	if len(results) != 0 {
		if cli.OutConsole {
			fmt.Println(results)
		}
		if cli.OutPath != "" {
			libs.OutFile(results, cli.OutPath)

		}
	}
	elapsed := time.Since(start)
	log.Debug().Msg("[*] Execution took " + elapsed.String())
}
