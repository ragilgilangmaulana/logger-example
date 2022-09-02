package serror

import (
	"fmt"
	"logging-example-go/cservice"
	"logging-example-go/utils/utstring"
	"os"
	"strings"
	"syscall"
)

func isLocal() bool {
	return strings.ToLower(utstring.Env(cservice.AppEnv, cservice.EnvLocal)) == cservice.EnvLocal
}

func printErr(m string) {
	fmt.Fprintln(os.Stderr, m)
}

func exit() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		os.Exit(1)
	}
}

func getPath(val string) string {
	for _, v := range rootPaths {
		if strings.HasPrefix(val, v) {
			val = utstring.Sub(val, len(v), 0)
			return val
		}
	}

	return val
}
