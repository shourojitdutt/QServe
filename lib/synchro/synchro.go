package synchro

import (
	"runtime"
	"sync"
)

/*
	@author Shourojit Dutt

	SYNCHRO LIBRARY
	–––––––––––––––
	The whole idea with this library is to
	maintain the synchronization and paralellism
	of the entire backend in one place.
	This enables better security and stability
	with regards to the manipulation of the
	WaitGroup, the overloading of which will
	cause the server to crash. It also enables
	for better code readibility since the whole
	synchronization logic and code are all in one
	centralized location.

*/

// Defining function type to be able to
// use it below.
type function func(map[string]interface{})

var globalMutex sync.Mutex
var globalWaitGroup sync.WaitGroup

// SpawnGoroutine : The crux of this library
func SpawnGoroutine(inputFunc function, args map[string]interface{}) {
	/*
		Waiting for the current server backlog to finish
		in the event that more than 1048000 or more
		goroutines exist since the hard limit for Golang
		is 1048576 and we want to leave some breathing room.
	*/
	if runtime.NumGoroutine() >= 1048000 {
		globalWaitGroup.Wait()
	}
	globalWaitGroup.Add(1)
	go inputFunc(args)
}

// Mutex functions for when DB write is happening.

// Lock : Global Mutex Lock
func Lock() {
	globalMutex.Lock()
}

// Unlock : Global Mutex Unlock
func Unlock() {
	globalMutex.Unlock()
}
