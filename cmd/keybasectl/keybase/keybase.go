/*
Copyright 2015 All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keybase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/stefancocora/keybasectl/internal/log"
)

// ErrorUserNotFound is the error returned when the user isn't found
type ErrorUserNotFound struct {
	err    error
	errmsg string
}

// Error implements the error interface for a type of ErrorUserNotFound
func (unf ErrorUserNotFound) Error() string {
	return unf.errmsg
}

// ErrorPKNotFound is the error returned when the user's public key isn't found
type ErrorPKNotFound struct {
	err    error
	errmsg string
}

// Error implements the error interface for a type of ErrorPKNotFound
func (pknf ErrorPKNotFound) Error() string {
	return pknf.errmsg
}

var kbdebug bool

var keybaseUserLookupURL = "https://keybase.io/_/api/1.0/user/lookup.json?usernames="

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

// Status contains the API call status results from keybase.io.
type Status struct {
	Desc string `json:"desc"`
	Code int    `json:"code"`
	Name string `json:"name"`
}

// User contains information regarding a user coming from the "them" response from the keybase API
type User struct {
	ID     string `json:"id"`
	Basics Basics `json:"basics"`
	// Invitations InvitationStats `json:"invitation_stats"`
	// Profile     Profile         `json:"profile"`
	// Emails      Emails          `json:"emails"`
	// PublicKeys  map[string]*Key `json:"public_keys"`
	// PrivateKeys map[string]*Key `json:"private_keys"`
}

// Basics contain basic information about the user.
type Basics struct {
	Username string `json:"username_cased"`
	// Created      int    `json:"ctime"`
	// Modified     int    `json:"mtime"`
	// IDVersion    int    `json:"id_version"`
	TrackVersion int `json:"track_version"`
	// LastIDChange int    `json:"last_id_change"`
	Salt string `json:"salt"`
}

// A KeyType is used to denote whether a key is public or private.
type KeyType int

// These constants provide friendly names for the key types returned by
// the API.
const (
	PublicKey  KeyType = 1
	PrivateKey KeyType = 2
)

// A Key contains information about a public or private key.
type Key struct {
	KeyID       string  `json:"kid"`
	Fingerprint string  `json:"key_fingerprint"`
	KeyType     KeyType `json:"key_type"`
	// Bundle      string  `json:"bundle"`
	// Modified    int     `json:"mtime"`
	// Created     int     `json:"ctime"`
}

func init() {

	// logging setup
	if kbdebug {
		log.LoggingInit(os.Stdout, "short", os.Stderr, "short", os.Stderr, "short")
	} else {
		log.LoggingInit(ioutil.Discard, "short", os.Stderr, "short", ioutil.Discard, "short")
	}

}

// UserLookup is used to lookup users using the keybase API
func UserLookup(username []string) ([]string, []string, error) {

	log.DebugLog.Printf("lookup username(s): %v", username)
	var uf, unf []string

	// step: lookup username
	uf, unf, errl := lookupUser(username)
	if errl != nil {

		if unfe, ok := errl.(ErrorUserNotFound); ok {

			log.DebugLog.Printf("received a ErrorUserNotFound error: %v", unfe)
		}
		return uf, unf, errl
	}

	return uf, unf, nil
}

// lookupUser uses the keybase API to lookup the given user
func lookupUser(username []string) ([]string, []string, error) {

	var userResponse struct {
		Status *Status `json:"status"`
		User   []*User `json:"them"`
	}
	var empty []string
	var userFound []string
	var userNotFound []string

	uname := strings.Join(username, ",")

	url := fmt.Sprintf("%s%s&fields=basics", keybaseUserLookupURL, uname)
	log.DebugLog.Printf("targeting keybase API url: %s", url)
	res, errlu := http.Get(url)
	if errlu != nil {

		return empty, empty, errlu
	}

	respb, errRA := ioutil.ReadAll(res.Body)
	log.DebugLog.Printf("response body: %s", respb)
	if errRA != nil {

		return empty, empty, errRA
	}

	defer res.Body.Close()
	errDec := json.Unmarshal(respb, &userResponse)

	if errDec != nil {

		return empty, empty, errDec
	}

	for u := range userResponse.User {

		if userResponse.User[u] != nil {

			log.DebugLog.Printf("user %s found", username[u])
			userFound = append(userFound, username[u])
			// log.DebugLog.Printf("unmarshalled resp user: %v", userResponse.User[u].Basics.Username)
			// log.DebugLog.Printf("unmarshalled resp salt: %v", userResponse.User[u].Basics.Salt)
			// log.DebugLog.Printf("unmarshalled resp track_version: %v", userResponse.User[u].Basics.TrackVersion)

		} else {

			msg := fmt.Sprintf("user %s not found", username[u])
			log.DebugLog.Printf(msg)
			userNotFound = append(userNotFound, username[u])
		}
	}

	if len(userNotFound) == 0 {

		return userFound, empty, nil
	} else {

		msg := fmt.Sprintf("user(s) %v not found", userNotFound)
		var eunf ErrorUserNotFound
		eunf.errmsg = msg
		return userFound, userNotFound, eunf
	}
}

// PubKeyLookup is used to lookup pubkeys using the keybase API
func PubKeyLookup(username []string) ([]string, []string, error) {

	log.DebugLog.Printf("lookup pubkey for username(s): %v", username)
	var kf, knf []string

	// step: lookup username's pubkey
	kf, knf, errl := lookupPubKey(username)
	if errl != nil {

		if pknfe, ok := errl.(ErrorPKNotFound); ok {

			log.DebugLog.Printf("received a ErrorPKNotFound error: %v", pknfe)
		}
		return kf, knf, errl
	}

	return kf, knf, nil
}

// lookupPubKey uses the keybase API to lookup the given user's pubkey
func lookupPubKey(username []string) ([]string, []string, error) {

	var pubKeyResponse struct {
		Status *Status `json:"status"`
		Key    []*Key  `json:"them"`
	}
	var empty, pubKeyFound, pubKeyNotFound []string

	uname := strings.Join(username, ",")

	url := fmt.Sprintf("%s%s&fields=public_keys", keybaseUserLookupURL, uname)
	log.DebugLog.Printf("targeting keybase API url: %s", url)
	res, errlu := http.Get(url)
	if errlu != nil {

		return empty, empty, errlu
	}

	respb, errRA := ioutil.ReadAll(res.Body)
	// log.DebugLog.Printf("response body: %s", respb)
	if errRA != nil {

		return empty, empty, errRA
	}

	defer res.Body.Close()
	errDec := json.Unmarshal(respb, &pubKeyResponse)

	if errDec != nil {

		return empty, empty, errDec
	}

	for u := range pubKeyResponse.Key {

		if pubKeyResponse.Key[u] != nil {

			log.DebugLog.Printf("public key for user %s found", username[u])
			pubKeyFound = append(pubKeyFound, username[u])
			// log.DebugLog.Printf("unmarshalled resp user: %v", pubKeyResponse.Key[u].Basics.Username)
			// log.DebugLog.Printf("unmarshalled resp salt: %v", pubKeyResponse.Key[u].Basics.Salt)
			// log.DebugLog.Printf("unmarshalled resp track_version: %v", pubKeyResponse.Key[u].Basics.TrackVersion)

		} else {

			msg := fmt.Sprintf("public key for user %s not found", username)
			log.DebugLog.Printf(msg)
			pubKeyNotFound = append(pubKeyNotFound, username[u])
		}
	}

	if len(pubKeyNotFound) == 0 {

		return pubKeyFound, empty, nil
	} else {

		msg := fmt.Sprintf("public key for user(s) %v not found", pubKeyNotFound)
		var eunf ErrorPKNotFound
		eunf.errmsg = msg
		return pubKeyFound, pubKeyNotFound, eunf
	}

}
