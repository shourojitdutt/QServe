package security

import (
	"lib/constants"
	"lib/database"
	"math/rand"
	"strings"
	"time"
)

// TwoFAGen : Generates OTP
func TwoFAGen(username string) string {
	var nums = "0123456789"
	var otp strings.Builder
	var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := 0; index < 6; index++ {
		otp.WriteString(string(nums[randSeed.Intn(len(nums)-1)]))
	}
	writeOtpToDb(username, otp.String())
	return otp.String()
}

// Puts OTP in DB, associated to that particular user.
func writeOtpToDb(username string, otp string) {
	var collection = database.GetCollection(constants.AppDbName, constants.UserCollectionName)
	var currentUser = database.SearchInCollection(collection, "username", username)
	database.UpdateOneInCollection(collection, currentUser.Lookup("_id").ObjectID(), "otp", otp)
}

// PerformTwoFA : Validates entered OTP
func PerformTwoFA(username string, enteredOtp string) bool {
	var collection = database.GetCollection(constants.AppDbName, constants.UserCollectionName)
	var currentUser = database.SearchInCollection(collection, "username", username)
	if currentUser.Lookup("otp").String() == enteredOtp {
		return true
	}
	return false
}
