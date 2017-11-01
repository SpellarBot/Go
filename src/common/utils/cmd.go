// Common Bash Cmds
// @Author: Golion
// @Date: 2017.5

package utils

import (
	"os"
	"fmt"
	"os/exec"
)

func IsDirExist(path string) bool {
	p, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return p.IsDir()
	}
}

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func Remove(fileName string) {
	if IsFileExist(fileName) {
		os.Remove(fileName)
	}
}

func Exec(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("[Exec][Output] error=[%v]", err.Error())
	}
	return string(output), nil
}

func GetHostName() string {
	if hostName, err := os.Hostname(); err != nil {
		return ""
	} else {
		return hostName
	}
}

func Cp(fileName string, cpSourceDir string, cpDestinationDir string) error {
	Remove(cpDestinationDir + "/" + fileName)
	cpCmdAttr1 := cpSourceDir + "/" + fileName
	cpCmdAttr2 := cpDestinationDir
	_, err := Exec("cp", cpCmdAttr1, cpCmdAttr2)
	if err != nil {
		return fmt.Errorf("[Cp] error=[%v]", err.Error())
	}
	if IsFileExist(cpDestinationDir + "/" + fileName) {
		return nil
	} else {
		return fmt.Errorf("[Cp] Error! File Not Exist After `cp`!")
	}
}

func Mv(from string, to string) error {
	if !IsFileExist(from) {
		return fmt.Errorf("[Mv] Error! from=[%v] Not Exist!", from)
	}
	Remove(to)
	_, err := Exec("mv", from, to)
	if err != nil {
		return fmt.Errorf("[Mv] error=[%v]", err.Error())
	}
	if IsFileExist(to) {
		return nil
	} else {
		return fmt.Errorf("[Mv] Error! File Not Exist After `mv`!")
	}
}

// scpConnStr="webserver@100.84.35.57"
func Scp(fileName string, scpConnStr string, scpSourceDir string, scpDestinationDir string) error {
	Remove(scpDestinationDir + "/" + fileName)
	scpCmdAttr1 := scpConnStr + ":" + scpSourceDir + "/" + fileName
	scpCmdAttr2 := scpDestinationDir
	_, err := Exec("scp", "-P", "9922", scpCmdAttr1, scpCmdAttr2)
	if err != nil {
		return fmt.Errorf("[Scp] error=[%v]", err.Error())
	}
	if IsFileExist(scpDestinationDir + "/" + fileName) {
		return nil
	} else {
		return fmt.Errorf("[Scp] Error! File Not Exist After `scp`!")
	}
}

// 解压.gz文件
func Gunzip(fileName string) error {
	if !IsFileExist(fileName) {
		return fmt.Errorf("[Gnuzip] Error! fileName=[%v] Not Exist!", fileName)
	}
	_, err := Exec("gunzip", "-f", fileName)
	if err != nil {
		return fmt.Errorf("[Gnuzip] error=[%v]", err.Error())
	}
	return nil
}