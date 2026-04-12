package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(`
    ███████╗ ██████╗ ██╗  ██╗███████╗██╗     ███████╗████████╗ ██████╗ ███╗   ██╗
    ██╔════╝██╔═══██╗██║ ██╔╝██╔════╝██║     ██╔════╝╚══██╔══╝██╔═══██╗████╗  ██║
    ███████╗██║   ██║█████╔╝ █████╗  ██║     █████╗     ██║   ██║   ██║██╔██╗ ██║
    ╚════██║██║   ██║██╔═██╗ ██╔══╝  ██║     ██╔══╝     ██║   ██║   ██║██║╚██╗██║
    ███████║╚██████╔╝██║  ██╗███████╗███████╗███████╗   ██║   ╚██████╔╝██║ ╚████║
    ╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═══╝

    [INFO]  Go Skeleton Server v0.0.1 — "It's not a bug, it's a feature... of laziness"
    [WARN]  This is 100% skeleton. No real server was harmed (or started) in the making of this log.
	`)

	fakeSteps := []string{
		"Waking up the Gopher...",
		"Feeding him virtual carrots...",
		"Convincing the goroutines to actually run (they're on strike)",
		"Checking if port :8080 exists in this dimension...",
		"Loading absolutely nothing into memory...",
		"Searching for real code... (404 - Not Found)",
		"Compiling dreams into reality...",
		"Almost there... just pretending to bind the listener...",
	}

	for _, step := range fakeSteps {
		fmt.Printf("    [STARTUP] %s\n", step)
		time.Sleep(800 * time.Millisecond)
	}

	fmt.Println(`
    [OK]  Fake server "started" successfully on :666 (the number of the skeleton)
          (You see anything? Cuz the Gopher doesn't.)

          Current status: Running on pure vibes and copium.
          Real implementation goes here → when you stop procrastinating.
          Blocking until you decide to press Ctrl+C to "shut down" the non-existent server.
	`)

	for {
	}
}
