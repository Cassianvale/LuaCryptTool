package main

// go build -ldflags="-s -w" -o encryptor_decryptor.exe

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed *
var embeddedFiles embed.FS

// 定义函数类型，该函数接受两个string参数并返回一个int和一个error
type actionFuncType func(string, string) (int, error)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请选择模式( 1: 加密，2: 解密 ): ")
		if !scanner.Scan() {
			return
		}
		mode := scanner.Text()

		folder, outfolder := getFolders(scanner)

		var actionFunc actionFuncType
		switch mode {
		case "1":
			actionFunc = encryptFile
		case "2":
			actionFunc = decryptFile
		default:
			fmt.Println("无效的模式选择。")
			continue
		}

		if filenums, err := processFiles(folder, outfolder, actionFunc); err != nil {
			fmt.Printf("操作失败: %v\n", err)
		} else {
			fmt.Printf("操作完成, 共计 %d 个文件被处理。\n", filenums)
		}

		if !shouldContinue(scanner) {
			break
		}
	}
}

func processFiles(src, out string, actionFunc actionFuncType) (int, error) {
	if err := clearFolder(out); err != nil {
		return 0, err
	}
	return actionFunc(src, out)
}

func clearFolder(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

func encryptFile(src, out string) (int, error) {
	luacData, err := embeddedFiles.ReadFile("luac.exe")
	if err != nil {
		return 0, err
	}
	luacPath := filepath.Join(os.TempDir(), "luac.exe")
	if err := os.WriteFile(luacPath, luacData, 0755); err != nil {
		return 0, err
	}
	defer os.Remove(luacPath)

	return processDirectory(src, out, func(srcFile, outFile string) error {
		cmd := exec.Command(luacPath, "-o", outFile, srcFile)
		return cmd.Run()
	})
}

func decryptFile(src, out string) (int, error) {
	unluacData, err := embeddedFiles.ReadFile("unluac.jar")
	if err != nil {
		return 0, err
	}
	unluacPath := filepath.Join(os.TempDir(), "unluac.jar")
	if err := os.WriteFile(unluacPath, unluacData, 0644); err != nil {
		return 0, err
	}
	defer os.Remove(unluacPath)

	return processDirectory(src, out, func(srcFile, outFile string) error {
		cmdStr := fmt.Sprintf("java -jar %s --rawstring %s > %s", unluacPath, srcFile, outFile)
		cmd := exec.Command("cmd", "/C", cmdStr)
		return cmd.Run()
	})
}

func processDirectory(src, out string, processFile func(string, string) error) (int, error) {
	files, err := os.ReadDir(src)
	if err != nil {
		return 0, err
	}

	var filenums int
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		outPath := filepath.Join(out, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))+".lua")
		if file.IsDir() {
			subFileNums, err := processDirectory(srcPath, outPath, processFile)
			if err != nil {
				return filenums, err
			}
			filenums += subFileNums
		} else {
			if err := processFile(srcPath, outPath); err == nil {
				filenums++
			} else {
				return filenums, err
			}
		}
	}
	return filenums, nil
}

func getFolders(scanner *bufio.Scanner) (string, string) {
	fmt.Print("请拖拽源文件夹: ")
	scanner.Scan()
	folder := scanner.Text()

	fmt.Print("请拖拽输出文件夹: ")
	scanner.Scan()
	outfolder := scanner.Text()

	return folder, outfolder
}

func shouldContinue(scanner *bufio.Scanner) bool {
	fmt.Print("输入数字键3继续进行操作，输入其他键退出: ")
	scanner.Scan()
	continueCode := scanner.Text()
	return continueCode == "3"
}
