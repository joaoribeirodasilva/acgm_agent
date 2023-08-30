package cmd

import (
	"flag"
)

type CmdOptions struct {
	Service *bool
	ConfigFile *string
}

func Parse() (*CmdOptions, error) {
	service := flag.Bool("service", false, "Run as service");
	config_file := flag.String("config","", "Full path to config file");
	flag.Parse()
	return &CmdOptions{
		Service: service,
		ConfigFile: config_file,
	}, nil;

}