@echo off
setlocal EnableDelayedExpansion

cd /d "%~dp0.."
if errorlevel 1 (
    echo Failed to change directory to repo root.
    exit /b 1
)

call "%~dp0run-servers.bat" %*
if errorlevel 1 exit /b 1

set "GAME_TEMPLATE=%~2"
if "%GAME_TEMPLATE%"=="" set "GAME_TEMPLATE=config\game.yaml"
if not exist "%GAME_TEMPLATE%" (
    echo Game template not found: %GAME_TEMPLATE%
    exit /b 1
)

if not exist "config\dev" mkdir "config\dev"

set "CH_CFG=config\dev\game-ch0.yaml"
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0generate-game-channel-config.ps1" -TemplatePath "%GAME_TEMPLATE%" -ChannelId 0 -OutputPath "%CH_CFG%"
if errorlevel 1 exit /b 1

echo Starting game server channel 0 with %CH_CFG%...
start "FM Game Ch0" cmd /k "cd /d services\game && go run . -config=%CH_CFG%"

echo.
echo Bot servers started, including game channel 0.
endlocal
