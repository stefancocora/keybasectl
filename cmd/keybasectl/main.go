package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"strings"

	"github.com/stefancocora/keybasectl/cmd/keybasectl/keybase"
	log "github.com/stefancocora/keybasectl/internal/log"
	"github.com/stefancocora/keybasectl/internal/version"
)

var debug bool

//---

// userFlag is the struct that get populated when the --auth cli flag is provided
type userFlag struct {
	set   bool
	value []string
}

func (us *userFlag) Set(val string) error {

	us.value = strings.Split(val, ",")
	us.set = true
	return nil
}

func (us *userFlag) String() string {

	return fmt.Sprintf("%v", us.value)
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
	var errl, errpkl error
	var kbFl keybase.DebugFlag
	var uf, unf []string // captures the users found and not found
	var kf, knf []string // captures the user's pubkey found and not found

	if !flag.Parsed() {

		flag.Parse()
	}

	// logging setup
	if debug {
		log.LoggingInit(os.Stdout, "short", os.Stderr, "short", os.Stderr, "short")
	} else {
		log.LoggingInit(ioutil.Discard, "short", ioutil.Discard, "short", ioutil.Discard, "short")
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
		goto exitAll
	}

	log.DebugLog.Printf("--user flag arguments: %#v", flag.Args())

	kbFl.NewDebugFlag(debug)
	log.DebugLog.Printf("current setting for the debug flag inside the keybase pkg: %v", kbFl.DebugSetting())

	// step: lookup user against keybase
	uf, unf, errl = keybase.UserLookup(usfL.value)
	if errl != nil {

		exitVal++
		if _, ok := errl.(keybase.ErrorUserNotFound); ok {

			if len(uf) > 0 {
				fmt.Fprintf(os.Stdout, "user(s): %v found during keybase lookup\n", uf)
			}
			fmt.Fprintf(os.Stdout, "user(s): %v not found during keybase lookup\n", unf)
			log.ErrorLog.Printf("error during keybase user lookup: %s", errl.Error())
			goto exitAll
		} else {

			fmt.Fprintf(os.Stdout, "error : %s\n", errl.Error())
			log.ErrorLog.Printf("error : %s", errl.Error())
			goto exitAll
		}
	} else {
		fmt.Fprintf(os.Stdout, "user(s): %v found during keybase lookup\n", uf)
	}

	// step: lookup user's pubkey against keybase
	kf, knf, errpkl = keybase.PubKeyLookup(usfL.value)
	if errpkl != nil {

		exitVal++
		if _, ok := errpkl.(keybase.ErrorPKNotFound); ok {

			if len(kf) > 0 {

				fmt.Fprintf(os.Stdout, "user(s): %v public key found during keybase public key lookup\n", kf)
			}
			fmt.Fprintf(os.Stdout, "user(s): %v public key not found during keybase public key lookup\n", knf)
			log.ErrorLog.Printf("error during keybase public key lookup: %s", errpkl.Error())
			goto exitAll
		} else {

			fmt.Fprintf(os.Stdout, "error : %s\n", errpkl.Error())
			log.ErrorLog.Printf("error : %s", errpkl.Error())
			goto exitAll
		}
	} else {
		fmt.Fprintf(os.Stdout, "user(s): %v public key found during keybase public key lookup\n", kf)

	}

exitAll:
	if exitVal > 0 {

		log.InfoLog.Println("stopping engines, we're done")
		os.Exit(1)
	} else {

		log.InfoLog.Println("stopping engines, we're done")
	}

}
