package security

import (
	"lib/constants"
	"lib/database"
	"lib/synchro"

	"go.mongodb.org/mongo-driver/bson"
)

// Reject : Invokes rejection logic for IP Address
func Reject(ipAddress string) {
	// Spawning a Goroutine since concurrency will be needed here
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		// Locking since we're gonna be writing to DB
		synchro.Lock()
		// Getting collection for rejection and writing this IP to it.
		var collection = database.GetCollection(constants.AppDbName, constants.RejectedIPsCollectionName)

		if database.CheckExistsInCollection(collection, "ipAddr", ipAddress) {
			// Update if exists
			var entry = database.SearchInCollection(collection, "ipAddr", ipAddress)
			var id = entry.Lookup("_id").ObjectID()
			var numRejections = entry.Lookup("numRejections").Int64() + 1
			if numRejections >= 10 {
				synchro.Unlock()
				Ban(ipAddress)
				ClearNumRejections(ipAddress)
				return
			}
			database.UpdateOneInCollection(collection, id, "numRejections", numRejections)
			synchro.Unlock()
			return
		}
		// Insert if it doesn't
		database.InsertOneIntoCollection(collection, bson.M{"ipAddr": ipAddress, "numRejections": 1})
		synchro.Unlock()
		// Now Unlocking after the writing is done.
	}, map[string]interface{}{})
}
