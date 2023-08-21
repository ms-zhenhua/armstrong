package report

import (
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	LogFileName     = "log.txt"
	TraceLogDirName = "trace_logs"
)

func StoreApiTestReport(wd string, swaggerPath string) {
	storeTraceLogs(wd, swaggerPath)
}

func storeTraceLogs(wd string, swaggerPath string) {
	if _, err := os.Stat(path.Join(wd, TraceLogDirName)); !os.IsNotExist(err) {
		if err = os.RemoveAll(path.Join(wd, TraceLogDirName)); err != nil {
			log.Fatalf("[ERROR] error removing trace log dir %s: %+v", TraceLogDirName, err)
		}
	}

	if err := os.Mkdir(path.Join(wd, TraceLogDirName), 0755); err != nil {
		log.Fatalf("[ERROR] error creating trace log dir %s: %+v", TraceLogDirName, err)
	}

	cmd := exec.Command("pal", "-i", path.Join(wd, "log.txt"), "-m", "oav", "-o", path.Join(wd, TraceLogDirName))
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

}
