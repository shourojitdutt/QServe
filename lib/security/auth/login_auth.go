package security

import (
	"lib/constants"
	"lib/database"
	"lib/synchro"
)

// AuthenticateUserLogin : Checks user credentials and returns map[string]interface{}
// so that it can appropriately send errors.
func AuthenticateUserLogin(appAuth string, username string, password string, ipAddr string) map[string]interface{} {
	// Locking to prevent dirty reads
	// because I'm a good Software Engineer.
	// Remove this if the use-case scenario
	// is devoid of dirty reads in all the
	// legitimate use-case scenarios.
	synchro.Lock()
	var retVal = map[string]interface{}{"success": true}
	if !PerformCoreAuth(appAuth, ipAddr) {
		// Invalid appAuth
		retVal["success"] = false
		retVal["error"] = "Nice Try Hacker Boi!"
		synchro.Unlock()
		return retVal
	}
	var usersCollection = database.GetCollection(constants.AppDbName, constants.UserCollectionName)
	if !database.CheckExistsInCollection(usersCollection, "username", username) {
		// User does not exist
		retVal["success"] = false
		retVal["error"] = "User does not exist!"
		synchro.Unlock()
		return retVal
	}
	var fetchedUser = database.SearchInCollection(usersCollection, "username", username)
	var pwSalt = fetchedUser.Lookup("password_salt").String()
	if pwSalt == "" {
		// No password salt exists for this user.
		// Shouldn't occur, but if it does, then
		// it indicates an issue with the database entry.
		retVal["success"] = false
		// Showing internal server error here because
		// we don't want any regular user to see
		// that we somehow screwed up the password salt
		retVal["error"] = "Internal Server Error"
		synchro.Unlock()
		return retVal
	}
	var hashedPassword = HashPassword(password, pwSalt)
	var receivedUsername = fetchedUser.Lookup("username").String()
	var receivedPassword = fetchedUser.Lookup("password").String()
	synchro.Unlock()
	if (receivedUsername == username) && (hashedPassword == receivedPassword) {
		return retVal
	}
	retVal["success"] = false
	retVal["error"] = "Invalid Credentials!"
	return retVal

}
