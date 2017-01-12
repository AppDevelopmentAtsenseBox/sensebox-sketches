package main

import (
	"github.com/honeybadger-io/honeybadger-go"
)

type senseBoxSketchesServer struct {}

func main () {

	defer honeybadger.Monitor()

	initConfigFromEnv()

	sketches := senseBoxSketchesServer{}
	sketches.StartHTTPSServer()
}