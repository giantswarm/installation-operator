package ansible

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Env    map[string]string
	Logger micrologger.Logger
	Out    io.Writer
}

type Executor struct {
	env    map[string]string
	logger micrologger.Logger
	out    io.Writer
}

func New(config Config) Executor {
	return Executor{
		env:    config.Env,
		logger: config.Logger,
		out:    config.Out,
	}
}

func (e Executor) Execute(command string, args []string, _ string) error {
	stderr := &bytes.Buffer{}

	cmd := exec.Command(command, args...)
	cmd.Stderr = stderr

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return microerror.Mask(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Fprintf(e.out, "%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return microerror.Mask(err)
	}

	err = cmd.Wait()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
