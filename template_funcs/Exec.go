package template_funcs

import (
	"bytes"
	"os/exec"
	"strings"
)

func Exec(cmdStr string) (string, error) {
	parts := strings.Split(cmdStr, " ")
	cmd := exec.Command(parts[0], parts[0:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
