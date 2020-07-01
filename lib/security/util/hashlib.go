package security

import (
	"crypto"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateSalt : Generates randomly created, (hopefully) unique salt
func GenerateSalt() string {
	var potentialCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()-_=+[]{};:,.<>?"
	var retVal strings.Builder
	var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generating a 96-Character Salt
	for index := 0; index < 96; index++ {
		retVal.WriteString(string(potentialCharacters[randSeed.Intn(len(potentialCharacters)-1)]))
	}
	return retVal.String()
}

// Salts password
func saltPassword(inputString string, existingSalt string) string {
	var salt string
	if existingSalt != "" {
		salt = existingSalt
	} else {
		// Generating 64-character Salt
		// since one does not currently
		// exist.
		salt = GenerateSalt()
	}
	// Will hash at index 0, -1, len(pass)/2 and len(pass)/3
	// This way, the positions will be different for each password
	// depending on the length of the password.
	// Using a map to get the position check logic to perform in O(1)
	var hashPositions = map[int64]bool{0: true, -1: true, int64(len(inputString) / 2): true, int64(len(inputString) / 3): true}
	var retVal strings.Builder
	for index := int64(0); index < int64(len(inputString)); index++ {
		if hashPositions[index] {
			retVal.WriteString(string(inputString[index]) + salt)
		} else {
			retVal.WriteString(string(inputString[index]))
		}
	}
	return retVal.String()
}

// HashPassword : Using SHA3-512 as the hashing standard since it's strong.
// Used when a user creates a new account or changes password.
// Input the value of saltForThisPassword as "" for new user, or whatever
// the existing value is for an existing user.
func HashPassword(inputString string, saltForThisPassword string) string {
	var saltedPassword string
	if saltForThisPassword == "" {
		saltedPassword = saltPassword(inputString, "")
	} else {
		saltedPassword = saltPassword(inputString, saltForThisPassword)
	}
	var hashObj = crypto.SHA3_512.New()
	hashObj.Write([]byte(saltedPassword))
	return fmt.Sprintf("%x", hashObj.Sum(nil))
}
