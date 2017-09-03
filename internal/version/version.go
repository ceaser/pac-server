package version

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	showVersion = flag.Bool("version", false, "show app version")

	// specifiy version, BuildTimeUTC, AppName at build time with `-ldflags "-X path.to.package.Version x.x.x"` etc...
	Version      = "-"
	BuildTimeUTC = "-"
	AppName      = "-"
)

func ShowVersion() {
	if *showVersion {
		fmt.Println(String())
		os.Exit(0)
	}
}

func String() string {
	return fmt.Sprintf(
		"%s %s (built w/%s)\nUTC Build Time: %v",
		AppName,
		Version,
		runtime.Version(),
		BuildTimeUTC,
	)
}
