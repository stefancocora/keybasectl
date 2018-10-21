package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/stefancocora/keybasectl/cmd/keybasectl/keybase"
	log "github.com/stefancocora/keybasectl/internal/log"
	"github.com/stefancocora/keybasectl/internal/version"
)

var debug bool

//---

// userFlag is the struct that get populated when the --auth cli flag is provided
type userFlag struct {
	set   bool
	value string
}

func (us *userFlag) Set(val string) error {

	us.value = val
	us.set = true
	return nil
}

func (us *userFlag) String() string {

	return us.value
}

var usfL userFlag
var usEnv = "KEYBASECTL_USER"
var usUsage = fmt.Sprintf("User to lookup. Alternatively sourced from %s [required]", usEnv)
var usName = "user"

//---

// apiEndpointFlag is the struct that get populated when the --api cli flag is provided
// this switches the keybase endpoint to either their prod or staging API endpoints
type apiEndpointFlag struct {
	set   bool
	value string
}

func (us *apiEndpointFlag) Set(val string) error {

	us.value = val
	us.set = true
	return nil
}

func (us *apiEndpointFlag) String() string {

	return us.value
}

var apifL apiEndpointFlag
var apiEnv = "KEYBASECTL_API_ENDPOINT"
var apiUsage = fmt.Sprintf("Keybase API endpoint to target. Alternatively sourced from %s. Default to [production]. Disabled for now since the staging URL from the keybase client doesn't DNS resolve", apiEnv)
var apiName = "api"

//---

func init() {

	flag.BoolVar(&debug, "debug", false, "turn on or off debugging")
	flag.Var(&apifL, apiName, apiUsage)
	flag.Var(&usfL, usName, usUsage)

}

func main() {

	var exitVal = 0

	if !flag.Parsed() {

		flag.Parse()
	}

	// logging setup
	if debug {
		log.LoggingInit(os.Stdout, "short", os.Stderr, "short", os.Stderr, "short")
	} else {
		log.LoggingInit(ioutil.Discard, "short", os.Stderr, "short", ioutil.Discard, "short")
	}

	log.InfoLog.Println("starting engines")

	_, err := version.BuildContext()
	if err != nil {

		stdlog.Fatalf("[FATAL] unable to get the binary version: %v", err)
	}

	// step: check required flag/envvar
	use, okOaEnv := os.LookupEnv(usEnv)
	log.DebugLog.Printf("environment variable lookup result: %s", use)
	log.DebugLog.Printf("cli flag: %s set to: %s, set: %v", usName, usfL.value, usfL.set)
	if !usfL.set && !okOaEnv {

		log.ErrorLog.Printf("required flag or environment variable not set! flag: %s, environmentVariable: %v", usName, usEnv)
		fmt.Fprintf(os.Stdout, "required flag or environment variable not set! flag: \"%s\", environmentVariable: \"%v\"\n", usName, usEnv)
		exitVal++
	}

	// step: lookup user against keybase
	kbFl := new(keybase.DebugFlag)
	kbFl.NewDebugFlag(debug)
	log.DebugLog.Printf("current setting for the debug flag inside the keybase pkg: %v", kbFl.DebugSetting())

	errl := keybase.Lookup(usfL.value)
	if errl != nil {

		exitVal++
		if _, ok := errl.(keybase.ErrorUserNotFound); ok {

			fmt.Fprintf(os.Stdout, "error during keybase user lookup: %s\n", errl.Error())
			log.ErrorLog.Printf("error during keybase user lookup: %s", errl.Error())
		} else {

			fmt.Fprintf(os.Stdout, "error : %s\n", errl.Error())
			log.ErrorLog.Printf("error : %s", errl.Error())
		}
	}

	log.InfoLog.Println("stopping engines, we're done")

	if exitVal > 0 {

		os.Exit(1)
	}

}
