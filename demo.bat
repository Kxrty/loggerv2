@echo off
chcp 65001 >nul
echo ================================================
echo    ДЕМОНСТРАЦИЯ LoggerV2
echo    Обработчик логов с приведением к ГОСТ РФ
echo ================================================
echo.

echo [1/5] Обработка Syslog (RFC 3164)
echo --------------------------------------------
echo ^<134^>Oct 11 22:14:15 server su: login failed | logger.exe
echo.

timeout /t 2 >nul

echo [2/5] Обработка CEF (Common Event Format)
echo --------------------------------------------
echo CEF:0^|Security^|IDS^|1.0^|100^|Attack^|10^|src=10.0.0.1 dst=192.168.1.1 | logger.exe
echo.

timeout /t 2 >nul

echo [3/5] Обработка LEEF (Log Event Extended Format)
echo --------------------------------------------
echo LEEF:1.0^|IBM^|QRadar^|7.3^|Login^|usrName=admin	result=success	sev=3 | logger.exe
echo.

timeout /t 2 >nul

echo [4/5] Обработка файла с примерами (9 логов)
echo --------------------------------------------
logger.exe -input examples\sample_logs.txt | findstr "Обработка"
echo.

timeout /t 2 >nul

echo [5/5] Запуск примера API
echo --------------------------------------------
go run examples\api_example\main.go | findstr /C:"=== Лог" /C:"Обработано"
echo.

echo ================================================
echo    Демонстрация завершена!
echo ================================================
echo.
echo Документация:
echo   - README.md        - Основное описание
echo   - QUICKSTART.md    - Быстрый старт
echo   - API.md           - API документация
echo   - COMMANDS.md      - Справочник команд
echo.
pause
