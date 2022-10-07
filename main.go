package main

import (
	"fmt"
	"log"
	"net/http"
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
			//TODO: Create Output file to JSON
		},
		Action: func(ctx *cli.Context) error {

			port := ctx.Args().Get(0)
			if num, _ := strconv.Atoi(port); num > 99999 {
				cli.Exit("Port must be in range 0-99999", 2)
			}

			log.Println(port)
			if port == "" {
				port = "5000"
			}

			startAPIServer(port, jsonData, fileName)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func startAPIServer(port string, postData string, fileName string) {
	color.HiYellow("Starting the API endpoint on port: " + port)

	if fileName != "" && postData != "" {
		color.Red("Cannot use --file and --json at the same time.")
		return
	} else if fileName != "" {
		color.Green("Using file data as POST")
		// TODO: Read the file and send it to server
	} else if postData != "" {
		color.Green("Using JSON argument as POST")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			//* Returns the JSON and prints the data about it
			fmt.Fprintf(w, postData)
			requestDebugger(w, r)
		})
		log.Printf(":" + port)
		http.ListenAndServe(":"+port, nil)
	} else {
		color.Green("Starting basic API server")

		// TODO: Don't bind only to /, rather to everything!
		http.HandleFunc("/", requestDebugger)
		http.ListenAndServe(":"+port, nil)
	}
}

func requestDebugger(w http.ResponseWriter, r *http.Request) {
	//* Prints the details of the request
	log.Println("Incoming Request on:", r.URL)
}
