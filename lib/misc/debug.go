package misc

import "log"

// Must : Logs error if one occurs
func Must(err error) bool {
	if err != nil {
		// Logging an error in case it exists.
		// Don't want to panic and stop the server.
		log.Println("[!] Error occurred : ", err)
	}
	return err != nil
}
