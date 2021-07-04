package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthony-dong/go-tool/commons/logger"
	"github.com/juju/errors"
)

var (
	shell       string
	shellPrefix string
)

/**
获取当前用户执行的shell.
*/
func init() {
	shell = os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	shellPrefix = shell + " -c "
}

func GitClone(sshAddr, dir string) error {
	return errors.Trace(GitCloneBranch("master", sshAddr, dir))
}

func GitCloneBranch(branch, sshAddr, dir string) error {
	gitCmd := fmt.Sprintf("git clone -b %s %s %s", branch, sshAddr, dir)
	return errors.Trace(Cmd(gitCmd))
}

func Run(shell string) error {
	gitCmd := fmt.Sprintf("%s", shell)
	return Cmd(gitCmd)
}

func Copy(src, dest string) (err error) {
	gitCmd := fmt.Sprintf("cp -R %s %s", src, dest)
	return Cmd(gitCmd)
}

func Mv(src, dest string) (err error) {
	gitCmd := fmt.Sprintf("mv '%s' '%s'", src, dest)
	return Cmd(gitCmd)
}

// delete file.
func Delete(file ...string) (err error) {
	if file == nil || len(file) == 0 {
		return
	}
	for _, elem := range file {
		if elem == "/" || strings.Contains(elem, "*") {
			return errors.New("can not delete * file")
		}
	}
	gitCmd := fmt.Sprintf("rm -r '%s'", strings.Join(file, " "))
	return errors.Trace(Cmd(gitCmd))
}

func Cmd(cmd string) error {
	command := exec.Command(shell, "-c", cmd)
	_shellLogger.Infof("exec: %s", strings.TrimPrefix(strings.TrimPrefix(command.String(), " "), shellPrefix))
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return errors.Trace(command.Run())
}

var (
	_shellLogger = logger.NewStdOutLogger(logger.NameOp("[Shell]"), logger.FlagOp(logger.TimeLoggerFormat))
)
