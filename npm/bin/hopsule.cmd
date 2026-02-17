@echo off
:: Hopsule CLI - Windows npm wrapper
:: Launches the Go binary installed by the npm postinstall script.

setlocal

set "BASEDIR=%~dp0"
set "BINARY=%BASEDIR%hopsule.exe"

if exist "%BINARY%" (
    "%BINARY%" %*
    exit /b %ERRORLEVEL%
)

echo Error: hopsule binary not found. 1>&2
echo Try reinstalling: npm install -g hopsule 1>&2
exit /b 1
