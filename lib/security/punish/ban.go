package security

import (
	"lib/misc"
	"lib/synchro"
)

// Ban : Processes banning logic for IP Address
func Ban(ipAddr string) {
	// Spawning a Goroutine to handle this in the background
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		// No need to lock here since we're not interacting with
		// any databases, and just launching terminal commands.

		// First, banning the provided IP Address using iptables.
		misc.ExecuteUsingShell("iptables -A INPUT -s " + ipAddr + " -j DROP")

		// Then, launching the Auto-Unban function.
		UnbanAfterThirtyMins(ipAddr)

		// Finally, returning to end this Goroutine
		return

	}, map[string]interface{}{})
}
