package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"printers/internal/config"
	"strconv"
	"time"

	"printers/internal/interfaces"
)

func GetPrintersInfo(ip string) (interfaces.PrinterInfo, error) {
	cfg, err := config.LoadJSON()
	if err != nil {
		log.Fatal(err)
	}

	queryURL := fmt.Sprintf("http://%s:%s/printer/objects/query?print_stats", ip, cfg["PORT_MOON"])
	resp, err := http.Get(queryURL)

	if err != nil {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка запроса print_stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка в нахождении файла: %s", resp.Status)
	}

	var queryResp interfaces.QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка разбора JSON print_stats: %w", err)
	}

	filename := queryResp.Result.Status.PrintStats.Filename
	if filename == "" {
		return interfaces.PrinterInfo{}, fmt.Errorf("файл не найден")
	}

	escapedFilename := url.QueryEscape(filename)
	metdataURL := fmt.Sprintf("http://%s:%s/server/files/metadata?filename=%s", ip, cfg["PORT_MOON"], escapedFilename)
	metaResp, err := http.Get(metdataURL)

	if err != nil {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка запроса metadata: %w", err)
	}
	defer metaResp.Body.Close()

	if metaResp.StatusCode != http.StatusOK {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка в парсинге файла: %s", metaResp.Status)
	}

	var metadata interfaces.MetadataResponse
	if err := json.NewDecoder(metaResp.Body).Decode(&metadata); err != nil {
		return interfaces.PrinterInfo{}, fmt.Errorf("ошибка разбора JSON metadata: %w", err)
	}

	estimatedTime := metadata.Result.EstimatedTime
	printStartTime := metadata.Result.PrintStartTime
	currentTime := float64(time.Now().Unix())

	if estimatedTime == 0 {
		return interfaces.PrinterInfo{}, interfaces.PrinterInfo{}
	}

	progress := ((currentTime - printStartTime) / estimatedTime) * 100
	timeRemaining := printStartTime + estimatedTime - currentTime
	utcOffsetStr := cfg["UTC_FORMAT"]
	utcOffset, _ := strconv.Atoi(utcOffsetStr)

	dateEnd := time.Unix(int64(printStartTime+estimatedTime+(3600*float64(utcOffset))), 0).UTC()

	return interfaces.PrinterInfo{
		Success:       fmt.Sprintf("%.2f", progress),
		DateEnd:       dateEnd.Format("2006-01-02 15:04:05"),
		EstimatedTime: (time.Duration(timeRemaining) * time.Second).String(),
	}, nil
}

func GetPhoto(hostIP string) (string, error) {
	imageURL := fmt.Sprintf("http://%s/timelapse/timelapse.jpg?cacheBust=%d", hostIP, time.Now().Unix())

	resp, err := http.Get(imageURL)
	if err != nil {
		log.Println("Ошибка скачивания изображения: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.Header.Get("Content-Type") != "image/jpeg" {
		log.Printf("Неверный ответ при загрузке картинки: статус %d, тип %s\n", resp.StatusCode, resp.Header.Get("Content-Type"))
		return "", err
	}

	photoPath := filepath.Join(os.TempDir(), "photo.jpg")

	file, err := os.Create(photoPath)
	if err != nil {{
		log.Println("Ошибка создания файла: ", err)
		return "", err
	}}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Ошибка записи файла: ", err)
		return "", err
	}

	return photoPath, err
}