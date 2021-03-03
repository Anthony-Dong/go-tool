package hexo

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthony-dong/go-tool/command/log"
	"github.com/anthony-dong/go-tool/shell"
	"github.com/anthony-dong/go-tool/util"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

var (
	delimiter = "---"
)

type Config struct {
	Title      string   `yaml:"title"`       // 标题(如果没有设置，为源文件的名称)
	TargetFile string   `yaml:"target_file"` // 目标文件，值得是生成的文件
	OriginFile string   `yaml:"origin_file"` // 原文件，指的是我们写的文件
	Date       string   `yaml:"date"`        // 日期(为文件的修改日期)
	Tags       []string `yaml:"tags"`
	Categories []string `yaml:"categories"`
}

type CheckFileCanHexoResult struct {
	CanHexo       bool
	HasAbstract   bool
	Content       []string
	Config        *Config
	FileInfo      os.FileInfo
	FileName      string // 源文件
	WriteFilePath string // 写入的目录
	NeedWrite     bool   // 是否需要重新写入
	ContentData   []byte
}

func (c *CheckFileCanHexoResult) WriteFile() error {
	checkFileNeedReWrite := func(data []byte) (bool, error) {
		// 当前文件
		oldData, err := util.ReadFile(c.FileName)
		if err != nil {
			return false, err
		}
		return !(util.Md5(oldData) == util.Md5(data)), nil
	}

	// 组装文件
	buffer := &bytes.Buffer{}
	buffer.Write([]byte(delimiter))
	buffer.WriteByte('\n')
	config, err := yaml.Marshal(c.Config)
	if err != nil {
		return err
	}
	buffer.Write(config)
	buffer.Write([]byte(delimiter))
	buffer.WriteByte('\n')
	buffer.WriteByte('\n')
	content := util.SliceLineToString(c.Content)
	buffer.Write(util.String2Slice(content))

	data := buffer.Bytes()
	// 检测文件是否需要重写，通过MD5校验
	needWrite, err := checkFileNeedReWrite(data)
	if err != nil {
		return err
	}
	if needWrite {
		log.Infof("[Hexo] 发现MD5比较不一致需要重新写入到源文件, 文件: %s", c.FileName)
		if err := util.WriteFileBody(c.FileName, data); err != nil {
			return err
		}
	}
	c.NeedWrite = needWrite
	c.ContentData = data
	return nil
}

