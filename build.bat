@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "ROOT=%~dp0"
if "%ROOT:~-1%"=="\" set "ROOT=%ROOT:~0,-1%"
set "BUILD_ROOT=%ROOT%\build"
set "TARGET=%~1"
if "%TARGET%"=="" set "TARGET=all"

set "GOOS=linux"
set "GOARCH=amd64"
set "CGO_ENABLED=0"

if not exist "%BUILD_ROOT%" mkdir "%BUILD_ROOT%"

if /I "%TARGET%"=="all" (
    call :build_go_service "gofurry-nav-backend" || exit /b 1
    call :build_go_service "gofurry-nav-collector" || exit /b 1
    call :build_go_service "gofurry-game-backend" || exit /b 1
    call :build_go_service "gofurry-game-collector" || exit /b 1
    call :build_admin || exit /b 1
    echo Build completed. Artifacts are in "%BUILD_ROOT%".
    exit /b 0
)

if /I "%TARGET%"=="gofurry-nav-backend" (
    call :build_go_service "gofurry-nav-backend" || exit /b 1
    exit /b 0
)
if /I "%TARGET%"=="gofurry-nav-collector" (
    call :build_go_service "gofurry-nav-collector" || exit /b 1
    exit /b 0
)
if /I "%TARGET%"=="gofurry-game-backend" (
    call :build_go_service "gofurry-game-backend" || exit /b 1
    exit /b 0
)
if /I "%TARGET%"=="gofurry-game-collector" (
    call :build_go_service "gofurry-game-collector" || exit /b 1
    exit /b 0
)
if /I "%TARGET%"=="gofurry-admin" (
    call :build_admin || exit /b 1
    exit /b 0
)

echo Unknown target: %TARGET%
echo Supported targets:
echo   all
echo   gofurry-nav-backend
echo   gofurry-nav-collector
echo   gofurry-game-backend
echo   gofurry-game-collector
echo   gofurry-admin
exit /b 1

:build_admin
echo [BUILD] gofurry-admin web
call :build_frontend_in_place "gofurry-admin\web" || exit /b 1
echo [BUILD] gofurry-admin binary
call :build_go_binary "gofurry-admin" "gofurry-admin" || exit /b 1
call :copy_dir "gofurry-admin\internal\transport\http\webui\dist" "gofurry-admin\dist" || exit /b 1
exit /b 0

:build_go_service
set "SERVICE=%~1"
echo [BUILD] %SERVICE%
call :build_go_binary "%SERVICE%" "%SERVICE%" || exit /b 1
exit /b 0

:build_frontend_in_place
set "SERVICE_REL=%~1"
call :npm_ci_and_build "%ROOT%\%SERVICE_REL%" || exit /b 1
exit /b 0

:build_go_binary
set "SERVICE=%~1"
set "BINARY_NAME=%~2"
set "OUTPUT_DIR=%BUILD_ROOT%\%SERVICE%"
set "OUTPUT_BIN=%OUTPUT_DIR%\%BINARY_NAME%"

if exist "%OUTPUT_DIR%" rmdir /s /q "%OUTPUT_DIR%"
mkdir "%OUTPUT_DIR%" || exit /b 1

pushd "%ROOT%\%SERVICE%" || exit /b 1
call go build -trimpath -ldflags="-s -w" -o "%OUTPUT_BIN%" .
if errorlevel 1 (
    popd
    exit /b 1
)
popd
exit /b 0

:npm_ci_and_build
set "WORKDIR=%~1"
pushd "%WORKDIR%" || exit /b 1
call npm ci
if errorlevel 1 (
    popd
    exit /b 1
)
call npm run build
if errorlevel 1 (
    popd
    exit /b 1
)
popd
exit /b 0

:copy_dir
set "SRC=%ROOT%\%~1"
set "DEST=%BUILD_ROOT%\%~2"

if not exist "%SRC%" (
    echo Source directory not found: "%SRC%"
    exit /b 1
)

if exist "%DEST%" rmdir /s /q "%DEST%"
mkdir "%DEST%" || exit /b 1
xcopy "%SRC%\*" "%DEST%\" /E /I /Y >nul
if errorlevel 1 exit /b 1
exit /b 0
