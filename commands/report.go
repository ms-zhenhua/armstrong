package commands

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/report"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ReportCommand struct {
	Ui          cli.Ui
	workingDir  string
	swaggerPath string
}

func (c *ReportCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("report")
	fs.StringVar(&c.workingDir, "working-dir", "", "path that contains all the test results")
	fs.StringVar(&c.swaggerPath, "swagger-path", "", "path to the .json swagger which is being test")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c ReportCommand) Help() string {
	helpText := `
Usage: armstrong report [-working-dir <path that contains all the test results>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c ReportCommand) Synopsis() string {
	return "Generate test report for a set of test results"
}

func (c ReportCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c ReportCommand) Execute() int {
	log.Println("[INFO] ----------- generate report ---------")
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}

	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("working directory is invalid: %+v", err))
			return 1
		}
	}

	if report.MergeApiTestReports(wd, c.swaggerPath) != nil {
		log.Fatalf("[ERROR] failed to generate test reports: %+v", err)
	}

	return 0
}
