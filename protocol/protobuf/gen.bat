@echo off
REM Regenerate Go under protocol\protobuf\gengo\ from sources in protocol\protobuf\internal\proto\.
REM (Output avoids protocol/.../internal/go/... so other modules can import generated types without Go internal rules.)
REM Regenerate JS under services\internal\protobuf\.
REM Requires: protoc, protoc-gen-go, protoc-gen-go-grpc on PATH; npm install in services\internal.
setlocal EnableExtensions
cd /d "%~dp0"

pushd "%~dp0..\.." >nul
set "REPO_ROOT=%CD%"
popd >nul

set "PROTO_SRC=%~dp0internal\proto"
set "GO_OUT=%~dp0gengo"
set "JS_OUT=%REPO_ROOT%\services\internal\protobuf"
set "GRPC_TOOLS=%REPO_ROOT%\services\internal\node_modules\.bin\grpc_tools_node_protoc.cmd"

where protoc >nul 2>&1
if errorlevel 1 (
  echo [gen.bat] ERROR: protoc not found in PATH.
  echo Install: https://grpc.io/docs/protoc-installation/
  exit /b 1
)

where protoc-gen-go >nul 2>&1
if errorlevel 1 (
  echo [gen.bat] ERROR: protoc-gen-go not found in PATH.
  echo Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  exit /b 1
)

where protoc-gen-go-grpc >nul 2>&1
if errorlevel 1 (
  echo [gen.bat] ERROR: protoc-gen-go-grpc not found in PATH.
  echo Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  exit /b 1
)

if not exist "%GRPC_TOOLS%" (
  echo [gen.bat] ERROR: grpc_tools_node_protoc not found.
  echo Run: cd services\internal ^&^& npm install
  exit /b 1
)

if not exist "%GO_OUT%" mkdir "%GO_OUT%"
if not exist "%JS_OUT%" mkdir "%JS_OUT%"

echo [gen.bat] Generating Go under "%GO_OUT%" ...
for /r "%PROTO_SRC%" %%F in (*.proto) do (
  echo   %%F
  protoc ^
    -I "%PROTO_SRC%" ^
    --go_out="%GO_OUT%" ^
    --go_opt=paths=source_relative ^
    --go-grpc_out="%GO_OUT%" ^
    --go-grpc_opt=paths=source_relative ^
    "%%F"
  if errorlevel 1 exit /b 1
)

echo [gen.bat] Generating JS under "%JS_OUT%" ...
for /r "%PROTO_SRC%" %%F in (*.proto) do (
  echo   %%F
  call "%GRPC_TOOLS%" ^
    -I "%PROTO_SRC%" ^
    --js_out=import_style=commonjs,binary:"%JS_OUT%" ^
    --grpc_out=grpc_js:"%JS_OUT%" ^
    "%%F"
  if errorlevel 1 exit /b 1
)

echo [gen.bat] Done.
exit /b 0
