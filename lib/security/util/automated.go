package security

import (
	"fmt"
	"lib/constants"
	"lib/database"
	"lib/synchro"
	"lib/util"
	"time"
)

// ClearNumRejections : Clears number of rejections in table
func ClearNumRejections(ipAddr string) {
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		// Locking here because Software Engineering
		synchro.Lock()
		var collection = database.GetCollection(constants.AppDbName, constants.RejectedIPsCollectionName)
		if database.CheckExistsInCollection(collection, "ipAddr", ipAddr) {
			// Update if exists
			var entry = database.SearchInCollection(collection, "ipAddr", ipAddress)
			var id = entry.Lookup("_id").ObjectID()
			// Clearing out numRejections
			var numRejections = 0
			database.UpdateOneInCollection(collection, id, "numRejections", numRejections)
		}
		synchro.Unlock()
		return
	}, map[string]interface{}{})
}

// UnbanAfterThirtyMins : Unbans an IP After 30 mins
func UnbanAfterThirtyMins(ipAddr string) {
	// Spawning a Goroutine that unblocks in 30 mins
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		addr := fmt.Sprint(data["ipAddr"])
		timer := time.NewTimer(30 * time.Minute)
		<-timer.C
		clearBan(addr)
	}, map[string]interface{}{"ipAddr": util.BasicToInterface(ipAddr)})
	// Returning to give control back and avoid any
	// potential issues
	return
}

// Utility function that clears the rejection
// table that we can use with the reccuring
// rejection clearer to make the code easier to
// read and understand.
func clearRejectionTable() {
	var collection = database.GetCollection(constants.AppDbName, constants.RejectedIPsCollectionName)
	database.ClearCollection(collection)
	return
}

// RecurringRejectionClearer : Every 30 mins, clears out the rejectedIpAddresses table.
func RecurringRejectionClearer() {
	// Spawning a Goroutine to handle this since
	// we cannot afford to have this on the main thread.
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		for {
			// Wait 30 mins
			<-time.After(30 * time.Minute)
			// Clear the rejections table
			clearRejectionTable()
		}
	}, map[string]interface{}{})
}

func clearBan(ipAddr string) {
	// Spawning a Goroutine again because this should be
	// performed in the background.
	synchro.SpawnGoroutine(func(data map[string]interface{}) {

		// Unbanning using iptables
		util.ExecuteUsingShell("iptables -D INPUT -s " + ipAddr + " -j DROP")

	}, map[string]interface{}{})
}
