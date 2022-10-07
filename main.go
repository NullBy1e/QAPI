package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	Filename       string
	OutputFilename string
	PostData       string
	RequestList    []RequestData
)

type RequestData struct {
	url            string
	method         string
	protocol       string
	headers        http.Header
	content_lenght int64
	host           string
	body           string
}

func main() {
	//* Init the Request list logger
	RequestList = []RequestData{}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Value:       "",
				Aliases:     []string{"f"},
				Usage:       "Reads JSON from file, returned on POST",
				Destination: &Filename,
			},
			&cli.StringFlag{
				Name:        "data",
				Value:       "",
				Aliases:     []string{"j"},
				Usage:       "Data returned on POST",
				Destination: &PostData,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "",
				Aliases:     []string{"o"},
				Usage:       "Output file for the requests data",
				Destination: &PostData,
			},
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

			startAPIServer(port, PostData, Filename)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func startAPIServer(port string, postData string, filename string) {
	color.HiYellow("Starting the API endpoint on port: " + port)

	if filename != "" && postData != "" {
		color.Red("Cannot use --file and --json at the same time.")
		return
	} else if filename != "" {
		color.Green("Using file data as POST response")
		// TODO: Read the file and send it to server
	} else if postData != "" {
		color.Green("Using data argument as POST response")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			//* Returns the JSON and prints the data about it
			fmt.Fprintf(w, postData)
			requestDebugger(w, r)
		})
		http.ListenAndServe(":"+port, nil)
	} else {
		color.Green("Starting basic API server")
		http.HandleFunc("/", requestDebugger)
		http.ListenAndServe(":"+port, nil)
	}
}

func requestDebugger(w http.ResponseWriter, r *http.Request) {
	//* Prints the details of the request

	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	request := RequestData{url: r.URL.Path, method: r.Method, protocol: r.Proto, headers: r.Header, host: r.Host, body: string(req_body)}

	fmt.Println("")
	color.Magenta("Request on: " + color.New(color.FgHiMagenta).Sprint(r.URL.Path))
	color.Yellow("Method: " + color.New(color.FgHiYellow).Sprint(request.method))
	color.Cyan("Protocol: " + color.New(color.FgHiCyan).Sprint(request.protocol))
	for i, j := range request.headers {
		fmt.Print(color.RedString(i) + " : ")
		for _, l := range j {
			fmt.Print(color.HiRedString(l) + "\n")
		}
	}
	color.Green("Body: " + color.New(color.FgHiGreen).Sprint(request.body))

	RequestList = append(RequestList, request)
}
