#!/bin/bash

# 设置版本号和程序名
version="1.0.0"
appname="MachineIDKey"

# 设置输出目录
outdir="./build"

# 设置交叉编译参数
# 可以根据需要修改目标操作系统和架构
targets=(
    "linux/amd64"
    "windows/amd64"
    "darwin/amd64"
)

# 构建并打包可执行文件
for target in "${targets[@]}"; do
    # 解析目标操作系统和架构
    os="$(echo ${target} | cut -d'/' -f1)"
    arch="$(echo ${target} | cut -d'/' -f2)"

    # 设置输出文件名
    if [ "${os}" == "windows" ]; then
        outfile="${outdir}/${appname}-${version}-${os}-${arch}.exe"
    else
        outfile="${outdir}/${appname}-${version}-${os}-${arch}"
    fi

    # 执行交叉编译
    env GOOS=${os} GOARCH=${arch} go build -o "${outfile}" main.go

    # 打包成 tar.gz 或 zip 文件
    if [ "${os}" == "windows" ]; then
        zip -j "${outfile}.zip" "${outfile}"
    else
        tar -czf "${outfile}.tar.gz" -C "${outdir}" "${appname}-${version}-${os}-${arch}"
    fi
done
