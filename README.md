# LuaCryptTool
[简体中文](README_ZH.md) / English

## Project Introduction

This tool is primarily designed for batch **bytecode conversion** or **decompilation of bytecode** for Lua files.
The project utilizes `go build` for packaging, and employs `upx` for high-quality compression of the packaged files.

- Encryptor: For batch encryption
`go build -ldflags="-s -w" -o encryptor.exe`
`upx -9 encryptor.exe`

- EncryptorDecryptor: For batch encryption & decryption
`upx -9 encryptor_decryptor.exe`
`go build -ldflags="-s -w" -o encryptor_decryptor.exe`

## How to Use
Navigate to the `dist` directory, install `./dist/jre-8u251-windows-x64.exe` and set up the environment variables, then run the packaged `./dist/encryptor.exe` or `./dist/encryptor_decryptor.exe` files.