package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	InputFile    string
	OutputFile   string
	ResponseData string
	RequestList  []RequestData
)

type RequestData struct {
	Url      string      `json:"url"`
	Method   string      `json:"method"`
	Protocol string      `json:"protocol"`
	Headers  http.Header `json:"headers"`
	Host     string      `json:"host"`
	Body     string      `json:"body"`
}

func main() {
	//* Init the Request list
	RequestList = []RequestData{}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Value:       "",
				Aliases:     []string{"f"},
				Usage:       "Reads JSON from file returns on request",
				Destination: &InputFile,
			},
			&cli.StringFlag{
				Name:        "data",
				Value:       "",
				Aliases:     []string{"d"},
				Usage:       "Data returned on request",
				Destination: &ResponseData,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "",
				Aliases:     []string{"o"},
				Usage:       "Writes the requests data to the specified filename",
				Destination: &OutputFile,
			},
		},
		Action: func(ctx *cli.Context) error {

			port := ctx.Args().Get(0)
			if num, _ := strconv.Atoi(port); num > 99999 {
				cli.Exit("Port must be in range 0-99999", 2)
			}

			if port == "" {
				port = "5000"
			}

			startAPIServer(port)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func startAPIServer(port string) {
	color.HiYellow("Starting the API endpoint on port: " + port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		//* Checks if user passed -d and -f
		if InputFile != "" && ResponseData != "" {
			color.Red("Cannot use --file and --json at the same time.")
			return
		}
		// * Checks for -f
		if InputFile != "" {
			color.Green("Using file data as response")

			file_content, err := ioutil.ReadFile(InputFile)
			if err != nil {
				log.Panic(err)
			}

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, string(file_content))
				requestDebugger(w, r)
			})
			// * User passed -d
		} else if ResponseData != "" {
			color.Green("Using data argument as response")

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, ResponseData)
				requestDebugger(w, r)
			})
			//* User didn't provide any flag
		} else {
			color.Green("Starting basic API server")
			http.HandleFunc("/", requestDebugger)
		}
		http.ListenAndServe("127.0.0.1:"+port, nil)
	}()

	<-done
	//* Dumps the requests to file on request
	if OutputFile != "" {
		color.HiBlue("\nWriting requests to file...")
		buff, _ := json.Marshal(RequestList)
		ioutil.WriteFile(OutputFile, buff, 0644)
	}
	color.HiRed("\nServer stopped")
}

func requestDebugger(w http.ResponseWriter, r *http.Request) {
	//* Prints and formats the details of the request

	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	request := RequestData{
		Url:      r.URL.Path,
		Method:   r.Method,
		Protocol: r.Proto,
		Headers:  r.Header,
		Host:     r.Host,
		Body:     string(req_body),
	}

	fmt.Println("")
	color.Magenta("Request on: " + color.New(color.FgHiMagenta).Sprint(request.Url))
	color.Yellow("Method: " + color.New(color.FgHiYellow).Sprint(request.Method))
	color.Cyan("Protocol: " + color.New(color.FgHiCyan).Sprint(request.Protocol))

	for i, j := range request.Headers {
		fmt.Print(color.RedString(i) + " : ")
		for _, l := range j {
			fmt.Print(color.HiRedString(l) + "\n")
		}
	}

	color.Green("Body: " + color.New(color.FgHiGreen).Sprint(request.Body))

	RequestList = append(RequestList, request)
}
