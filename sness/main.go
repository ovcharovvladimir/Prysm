// Package sness defines all the utlities needed for a beacon chain node.
package main

import (
	"os"
	"runtime"

	"strconv"

	"github.com/mattn/go-colorable"
	"github.com/ovcharovvladimir/Prysm/shared/cmd"
	"github.com/ovcharovvladimir/Prysm/shared/debug"
	"github.com/ovcharovvladimir/Prysm/sness/node"
	"github.com/ovcharovvladimir/Prysm/sness/utils"
	"github.com/ovcharovvladimir/essentiaHybrid/log"

	//"github.com/urfave/cli"
	"gopkg.in/urfave/cli.v1"
)

var (
	err error
)

func startNode(ctx *cli.Context) error {
	var verbosity = 3
	if ctx.IsSet(cmd.VerbosityFlag.Name) {
		vr := ctx.GlobalString(cmd.VerbosityFlag.Name)
		verbosity, err = strconv.Atoi(vr)
		if err != nil {
			// handle error
			log.Crit(err.Error())
		}
	}

	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(verbosity), log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat(true))))

	beacon, err := node.NewBeaconNode(ctx)
	if err != nil {
		return err
	}
	beacon.Start()
	return nil
}

func main() {

	app := cli.NewApp()
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	app.Name = "sness"
	app.Usage = "this is a supernode implementation for Essentia"
	app.Action = startNode

	app.Flags = []cli.Flag{
		utils.DemoConfigFlag,
		utils.SimulatorFlag,
		utils.VrcContractFlag,
		utils.PubKeyFlag,
		utils.Web3ProviderFlag,
		utils.RPCPort,
		utils.CertFlag,
		utils.KeyFlag,
		utils.GenesisJSON,
		utils.EnableCrossLinks,
		utils.EnableRewardChecking,
		utils.EnableAttestationValidity,
		utils.EnablePOWChain,
		cmd.DataDirFlag,
		cmd.VerbosityFlag,
		cmd.EnableTracingFlag,
		cmd.TracingEndpointFlag,
		cmd.TraceSampleFractionFlag,
		debug.PProfFlag,
		debug.PProfAddrFlag,
		debug.PProfPortFlag,
		debug.MemProfileRateFlag,
		debug.CPUProfileFlag,
		debug.TraceFlag,
	}

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return debug.Setup(ctx)
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
