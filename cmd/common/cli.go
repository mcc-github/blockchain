/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"fmt"
	"os"
	"path/filepath"

	"io"

	"github.com/mcc-github/blockchain/cmd/common/comm"
	"github.com/mcc-github/blockchain/cmd/common/signer"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	saveConfigCommand = "saveConfig"
)

var (
	
	terminate = os.Exit
	
	outWriter io.Writer = os.Stderr

	
	mspID                                     *string
	tlsCA, tlsCert, tlsKey, userKey, userCert **os.File
	configFile                                *string
)



type CLICommand func(Config) error


type CLI struct {
	app         *kingpin.Application
	dispatchers map[string]CLICommand
}


func NewCLI(name, help string) *CLI {
	return &CLI{
		app:         kingpin.New(name, help),
		dispatchers: make(map[string]CLICommand),
	}
}


func (cli *CLI) Command(name, help string, onCommand CLICommand) *kingpin.CmdClause {
	cmd := cli.app.Command(name, help)
	cli.dispatchers[name] = onCommand
	return cmd
}


func (cli *CLI) Run(args []string) {
	configFile = cli.app.Flag("configFile", "Specifies the config file to load the configuration from").String()
	persist := cli.app.Command(saveConfigCommand, fmt.Sprintf("Save the config passed by flags into the file specified by --configFile"))
	configureFlags(cli.app)

	command := kingpin.MustParse(cli.app.Parse(args))
	if command == persist.FullCommand() {
		if *configFile == "" {
			out("--configFile must be used to specify the configuration file")
			return
		}
		persistConfig(parseFlagsToConfig(), *configFile)
		return
	}

	var conf Config
	if *configFile == "" {
		conf = parseFlagsToConfig()
	} else {
		conf = loadConfig(*configFile)
	}

	f, exists := cli.dispatchers[command]
	if !exists {
		out("Unknown command:", command)
		terminate(1)
		return
	}
	err := f(conf)
	if err != nil {
		out(err)
		terminate(1)
		return
	}
}

func configureFlags(persistCommand *kingpin.Application) {
	
	tlsCA = persistCommand.Flag("peerTLSCA", "Sets the TLS CA certificate file path that verifies the TLS peer's certificate").File()
	tlsCert = persistCommand.Flag("tlsCert", "(Optional) Sets the client TLS certificate file path that is used when the peer enforces client authentication").File()
	tlsKey = persistCommand.Flag("tlsKey", "(Optional) Sets the client TLS key file path that is used when the peer enforces client authentication").File()
	
	userKey = persistCommand.Flag("userKey", "Sets the user's key file path that is used to sign messages sent to the peer").File()
	userCert = persistCommand.Flag("userCert", "Sets the user's certificate file path that is used to authenticate the messages sent to the peer").File()
	mspID = persistCommand.Flag("MSP", "Sets the MSP ID of the user, which represents the CA(s) that issued its user certificate").String()
}

func persistConfig(conf Config, file string) {
	if err := conf.ToFile(file); err != nil {
		out("Failed persisting configuration:", err)
		terminate(1)
	}
}

func loadConfig(file string) Config {
	conf, err := ConfigFromFile(file)
	if err != nil {
		out("Failed loading config", err)
		terminate(1)
		return Config{}
	}
	return conf
}

func parseFlagsToConfig() Config {
	conf := Config{
		SignerConfig: signer.Config{
			MSPID:        *mspID,
			IdentityPath: evaluateFileFlag(userCert),
			KeyPath:      evaluateFileFlag(userKey),
		},
		TLSConfig: comm.Config{
			KeyPath:        evaluateFileFlag(tlsKey),
			CertPath:       evaluateFileFlag(tlsCert),
			PeerCACertPath: evaluateFileFlag(tlsCA),
		},
	}
	return conf
}

func evaluateFileFlag(f **os.File) string {
	if f == nil {
		return ""
	}
	if *f == nil {
		return ""
	}
	path, err := filepath.Abs((*f).Name())
	if err != nil {
		out("Failed listing", (*f).Name(), ":", err)
		terminate(1)
	}
	return path
}
func out(a ...interface{}) {
	fmt.Fprintln(outWriter, a...)
}
