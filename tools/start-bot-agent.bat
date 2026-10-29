@echo off
setlocal

rem start-bot-agent.bat <ws://master:8090> [token] [name=computer name] [seats=4]
rem Joins a master on another PC. The master needs -listen 0.0.0.0:8090 and the same token.
set "MASTER=%~1"
if "%MASTER%"=="" (
    echo usage: start-bot-agent.bat ws://master:8090 [token] [name] [seats]
    exit /b 1
)
set "TOKEN=%~2"
set "NAME=%~3"
if "%NAME%"=="" set "NAME=%COMPUTERNAME%"
set "SEATS=%~4"
if "%SEATS%"=="" set "SEATS=4"

cd /d "%~dp0..\bot"
if errorlevel 1 (
    echo Failed to change directory to bot.
    exit /b 1
)

echo Building bot...
go build -o out\bot.exe ..\cmd\bot
if errorlevel 1 exit /b 1

set "TOKEN_ARG="
if not "%TOKEN%"=="" set "TOKEN_ARG=-token %TOKEN%"

start "FM Bot Agent %NAME%" cmd /k "out\bot.exe -mode agent -config ..\config\bot.yaml -master %MASTER% -name %NAME% -seats %SEATS% %TOKEN_ARG%"
endlocal
exit /b 0
