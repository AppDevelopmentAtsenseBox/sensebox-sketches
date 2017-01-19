package main

import (
	// "encoding/json"
	"fmt"
	"os"

	"github.com/honeybadger-io/honeybadger-go"
)

const envPrefix = "SENSEBOX_SKETCHES_"

var ConfigCaCertBytes, ConfigServerCertBytes, ConfigServerKeyBytes []byte

// var DockerImage, ContainerName string

func initConfigFromEnv() {
	errors := make([]error, 0)

	// try to configure honeybadger integration..
	honeybadgerApiKey, _ := getStringFromEnv("HONEYBADGER_APIKEY")
	if honeybadgerApiKey != "" {
		honeybadger.Configure(honeybadger.Configuration{APIKey: honeybadgerApiKey})
		LogInfo("startup", "enabled honeybadger integration")
	}

	caCertBytes, caCertBytesErr := getBytesFromEnv("CA_CERT")
	if caCertBytesErr != nil {
		errors = append(errors, caCertBytesErr)
	}

	serverCertBytes, serverCertBytesErr := getBytesFromEnv("SERVER_CERT")
	if serverCertBytesErr != nil {
		errors = append(errors, serverCertBytesErr)
	}

	serverKeyBytes, serverKeyBytesErr := getBytesFromEnv("SERVER_KEY")
	if serverKeyBytesErr != nil {
		errors = append(errors, serverKeyBytesErr)
	}

	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	ConfigCaCertBytes = caCertBytes
	ConfigServerCertBytes = serverCertBytes
	ConfigServerKeyBytes = serverKeyBytes
}
