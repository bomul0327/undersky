package gamer

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
)

type Python3Driver struct{}

func (driver *Python3Driver) StartProcess(port int) error {
	var outb, errb bytes.Buffer

	cmd := exec.Command("python", "gamer/python3/main.py", "--src", "baseline.1000", "--port", strconv.Itoa(port))
	if workDir, ok := os.LookupEnv("COLOSSEUM_BASE_PATH"); ok {
		cmd.Dir = workDir
		cmd.Env = append(os.Environ(), "PYTHONPATH="+workDir)
	}
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	go func() { cmd.Run() }()

	return nil
}
