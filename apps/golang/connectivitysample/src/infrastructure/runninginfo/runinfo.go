package runninginfo

import (
	"flag"
	"log"
	"runtime/debug"
)

// Renders runtime information
func Print(componentName string) {

	// Build time
	log.Println("Component:", componentName)

	// Go version
	info, _ := debug.ReadBuildInfo()
	log.Println("GoVersion:", info.GoVersion)

	// Parameters
	flag.VisitAll(func(a *flag.Flag) {
		log.Println("-" + a.Name + " " + a.Value.String())
	})
}
