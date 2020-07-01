package security

import (
	punish "lib/security/punish"
)

// The main authentication key for the whole server.
// TODO : Edit Later obviously; this is shit.
const appAuth = "abcd"

// checkAuth : Core check for proper user authentication.
func checkAuth(auth string) bool {
	return auth == appAuth
}

// PerformCoreAuth : Performs core authentication to check
// the legitimacy of the request, and then performs rejection
// logic as well.
func PerformCoreAuth(auth string, ipAddr string) bool {
	// If Auth checks out, all is good
	if checkAuth(auth) {
		return true
	}
	// If not, performing rejection logic
	punish.Reject(ipAddr)
	return false
}
