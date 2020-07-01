package security

import (
	"fmt"
	"lib/constants"
	"lib/database"
	"lib/util"
	"time"
)

/*
	The request number sent from the client side must be
	the number of seconds since Jan 1st 1970 in order to ensure
	unique numbers each time. This will cause limitations in case
	the user decides to go abroad to a slower time-region, but oh well.
*/

// This is the time format for our server. Officially known as the RFC3339 format.
const timeFormat = "2006-01-02T15:04:05Z07:00"

// This is the time-limit beyond which we do not allow requests to live.
const timeLimit = 60 * time.Second

// requestNumberCheck : Checks if the request number in
// the user profile is less than the currently received number.
func requestNumberCheck(requestPacket map[string]interface{}) bool {
	var collection = database.GetCollection(constants.AppDbName, constants.UserCollectionName)
	var user = database.SearchInCollection(collection, "username", requestPacket["username"])
	if user.Lookup("requestNumber").String() < fmt.Sprint(requestPacket["requestNumber"]) {
		return true
	}
	return false
}

// requestTimeStampCheck : Checks timestamp of request to
// ensure that no request that has exceeded it's time-limit
// of <timeLimit> seconds can be reused to DDoS effectively.
func requestTimeStampCheck(requestPacket map[string]interface{}) bool {
	reqTimeStamp, reqErr := time.Parse(timeFormat, fmt.Sprint(requestPacket["requestTimeStamp"]))
	if !util.Must(reqErr) {
		/*
			Either this is actually difficult to process or my mind is shot to hell right now,
			but in any case, what this comparison is doing is saying that if we take the timestamp
			that's sent with the request, and add the timeLimit to it, that should exceed the current time;
			This holds true with our logic since to exceed the current time, the request would have to be
			within it's designated time frame of arrival and execution. For example, if a request was sent
			at 17:55:44PM on 21st Aug 2019, then by the current timeLimit (1 minute) at the moment of writing this comment,
			the request should expire at 17:56:44PM, so if it is received at 17:56:22PM, the request timestamp (17:55:44) plus
			the timeLimit (1 min) creates 17:56:44PM which is a greater time than the current time (17:56:22).
			The reason this should work in the first place is that we'll be using HTTPS for all our traffic, and possibly
			even doing a client-side encryption, and server-side decryption of any sent packets in order to ensure that the sent
			data is encrypted so that nobody can mess around with the timestamps, albeit we don't exclusively rely on timestamps anyway.
		*/
		if reqTimeStamp.Add(timeLimit).UnixNano() >= time.Now().UnixNano() {
			return true
		}
	}
	return false
}

// InvokeTwoProngedSniffingDefense : Invokes our two functions
// that defend against sniffing-based attacks. This function will reject
// malicious or suspicious packets right away, and return a boolean according
// to whether the packet is safe or not. In the event that this function returns
// "false", use the return keyword in whatever function is calling this one in order
// to avoid any unnecessary problems.
func InvokeTwoProngedSniffingDefense(requestPacket map[string]interface{}) bool {
	// Let's start off assuming that it's safe because why not
	var packetIsSafe = true
	if !requestNumberCheck(requestPacket) || !requestTimeStampCheck(requestPacket) {
		// Rejecting Packet
		Reject(fmt.Sprint(requestPacket["ipAddr"]))
		// Now, marking packet as unsafe.
		packetIsSafe = false
	}
	return packetIsSafe
}
