// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/conformal/btcutil"
	"github.com/conformal/go-flags"
)

// config defines the configuration options for blocksafari.
//
// See loadConfig for details on the configuration load process.
type config struct {
	ConfigFile  string   `short:"C" long:"configfile" description:"Path to configuration file"`
	Listeners   []string `long:"listen" description:"Add an interface/port to listen on"`
	RPCCert     string   `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	RPCServer   string   `short:"s" long:"rpcserver" description:"IP and port for rpcserver."`
	RPCUser     string   `short:"u" long:"rpcuser" description:"rpc username."`
	RPCPassword string   `short:"P" long:"rpcpass" description:"rpc password."`
}

const (
	defaultConfigFilename = "blocksafari.conf"
)

var (
	btcdHomeDir        = btcutil.AppDataDir("btcd", false)
	bsHomeDir          = btcutil.AppDataDir("blocksafari", false)
	cfg                *config
	defaultConfigFile  = filepath.Join(bsHomeDir, defaultConfigFilename)
	defaultRPCCertFile = filepath.Join(btcdHomeDir, "rpc.cert")

	pem []byte
)

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
//      1) Start with a default config with sane settings
//      2) Pre-parse the command line to check for an alternative config file
//      3) Load configuration file overwriting defaults with any specified options
//      4) Parse CLI options and overwrite/add any specified options
//
// Command line options always take precedence.
func loadConfig() (*config, []string, error) {
	cfg := config{
		ConfigFile: defaultConfigFile,
		RPCCert:    defaultRPCCertFile,
	}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.  Any errors can be ignored
	// here since they will be caught be the final parse below.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.None)
	preParser.Parse()

	// Load config from file.
	parser := flags.NewParser(&cfg, flags.Default)
	err := flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
	if err != nil {
		return nil, nil, err
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}

	pem, err = ioutil.ReadFile(cfg.RPCCert)
	if err != nil {
		return nil, nil, err
	}

	return &cfg, remainingArgs, nil
}
