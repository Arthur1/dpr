package cli

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/Arthur1/dpr"
)

type VersionCmd struct{}

func (c *VersionCmd) Run(globals *Globals) error {
	printVersion()
	return nil
}

func printVersion() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "dprc is a cli tool for deploy packages registry.\n")
	fmt.Fprintf(writer, "Version:\t%s\n", dpr.Version)
	fmt.Fprintf(writer, "Go version:\t%s\n", runtime.Version())
	fmt.Fprintf(writer, "Arch:\t%s\n", runtime.GOARCH)
	writer.Flush()
}
