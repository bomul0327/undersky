package gamer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Python3Driver struct {
	src string
}

func NewPython3Driver(srcModule string) *Python3Driver {
	return &Python3Driver{
		src: srcModule,
	}
}

func (driver *Python3Driver) StartProcess(port int) error {
	cmd := exec.Command("python", "gamer/python3/main.py", "--src", driver.src, "--port", strconv.Itoa(port))
	if workDir, ok := os.LookupEnv("COLOSSEUM_BASE_PATH"); ok {
		cmd.Dir = workDir
		cmd.Env = append(os.Environ(), "PYTHONPATH="+workDir)
	}

	go func() {
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb

		cmd.Run()

		fmt.Println("out: ", outb.String())
		fmt.Println("err: ", errb.String())
	}()

	return nil
}
