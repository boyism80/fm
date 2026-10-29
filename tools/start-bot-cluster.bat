@echo off
setlocal

rem start-bot-cluster.bat [agents=3] [seats=4] [listen=127.0.0.1:8090] [token]
set "AGENTS=%~1"
if "%AGENTS%"=="" set "AGENTS=3"
set "SEATS=%~2"
if "%SEATS%"=="" set "SEATS=4"
set "LISTEN=%~3"
if "%LISTEN%"=="" set "LISTEN=127.0.0.1:8090"
set "TOKEN=%~4"

cd /d "%~dp0..\bot"
if errorlevel 1 (
    echo Failed to change directory to bot.
    exit /b 1
)

powershell -NoProfile -Command "$c = New-Object Net.Sockets.TcpClient; try { $c.Connect('127.0.0.1', 8484); exit 0 } catch { exit 1 } finally { $c.Dispose() }"
if errorlevel 1 (
    echo login 127.0.0.1:8484 is not listening. Run tools\run-bot-servers.bat first.
    exit /b 1
)

echo Building bot...
go build -o out\bot.exe ..\cmd\bot
if errorlevel 1 exit /b 1

set "TOKEN_ARG="
if not "%TOKEN%"=="" set "TOKEN_ARG=-token %TOKEN%"

echo Starting master on %LISTEN%
start "FM Bot Master" cmd /k "out\bot.exe -mode master -config ..\config\bot.yaml -listen %LISTEN% %TOKEN_ARG%"

for /f "tokens=1,2 delims=:" %%A in ("%LISTEN%") do (
    set "HOST=%%A"
    set "PORT=%%B"
)
if "%HOST%"=="0.0.0.0" set "HOST=127.0.0.1"
powershell -NoProfile -Command "for ($i = 0; $i -lt 120; $i++) { $c = New-Object Net.Sockets.TcpClient; try { $c.Connect('%HOST%', %PORT%); exit 0 } catch { Start-Sleep -Milliseconds 500 } finally { $c.Dispose() } }; exit 1"
if errorlevel 1 (
    echo master did not open %LISTEN%
    exit /b 1
)

for /L %%N in (1,1,%AGENTS%) do (
    echo Starting agent-%%N with %SEATS% seats
    start "FM Bot Agent %%N" cmd /k "out\bot.exe -mode agent -config ..\config\bot.yaml -master ws://%HOST%:%PORT% -name agent-%%N -seats %SEATS% %TOKEN_ARG%"
)

start "" "http://%HOST%:%PORT%"
endlocal
exit /b 0
