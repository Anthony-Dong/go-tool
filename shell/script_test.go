package shell

import (
	"fmt"
	"github.com/juju/errors"
	"testing"
)

func TestCmd(t *testing.T) {
	err := Cmd("ls -al")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	fmt.Println(Delete("/data/test/go-template", "/data/test/go-script"))
}

func TestGit(t *testing.T) {
	err := GitClone("git@gitlab.corp.xxxx7.com:ebike-urban-op/go-script.git", "/Users/fanhaodong/data/test/demo")
	if err != nil {
		fmt.Println(errors.ErrorStack(err))
	}
}

func TestGitBranch(t *testing.T) {

}

func TestRun(t *testing.T) {

}

func TestCopy(t *testing.T) {

}

func TestMv(t *testing.T) {
	fmt.Println(Mv("/data/test/go2sky", "/data/test/go2sky-1"))
}
