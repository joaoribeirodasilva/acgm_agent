package main

import (
	"fmt"

	"biqx.com.br/acgm_agent/modules/cmd"
	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
)

func main() {

	options, err := cmd.Parse()
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	// fmt.Printf("%+v, %s, %t\n", options, *options.ConfigFile, *options.Service)

	config := &config.Config{}
	err = config.Read(options)
	if err != nil {
		panic("ERROR: " + err.Error())
	}

	database := database.New(config)
	err = database.Connect()
	if err != nil {
		panic(fmt.Sprintf("ERROR: Error connecting to the database! %s", err.Error()))
	}

	fmt.Printf("Configuration: \n%+v", config)
	fmt.Printf("Database: \n%+v", database)

	database.Disconnect()

}
