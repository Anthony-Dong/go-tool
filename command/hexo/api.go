package hexo

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/anthony-dong/go-tool/command/log"
	"github.com/anthony-dong/go-tool/shell"
	"github.com/anthony-dong/go-tool/util"
	"github.com/juju/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	TargetFileDir = ""
	OriginFileDir = ""
	md5File       = "/Users/fanhaodong/note/note/.config/check.json"
)

type File struct {
	OriginFile    string `json:"origin_file"`
	OriginFileMd5 string `json:"origin_file_md5"`

	TargetFile    string `json:"target_file"`
	TargetFileMd5 string `json:"target_file_md5"`
}

// 1、检查文件
func checkFile(dir string) ([]string, error) {
	files, err := getOriginDir(dir)
	if err != nil {
		return nil, err
	}
	resultFile := make([]string, 0, len(files))
	for _, elem := range files {
		newName := ChangeFileName(elem)
		if newName != elem {
			log.Infof("copy elem: %s, newName: %s start", elem, newName)
			if err := shell.Mv(elem, newName); err != nil {
				log.Errorf("copy elem: %s, newName: %s error", elem, newName)
				return nil, err
			}
			log.Infof("copy elem: %s, newName: %s end", elem, newName)
		}
		resultFile = append(resultFile, newName)
	}

	needCheckFiles := make([]string, 0, len(files))
	checked := make(map[string]string, len(files))
	for _, elem := range resultFile {
		base := filepath.Base(elem)
		if base == "README.md" || base == "Untitled.md" {
			continue
		}
		fileAbs, isExist := checked[base]
		if isExist {
			return nil, errors.Errorf("find multi file name is equal, file1: %s, file2: %s", fileAbs, elem)
		} else {
			checked[base] = elem
			needCheckFiles = append(needCheckFiles, elem)
		}
	}
	return needCheckFiles, err
}

func md5Check() {

}

func findFirmCode()  {
	util.ReplaceFileContent()
}

func getOriginDir(dir string) ([]string, error) {
	return util.GetAllFiles(dir, func(fileName string) bool {
		suffix := filepath.Ext(fileName)
		base := filepath.Base(fileName)
		if (suffix == ".md" || suffix == ".markdown") && !strings.Contains(base, "README") {
			return true
		}
		return false
	})
}

func GetFileMd5(filename string) (string, error) {
	pFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	return hex.EncodeToString(md5h.Sum(nil)), nil
}

func ChangeFileName(filename string) string {
	baseName := filepath.Base(filename)
	suffix := filepath.Ext(filename)
	dir := filepath.Dir(filename)
	filename = strings.TrimSuffix(baseName, suffix)
	if strings.ToLower(filename) == "readme" {
		return filepath.Join(dir, "README.md")
	}
	length := len(filename)
	builder := strings.Builder{}
	for x := 0; x < length; x++ {
		char := filename[x]
		if char == ' ' {
			builder.WriteByte('-')
		} else {
			builder.WriteByte(char)
		}
	}
	builder.WriteString(suffix)
	return filepath.Join(dir, builder.String())
}
