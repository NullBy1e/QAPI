package main

import (
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	fileName string
	jsonData string
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Value:       "",
				Aliases:     []string{"f"},
				Usage:       "Reads JSON from file, returned on POST",
				Destination: &fileName,
			},
			&cli.StringFlag{
				Name:        "json",
				Value:       "",
				Aliases:     []string{"j"},
				Usage:       "JSON returned on POST",
				Destination: &jsonData,
			},
		},
		Action: func(ctx *cli.Context) error {
			port, _ := strconv.Atoi(ctx.Args().Get(1))
			if port > 99999 {
				cli.Exit("Port must be in range 0-99999", 2)
			}

			if port == 0 {
				port = 5000
			}

			startAPIServer(port, jsonData, fileName)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func startAPIServer(port int, postData string, fileName string) {
	color.HiYellow("Starting the API endpoint on port: %d", port)
}
