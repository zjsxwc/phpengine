package phpengine

import (
	"fmt"
	"github.com/syyongx/php2go"
	"os"
	"os/exec"
	"strconv"
)

type PhpEngine struct {
	phpCliBinPath string
}

func NewPhpEngine(phpCliBinPath string) *PhpEngine {
	return &PhpEngine{phpCliBinPath: phpCliBinPath}
}

func (pe *PhpEngine) NewSession(scriptPath string) *Session {
	tmpScriptPath := "/tmp/phpengine_" + strconv.Itoa(int(php2go.Time())) + "_" + php2go.Md5(scriptPath) + "_" + strconv.Itoa(php2go.Rand(1, 99999)) + ".php"

	//todo
	//  prepare php $_GET $_POST $_SERVER $_REQUEST
	//  file_get_contents("php://input"); 读取http请求里的body也就是没有header等数据的内容

	err := php2go.FilePutContents(tmpScriptPath,
		fmt.Sprintf("<?php\nglobal $phpengineData;\n\ninclude_once \"%s\";file_put_contents(\"%s.out.json\", json_encode($phpengineData));", scriptPath, tmpScriptPath),
		os.FileMode(777))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &Session{phpEngine: pe, tmpScriptPath: tmpScriptPath, scriptPath: scriptPath}
}

type Session struct {
	phpEngine     *PhpEngine
	tmpScriptPath string
	scriptPath    string
	lastOutput    string
}

func (s *Session) GetLastOutput() string {
	return s.lastOutput
}
func (s *Session) Clear() {
	_ = php2go.Unlink(s.tmpScriptPath)
}

func (s *Session) Execute() string {
	f, err := exec.Command(s.phpEngine.phpCliBinPath, s.tmpScriptPath).Output()
	s.lastOutput = string(f)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	phpengineData, err := php2go.FileGetContents(s.tmpScriptPath + ".out.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = php2go.Unlink(s.tmpScriptPath + ".out.json")
	return phpengineData
}
