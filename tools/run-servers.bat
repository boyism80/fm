@echo off
setlocal EnableDelayedExpansion

cd /d "%~dp0.."
if errorlevel 1 (
    echo Failed to change directory to repo root.
    exit /b 1
)

set "LOGIN_CFG=%~1"
set "GAME_TEMPLATE=%~2"
set "GC1=%~3"
set "GC2=%~4"
set "GC3=%~5"

if "%GAME_TEMPLATE%"=="" set "GAME_TEMPLATE=config\game.yaml"
if not exist "%GAME_TEMPLATE%" (
    echo Game template not found: %GAME_TEMPLATE%
    exit /b 1
)

if not exist "config\dev" mkdir "config\dev"

set "LOGIN_CMD=go run ./services/login/main.go"
if not "%LOGIN_CFG%"=="" (
    set "LOGIN_CMD=!LOGIN_CMD! -config=%LOGIN_CFG%"
)

echo Starting login server...
start "FM Login" cmd /k !LOGIN_CMD!

set "IDX=0"
for %%C in (1 2 3) do (
    set /a IDX+=1
    set "CH_CFG="
    if !IDX!==1 if not "!GC1!"=="" set "CH_CFG=!GC1!"
    if !IDX!==2 if not "!GC2!"=="" set "CH_CFG=!GC2!"
    if !IDX!==3 if not "!GC3!"=="" set "CH_CFG=!GC3!"
    if "!CH_CFG!"=="" (
        set "CH_CFG=config\dev\game-ch%%C.yaml"
        powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0generate-game-channel-config.ps1" -TemplatePath "%GAME_TEMPLATE%" -ChannelId %%C -OutputPath "!CH_CFG!"
        if errorlevel 1 exit /b 1
    )
    echo Starting game server channel %%C with !CH_CFG!...
    start "FM Game Ch%%C" cmd /k go run ./services/game/main.go -config=!CH_CFG!
)

echo.
echo All servers started in separate windows.
echo   Login:  optional arg 1 = login config ^(default config/login.yaml^)
echo   Game:   optional arg 2 = template for generated configs ^(default config/game.yaml^)
echo           optional args 3-5 = per-channel game configs ^(skip generation^)
echo.
echo Ensure services/internal/config.yaml lists ports 8486/8487/8488 for channels 1-3 on world 0.
endlocal
