package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/nicus101/godyndns-ovh/internal/config"
)

func executeCommands(config *config.Config) {
	for appName, cmd := range config.Execute {
		cmd := exec.Command(cmd.Command, cmd.Arguments...)
		out, err := cmd.Output()

		if err != nil {
			log.Fatal("Error executing command: ", appName)
		}
		fmt.Println(string(out))
	}

}
