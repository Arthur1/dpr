package dpr

import (
	"github.com/alecthomas/kong"
)

var (
	Version string = "unknown"
)

type Globals struct {
	Version    VersionFlag `name:"version" alias:"v" help:"print version and quit"`
	ConfigFile string      `name:"config" alias:"c" help:"dprc config file" default:"./dprcconfig.yml" type:"path"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error {
	return nil
}
func (v VersionFlag) IsBool() bool {
	return true
}
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	printVersion()
	app.Exit(0)
	return nil
}

var cli struct {
	Globals
	Push                 PushCmd                 `cmd:"" help:"push deploy package"`
	Pull                 PullCmd                 `cmd:"" help:"pull deploy package"`
	ApplyLifecyclePolicy ApplyLifecyclePolicyCmd `cmd:"" help:"apply lifecycle policy and delete expired packages"`
	Version              VersionCmd              `cmd:"" help:"print version information"`
}

type Cli struct{}

func (c *Cli) Run() {
	kctx := kong.Parse(&cli)
	err := kctx.Run(&cli.Globals)
	kctx.FatalIfErrorf(err)
}
