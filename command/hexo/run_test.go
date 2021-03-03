package hexo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAllPage(t *testing.T) {
	list, err := GetAllPage("/Users/fanhaodong/note/note")
	if err != nil {
		t.Fatal(err)
	}
	for _, elem := range list {
		fmt.Println(elem)
	}
}

func TestReadFile(t *testing.T) {
	file, err := os.Open("/Users/fanhaodong/note/note/设计力架构力/测试/自动化测试.md")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := ReadFile(file, func(line string) error {
		fmt.Println(line)
		return nil
	}); err != nil {
		return
	}
}

func TestCheckFileCanHexo(t *testing.T) {
	result, err := CheckFileCanHexo("/Users/fanhaodong/go/code/go-tool/bin/test.md","/Users/fanhaodong/go/code/go-tool")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", result.Config)
	fmt.Printf("%#v", result)
}

func TestCheckFileHasFirmCode(t *testing.T) {
	firmCode := []string{"baidu", "ali"}
	result, err := CheckFileCanHexo("/Users/fanhaodong/go/code/go-tool/bin/test.md","/Users/fanhaodong/go/code/go-tool")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckFileHasFirmCode("/Users/fanhaodong/go/code/go-tool/bin/test.md", result.Content, firmCode); err != nil {
		t.Fatal(err)
	}

	fmt.Println(result.Content)

	if err := result.WriteFile(); err != nil {
		t.Fatal(err)
	}

}

func TestRun(t *testing.T) {
	dir := filepath.Clean("/Users/fanhaodong/note/note")
	targetDir := filepath.Clean("/Users/fanhaodong/note/note/hexo-home/source/_posts")
	firmCode := []string{"baidu", "ali"}
	if err := Run(context.Background(), dir, targetDir, firmCode); err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	println(CheckFileCanHexoPre("/Users/fanhaodong/go/code/go-tool/test/test.md"))
}
