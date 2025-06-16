@echo off
setlocal EnableDelayedExpansion

echo === Go cross-compilation script for Windows ===

set VALID_GOOS=linux windows darwin
set VALID_GOARCH=amd64 arm64 386

if "%~1"=="" (
    set /p GOOS=Выберите систему для сборки (linux/windows/darwin): 
) else (
    set GOOS=%~1
)

if "%~2"=="" (
    set /p GOARCH=Выберите архитектуру для сборки (amd64/arm64/386): 
) else (
    set GOARCH=%~2
)

set FOUND_GOOS=false
for %%G in (%VALID_GOOS%) do (
    if /I "%%G"=="%GOOS%" set FOUND_GOOS=true
)
if "!FOUND_GOOS!"=="false" (
    echo Ошибка: недопустимое значение GOOS: %GOOS%
    exit /b 1
)

set FOUND_GOARCH=false
for %%A in (%VALID_GOARCH%) do (
    if /I "%%A"=="%GOARCH%" set FOUND_GOARCH=true
)
if "!FOUND_GOARCH!"=="false" (
    echo Ошибка: недопустимое значение GOARCH: %GOARCH%
    exit /b 1
)

set /p OUTPUT_NAME=Введите имя выходного файла (по умолчанию PrintersControll): 
if "%OUTPUT_NAME%"=="" set OUTPUT_NAME=PrintersControll

if /I "%GOOS%"=="windows" (
    set OUTPUT_NAME=%OUTPUT_NAME%.exe
)

echo.
echo ▶ Сборка: GOOS=%GOOS%, GOARCH=%GOARCH% → %OUTPUT_NAME%
set GOOS=%GOOS%
set GOARCH=%GOARCH%
go build -o %OUTPUT_NAME%

if %errorlevel%==0 (
    echo ✅ Сборка успешно завершена: %OUTPUT_NAME%
) else (
    echo ❌ Ошибка при сборке
    exit /b 1
)

endlocal
