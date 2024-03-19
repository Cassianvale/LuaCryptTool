package main

// go build -ldflags="-s -w" -o encryptor.exe
import (
	"bufio"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed luac.exe
var embeddedFiles embed.FS

type actionFuncType func(string, string) (int, error)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请拖拽源文件夹: ")
		if !scanner.Scan() {
			return
		}
		folder := scanner.Text()

		fmt.Print("请拖拽输出文件夹: ")
		if !scanner.Scan() {
			return
		}
		outfolder := scanner.Text()

		if filenums, err := processFiles(folder, outfolder, encryptFile); err != nil {
			fmt.Printf("加密失败: %v\n", err)
		} else {
			fmt.Printf("加密完成, 共计 %d 个文件被处理。\n", filenums)
		}

		if !shouldContinue(scanner) {
			break
		}
	}
}

func processFiles(src, out string, actionFunc actionFuncType) (int, error) {
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		return 0, err
	}
	return actionFunc(src, out)
}

func encryptFile(src, out string) (int, error) {
	luacData, err := embeddedFiles.ReadFile("luac.exe")
	if err != nil {
		return 0, err
	}
	luacPath := filepath.Join(os.TempDir(), "luac.exe")
	if err := ioutil.WriteFile(luacPath, luacData, 0755); err != nil {
		return 0, err
	}
	defer os.Remove(luacPath)

	return processDirectory(src, out, func(srcFile, outFile string) error {
		cmd := exec.Command(luacPath, "-o", outFile, srcFile)
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
		outPath := filepath.Join(out, file.Name()+".bytes")
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

func shouldContinue(scanner *bufio.Scanner) bool {
	fmt.Print("输入数字键1继续加密其他文件，输入其他键退出: ")
	scanner.Scan()
	continueCode := scanner.Text()
	return continueCode == "1"
}
