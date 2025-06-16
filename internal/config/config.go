package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() (map[string]string, error) {
	rootDir, _ := os.Getwd()
	path := filepath.Join(rootDir, ".env")
	if _, err := os.Stat(path); err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Неудалось найти .env файл")
	}
	_ = godotenv.Load(path)
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Не удалось прочитать .env файл")
	}

	return envMap, nil
}

func LoadJSON() (map[string]string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить рабочую директорию: %w", err)
	}

	path := filepath.Join(rootDir, "config.json")

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config.json не найден: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var cfg map[string]string
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return cfg, nil
}