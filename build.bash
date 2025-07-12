#!/bin/bash

echo "==== Go cross-compilation script ===="

VALID_GOOS=("linux" "windows", "darwin")
VALID_GOARCH=("amd64" "arm64" "386")

if [ -z "$1" ]; then
    echo "Выберите систему для сборки (linux/windows/darwin): "
    read GOOS
else
    GOOS=$1
fi

if [ -z "$2" ]; then
    echo "Выберите архитектуру для сборки (amd64/arm64/386): "
    read GOARCH
else
    GOARCH=$2
fi

if [[ ! " ${VALID_GOOS[*]} " =~ " ${GOOS} "]]; then
    echo "❌ Недопустимое значение GOOS: $GOOS"
    exit 1
fi

if [[ !" ${VALID_GOARCH[*]} " =~ " ${GOARCH} "]]; then
    echo "❌ Недопустимое значение GOARCH: $GOARCH"
    exit 1
fi

echo "Введите имя итогового файла (по умолчанию PrintersControll): "
read OUTPUT_NAME
OUTPUT_NAME=${OUTPUT_NAME:-PrintersControll}

if [ "$GOOS" == "windows" ]; then
    OUTPUT_NAME="$OUTPUT_NAME.exe"
fi

echo "▶ Выполняется сборка: GOOS=$GOOS, GOARCH=$GOARCH -> $OUTPUT_NAME"
GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_NAME"

if [ $? -eq 0 ]; then
    echo "✅ Сборка успешно завершена: $OUTPUT_NAME"
else
    echo "❌ Ошибка при сборке"
    exit 1
fi