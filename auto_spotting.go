// Copyright (c) 2016-2022 Cristian Măgherușan-Stanciu
// Licensed under the Open Software License version 3.0

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	autospotting "github.com/AutoSpotting/AutoSpotting/core"
	"github.com/aws/aws-lambda-go/lambda"
)

var as *autospotting.AutoSpotting
var conf autospotting.Config

// Version represents the build version being used
var Version = "number missing"

// SavingsCut is populated at build time and controls the percentage of the savings charged for the stable builds
var SavingsCut = "0"

var eventFile string

func main() {
	eventFile = conf.EventFile

	if autospotting.RunningFromLambda() {
		lambda.Start(Handler)
	} else if eventFile != "" {
		parseEvent, err := ioutil.ReadFile(eventFile)
		if err != nil {
			log.Fatal(err)
		}
		Handler(context.TODO(), parseEvent)
	} else {
		eventHandler(nil)
	}
}

func eventHandler(event *json.RawMessage) {

	log.Println("Starting autospotting agent, build ", Version, "charging", SavingsCut, "percent of savings via AWS Marketplace")

	log.Printf("Configuration flags: %#v", conf)

	as.EventHandler(event)
	log.Println("Execution completed, nothing left to do")
}

// this is the equivalent of a main for when running from Lambda, but on Lambda
// the runFromCronEvent() is executed within the handler function every time we have an event
func init() {
	as = &autospotting.AutoSpotting{}

	sc, err := strconv.ParseFloat(SavingsCut, 64)
	if err != nil {
		log.Printf("Failed to convert savings cut %s to float\n", SavingsCut)
	}

	conf = autospotting.Config{
		Version:    Version,
		SavingsCut: sc,
	}

	autospotting.ParseConfig(&conf)
	as.Init(&conf)
}

// Handler implements the AWS Lambda handler interface
func Handler(ctx context.Context, rawEvent json.RawMessage) {
	eventHandler(&rawEvent)
}