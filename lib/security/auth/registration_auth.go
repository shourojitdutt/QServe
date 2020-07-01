package security

import (
	"lib/constants"
	"lib/database"
	"lib/synchro"

	"go.mongodb.org/mongo-driver/bson"
)

// AuthenticateUserRegistration : Checks basic authentication and then registers the user.
// Sending response as a map[string]interface{} so that the errors can be logged and sent
// back to the user appropriately.
func AuthenticateUserRegistration(appAuth string, username string, password string, ipAddr string) map[string]interface{} {
	// Locking to prevent dirty writes.
	// Not reccomended to remove this since
	// no real use-case scenarios would be
	// immune to dirty-writes in legitimate
	// use-case scenarios.
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
	if database.CheckExistsInCollection(usersCollection, "username", username) {
		// User already exists
		retVal["success"] = false
		retVal["error"] = "User already exists!"
		synchro.Unlock()
		return retVal
	}
	// Generating a salt
	var salt = GenerateSalt()
	var hashedPw = HashPassword(password, salt)
	var data = bson.M{"username": username, "password": hashedPw, "password_salt": salt}
	// Cannot directly return database.InsertOneIntoCollection(...) because
	// of the Lock, so putting the result into a variable and returning that
	// after the operation has finished executing instead.
	var res = database.InsertOneIntoCollection(usersCollection, data)
	synchro.Unlock()
	if !res {
		retVal["success"] = false
		// Showing internal server error here since
		// something would have gone wrong with Mongo
		// for this to happen.
		retVal["error"] = "Internal Server Error"
		return retVal
	}
	return retVal
}