func Run(ctx context.Context, dir string, targetDir string, firmCode []string) error {
	dir, err := util.Abs(filepath.Clean(dir))
	if err != nil {
		return errors.Trace(err)
	}
	targetDir, err = util.Abs(filepath.Clean(targetDir))
	if err != nil {
		return errors.Trace(err)
	}
	log.Debugf("[Hexo] 开始全部的Markdown的文件, 目录: %s", dir)
	allPage, err := GetAllPage(dir)
	if err != nil {
		return errors.Trace(err)
	}
	log.Debugf("[Hexo] 获取全部的Markdown文件成功, 目录: %s, 总数: %d", dir, len(allPage))

	log.Debugf("[Hexo] 开始全部的Target-Markdown的文件, 目录: %s", targetDir)
	targetPage, err := GetAllMarkDownPage(targetDir)
	if err != nil {
		return errors.Trace(err)
	}
	log.Debugf("[Hexo] 获取全部的Target-Markdown文件成功, 目录: %s, 总数: %d", targetDir, len(targetPage))
	targetPageSet := util.NewSet(targetPage)
	newTargetPage := util.NewSetInitSize(targetPageSet.Size())

	wg := sync.WaitGroup{}
	for _, file := range allPage {
		wg.Add(1)
		go func(fileName string) {
			needEnd := false
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("[Hexo] 运行期间发现 panic, 文件: %s, 异常: %v", fileName, err)
				}
				wg.Done()
				if needEnd {
					log.Debugf("[Hexo] 结束操作文件, 文件: %s", fileName)
				}
			}()
			// 检测是否是hexo
			result, err := CheckFileCanHexo(fileName, dir)
			if err != nil {
				log.Errorf("[Hexo] 检测文件是否是hexo文件发现异常, 文件: %s, 异常: %v", fileName, err)
				return
			}
			if result == nil {
				return
			}
			log.Debugf("[Hexo] 开始操作文件, 文件: %s", fileName)
			if !result.HasAbstract {
				log.Warnf("[Hexo] 警告, 发现没有摘要, 文件: %s", fileName)
			}
			needEnd = true
			// 开始检测是否有公司代码
			if err := CheckFileHasFirmCode(fileName, result.Content, firmCode); err != nil {
				log.Errorf("[Hexo] 检测公司代码失败, 异常: %s, 文件: %s", err, fileName)
				return
			}

			// 检测文件格式写入
			if err := result.WriteFile(); err != nil {
				log.Errorf("[Hexo] 检测原文件格式失败, 异常: %s, 文件: %s", err, fileName)
				return
			}

			// copy文件

			targetFile := filepath.Join(targetDir, result.Config.TargetFile)

			writeTargetFile := func() {
				if err := util.WriteFileBody(targetFile, result.ContentData); err != nil {
					log.Errorf("[Hexo] 发现需要写入到Hexo的post目录文件发现了异常, 文件: %s, 异常: %v", targetFile, err)
					return
				}
			}

			// 不存在
			if !targetPageSet.Contains(targetFile) {
				log.Infof("[Hexo] 发现Hexo的post目录文件不存在需要写入的文件, 目标文件: %s, 源文件: %s", targetFile, fileName)
				writeTargetFile()
			} else {
				// 比较MD5是否相同
				readBody, err := util.ReadFile(targetFile)
				if err != nil {
					log.Errorf("[Hexo] 发现读取Hexo的post目录文件发现了异常, 文件: %s, 异常: %v", targetFile, err)
					return
				}
				if util.Md5(result.ContentData) != util.Md5(readBody) {
					log.Infof("[Hexo] 发现读取Hexo的post目录文件和原文件MD5值不一样, 需要重写, 文件: %s,  源文件: %s", targetFile, fileName)
					writeTargetFile()
				}
			}

			// 操作完成
			newTargetPage.Put(targetFile)
			return
		}(file)
	}
	wg.Wait()

	// delete
	log.Infof("[Hexo] 操作脚本完成, 一共写入: %d", newTargetPage.Size())

	slice := targetPageSet.ToSlice()
	for _, elem := range slice {
		if newTargetPage.Contains(elem) {
			targetPageSet.Delete(elem)
		}
	}

	log.Infof("[Hexo] 操作脚本完成需要删除文件, 删除: %d", targetPageSet.Size())

	if targetPageSet.Size() == 0 {
		return nil
	}

	for _, elem := range targetPageSet.ToSlice() {
		log.Infof("[Hexo] 删除文件, 文件: %s", elem)
		if err := shell.Delete(elem); err != nil {
			log.Errorf("[Hexo] 删除文件失败, 文件: %s, 异常: %s", elem, err)
			return err
		}
	}

	return nil
}

// 是否有
func CheckFileHasFirmCode(fileName string, content []string, firmCode []string) error {
	if firmCode == nil || len(firmCode) == 0 || content == nil || len(content) == 0 {
		return nil
	}
	for index, line := range content {
		for _, elem := range firmCode {
			if strings.Contains(line, elem) {
				log.Warnf("[Hexo] 发现公司代码, 文件: %s, 公司代码: %s, 原文: %s", fileName, elem, line)
				newElem := util.NewString(len(elem), 'x')
				line = strings.ReplaceAll(line, elem, newElem)
			}
		}
		content[index] = line
	}
	return nil
}

