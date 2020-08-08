package shell

import (
	"fmt"
	"io"
	"os"
)

type shell struct {
	prefix string
	w      io.Writer
}

func NewShell(prefix string) *shell {
	return &shell{
		prefix: prefix,
		w:      os.Stdout,
	}
}

func (this *shell) PrintLine() {
	str := fmt.Sprintf("\033[32m%s\033[0m", this.prefix)
	fmt.Fprint(this.w, str)
}
func (this *shell) PrintInfo(str string) {
	fmt.Fprintln(this.w, str)
}

