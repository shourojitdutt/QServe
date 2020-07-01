package misc

import (
	"lib/synchro"
	"os/exec"
)

// ExecuteUsingShell : Executes shell commands
func ExecuteUsingShell(command string) {
	// Spawning Goroutine for this because concurrency is almost always a good idea
	synchro.SpawnGoroutine(func(data map[string]interface{}) {
		exec.Command("sh", "-c", command).Run()
	}, map[string]interface{}{})
}
