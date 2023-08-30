package report

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	TestReportDirName = "TestReport"
	TestResultDirName = "TestResult"
	TraceLogDirName   = "Traces"
)

func StoreApiTestReport(wd string, swaggerPath string) error {
	wd = filepath.ToSlash(wd)
	swaggerPath = filepath.ToSlash(swaggerPath)
	testResultPath := path.Join(wd, TestResultDirName)
	traceLogPath := path.Join(testResultPath, TraceLogDirName)
	if _, err := os.Stat(testResultPath); !os.IsNotExist(err) {
		if err = os.RemoveAll(testResultPath); err != nil {
			return fmt.Errorf("[ERROR] error removing trace log dir %s: %+v", testResultPath, err)
		}
	}

	if err := os.MkdirAll(traceLogPath, 0755); err != nil {
		return fmt.Errorf("[ERROR] error creating trace log dir %s: %+v", testResultPath, err)
	}

	cmd := exec.Command("pal", "-i", path.Join(wd, "log.txt"), "-m", "oav", "-o", traceLogPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run  `pal` command with %+v", err)
	}

	cmd = exec.Command("oav", "validate-traffic", traceLogPath, swaggerPath, "--report", path.Join(testResultPath, "ApiTestReport.html"))
	err = cmd.Run()
	if err != nil {
		log.Printf("`oav` command completed with err: %+v", err)
	}

	return nil
}

func MergeApiTestReports(wd string, swaggerPath string) error {
	wd = filepath.ToSlash(wd)
	swaggerPath = filepath.ToSlash(swaggerPath)
	tesReportPath := path.Join(wd, TestReportDirName)
	if _, err := os.Stat(tesReportPath); !os.IsNotExist(err) {
		if err = os.RemoveAll(tesReportPath); err != nil {
			return fmt.Errorf("[ERROR] error removing test report dir %s: %+v", tesReportPath, err)
		}
	}

	traceLogPath := path.Join(tesReportPath, TraceLogDirName)
	if err := os.MkdirAll(traceLogPath, 0755); err != nil {
		return fmt.Errorf("[ERROR] error creating test report dir %s: %+v", traceLogPath, err)
	}

	destIndex := 1
	err := filepath.WalkDir(wd, func(walkPath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			traceDir := filepath.Join(walkPath, TestResultDirName, TraceLogDirName)
			if _, err = os.Stat(traceDir); !os.IsNotExist(err) {
				destIndex, err = copyTraceFiles(traceDir, traceLogPath, destIndex)
				if err != nil {
					return fmt.Errorf("failed to copy trace files: %+v", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("impossible to walk directories: %s", err)
	}

	cmd := exec.Command("oav", "validate-traffic", traceLogPath, swaggerPath, "--report", path.Join(tesReportPath, "ApiTestReport.html"))
	err = cmd.Run()
	if err != nil {
		log.Printf("`oav` command completed with err: %+v", err)
	}

	return nil
}

func copyTraceFiles(src string, dest string, destIndex int) (int, error) {
	err := filepath.WalkDir(src, func(walkPath string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		srcFile, err := os.Open(walkPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destPath := path.Join(dest, fmt.Sprintf("trace-%d.json", destIndex))
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		destIndex++
		return nil
	})

	if err != nil {
		return destIndex, err
	}

	return destIndex, nil
}
