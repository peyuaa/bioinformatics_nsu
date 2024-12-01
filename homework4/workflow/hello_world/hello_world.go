package main

import (
	sp "github.com/scipipe/scipipe"
)

func main() {
	// Init workflow with a name, and max concurrent tasks
	wf := sp.NewWorkflow("hello_world", 4)

	// Initialize processes and set output file paths
	hello := wf.NewProc("hello", "echo 'Hello ' > {o:out}")
	hello.SetOut("out", "hello.txt")

	world := wf.NewProc("world", "echo $(cat {i:in}) World >> {o:out}")
	world.SetOut("out", "{i:in|%.txt}_world.txt")

	// Connect network
	world.In("in").From(hello.Out("out"))

	// Run workflow
	wf.Run()
}
