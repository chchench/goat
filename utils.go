package gotd

import (
	"os/exec"
	"strings"
)

func LaunchCommand(targetDir string, prog string, cmdLineArgs []string) (string, error) {

	cmd := exec.Command(prog, cmdLineArgs...)

	if targetDir != "" {
		cmd.Dir = targetDir
	}

	out, err := cmd.Output()
	return string(out), err
}

func parseCmdLine4ProgNParams(str string) (string, []string) {
	fields := splitNClean(str)
	prog := fields[0]
	return prog, fields[1:]
}

func splitNClean(str string) []string {
	var newList []string

	list := strings.Split(str, " ")

	for _, s := range list {
		if s == "" {
			continue
		}
		newList = append(newList, s)
	}
	return newList
}
