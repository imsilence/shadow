@echo off

cd %~dp0

cd ..
setlocal

set HOME=%cd%
echo WORKON: %HOME%

set GOPATH=%cd%

set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -o bin\agent.exe -i agent

if "%ERRORLEVEL%" == "0" (
    echo success build win x64
) else (
    echo error build win x64
)

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o bin\agent -i agent

if "%ERRORLEVEL%" == "0" (
    echo success build linux x64
) else (
    echo error build linux x64
)

endlocal
