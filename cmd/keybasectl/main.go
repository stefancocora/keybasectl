package main

import (
	"log"

	"github.com/stefancocora/keybasectl/internal/version"
)

func main() {

	log.Println("starting engines")

	bc, err := version.BuildContext()
	if err != nil {

		log.Fatalf("[FATAL] unable to get the binary version: %v", err)
	}

	log.Printf("build environment: %s", bc)

	bccli, err := version.BuildContextCli()
	if err != nil {

		log.Fatalf("[FATAL] unable to get the binary version: %v", err)
	}

	log.Printf("build environment: %s\n", bccli)

	log.Println("stopping engines, we're done")

}