// true 表示是
func CheckFileCanHexoPre(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer file.Close()
	trueError := errors.New("true")
	falseError := errors.New("false")
	count := 0
	if err := ReadFile(file, func(line string) error {
		if line == "" {
			return nil
		}
		count++
		if count == 1 && line == delimiter {
			return trueError
		} else {
			return falseError
		}
	}); err != nil {
		if err == trueError {
			return true
		}
		return false
	}
	return false
}

func CheckFileCanHexo(fileName string, filePath string) (*CheckFileCanHexoResult, error) {
	if !CheckFileCanHexoPre(fileName) {
		return nil, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 摘要分隔符
	abstractDelimiter := "<!-- more -->"
	hasAbstract := false

	// 头部信息分隔符
	delimiterCount := 0
	delimiterErr := errors.New("delimiter")

	// yaml-config
	yamlConfig := make([]string, 0)
	// body
	content := make([]string, 0)

	// 是否有空格 1、必须有分隔符，2、必须是 isspace=true 如果出现了非空格则设置为false
	isSpace := true
	if err := ReadFile(file, func(line string) error {
		// 如果分隔符
		if line == delimiter {
			delimiterCount++
			return nil
		}

		// 不是
		if !hasAbstract && strings.Contains(line, abstractDelimiter) {
			hasAbstract = true
		}

		// 如果是刚开始则为 yaml
		if delimiterCount == 1 {
			yamlConfig = append(yamlConfig, line)
			return nil
		}

		// 为正文
		if delimiterCount >= 2 {
			if line == "" && isSpace {
				return nil
			}
			isSpace = false
			content = append(content, line)
			return nil
		}
		return nil
	}); err != nil {
		if err != delimiterErr {
			return nil, err
		}
	}
	canHexo := delimiterCount == 2
	if !canHexo {
		return nil, nil
	}

	config := util.SliceLineToString(yamlConfig)
	fileConfig := new(Config)
	err = yaml.Unmarshal(util.String2Slice(config), fileConfig)
	if err != nil {
		return nil, err
	}

	// title为空
	if fileConfig.Title == "" {
		fileConfig.Title = strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name()))
	}

	// 源文件
	fileConfig.OriginFile = strings.TrimPrefix(strings.TrimPrefix(fileName, filePath), string(filepath.Separator))
	// 目标文件
	if fileConfig.TargetFile == "" {
		fileConfig.TargetFile = strings.ReplaceAll(fileConfig.Title, " ", "-") + ".md"
	}

	// 文件修改时间
	if fileConfig.Date == "" {
		fileConfig.Date = fileInfo.ModTime().Format(util.FromatTime_V1)
	}

	return &CheckFileCanHexoResult{
		CanHexo:     canHexo,
		HasAbstract: hasAbstract,
		Content:     content,
		FileInfo:    fileInfo,
		Config:      fileConfig,
		FileName:    fileName,
	}, nil
}

func ReadFile(file io.Reader, foo func(line string) error) error {
	reader := bufio.NewReader(file)
	for {
		lines, isEOF, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Trace(err)
		}
		if isEOF {
			break
		}
		if err := foo(string(lines)); err != nil {
			return err
		}
	}
	return nil
}

func GetAllPage(dir string) ([]string, error) {
	return util.GetAllFiles(dir, func(fileName string) bool {
		suffix := filepath.Ext(fileName)
		base := filepath.Base(fileName)
		if !(suffix == ".md" || suffix == ".markdown") {
			return false
		}
		dir := filepath.Dir(fileName)
		if strings.Contains(dir, "bin") || strings.Contains(dir, "hexo-home") || strings.Contains(dir, ".config") {
			return false
		}
		if strings.Contains(base, "README") {
			return false
		}
		return true
	})
}

func GetAllMarkDownPage(dir string) ([]string, error) {
	return util.GetAllFiles(dir, func(fileName string) bool {
		suffix := filepath.Ext(fileName)
		if !(suffix == ".md" || suffix == ".markdown") {
			return false
		}
		return true
	})
}
