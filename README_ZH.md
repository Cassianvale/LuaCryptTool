# LuaCryptTool
简体中文 / [English](README.md)

## 项目介绍

本工具主要针对Lua文件进行批量 **转字节码** 或 **反编译字节码**
本项目使用go build进行打包，并使用 upx 对打包文件进行了高质量压缩

- Encryptor: 批量加密
`go build -ldflags="-s -w" -o encryptor.exe`
`upx -9 encryptor.exe`

- EncryptorDecryptor: 批量加密 & 解密
`upx -9 encryptor_decryptor.exe`
`go build -ldflags="-s -w" -o encryptor_decryptor.exe`

## 使用方式
打开dist目录，安装`./dist/jre-8u251-windows-x64.exe`并配置环境变量，运行已经打包好的`./dist/encryptor.exe` 或 `./dist/encryptor_decryptor.exe` 文件
