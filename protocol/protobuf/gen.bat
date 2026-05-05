@echo off
REM Regenerate Go under protocol\protobuf\gengo\ from sources in protocol\protobuf\internal\proto\.
REM (Output avoids protocol/.../internal/go/... so other modules can import generated types without Go internal rules.)
REM Regenerate TypeScript source under services\internal\src\protobuf\generated\.
REM Requires: protoc, protoc-gen-go, protoc-gen-go-grpc on PATH; npm install in services\internal.
setlocal EnableExtensions
cd /d "%~dp0"

pushd "%~dp0..\.." >nul
set "REPO_ROOT=%CD%"
popd >nul

set "PROTO_SRC=%~dp0internal\proto"
set "GO_OUT=%~dp0gengo"
set "TS_PROTO_OUT=%REPO_ROOT%\services\internal\src\protobuf\generated"
set "GRPC_TOOLS=%REPO_ROOT%\services\internal\node_modules\.bin\grpc_tools_node_protoc.cmd"
set "TS_PROTO_PLUGIN=%REPO_ROOT%\services\internal\node_modules\.bin\protoc-gen-ts_proto.cmd"

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

if not exist "%TS_PROTO_PLUGIN%" (
  echo [gen.bat] ERROR: protoc-gen-ts_proto not found.
  echo Run: cd services\internal ^&^& npm install
  exit /b 1
)

if not exist "%GO_OUT%" mkdir "%GO_OUT%"
if not exist "%TS_PROTO_OUT%" mkdir "%TS_PROTO_OUT%"

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

echo [gen.bat] Generating TS source under "%TS_PROTO_OUT%" ...
for /r "%PROTO_SRC%" %%F in (*.proto) do (
  echo   %%F
  call "%GRPC_TOOLS%" ^
    -I "%PROTO_SRC%" ^
    --plugin=protoc-gen-ts_proto="%TS_PROTO_PLUGIN%" ^
    --ts_proto_out="%TS_PROTO_OUT%" ^
    --ts_proto_opt=outputServices=grpc-js,env=node,esModuleInterop=true,useOptionals=none ^
    "%%F"
  if errorlevel 1 exit /b 1
)

echo [gen.bat] Done.
exit /b 0
