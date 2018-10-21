package keybase

import (
	"io/ioutil"
	"os"

	log "github.com/stefancocora/keybasectl/internal/log"
)

var kbdebug bool

// DebugFlag holds the value from the main pkg of the debug flag
type DebugFlag struct {
	Debug bool
}

// NewDebugFlag propagates the value of the debug flag to this package
func (d *DebugFlag) NewDebugFlag(debug bool) {

	d.Debug = debug

	kbdebug = d.Debug

}

// DebugSetting returns the current setting of the debug flag for this pkg
func (d *DebugFlag) DebugSetting() bool {

	return d.Debug
}

func init() {

	// logging setup
	if kbdebug {
		log.LoggingInit(os.Stdout, "short", os.Stderr, "short", os.Stderr, "short")
	} else {
		log.LoggingInit(ioutil.Discard, "short", os.Stderr, "short", ioutil.Discard, "short")
	}

}

// Lookup is used to lookup users and pubkeys using the keybase API
func Lookup(user string) error {

	log.DebugLog.Println("nothing implemented yet")

	// step: lookup user
	errl := lookupUser(user)
	if errl != nil {

		return errl
	}

	return nil
}

// lookupUser uses the keybase API to lookup the given user
func lookupUser(user string) error {

	return nil
}
