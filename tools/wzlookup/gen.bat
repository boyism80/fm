@echo off
setlocal
cd /d %~dp0\..\..
go run ./tools/wzlookup -wz resources/wz -out wz_lookup %*
endlocal
