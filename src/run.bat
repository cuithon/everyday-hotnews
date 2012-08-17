:: Copyright 2012 The Go Authors. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.
@echo off

:: Keep environment variables within this script
:: unless invoked with --no-local.
if x%1==x--no-local goto nolocal
if x%2==x--no-local goto nolocal
setlocal
:nolocal

set GOBUILDFAIL=0

:: we disallow local import for non-local packages, if %GOROOT% happens
:: to be under %GOPATH%, then some tests below will fail
set GOPATH=

rem TODO avoid rebuild if possible

if x%1==x--no-rebuild goto norebuild
echo # Building packages and commands.
go install -a -v std
if errorlevel 1 goto fail
echo.
:norebuild

:: we must unset GOROOT_FINAL before tests, because runtime/debug requires
:: correct access to source code, so if we have GOROOT_FINAL in effect,
:: at least runtime/debug test will fail.
set GOROOT_FINAL=

:: get CGO_ENABLED
go env > env.bat
if errorlevel 1 goto fail
call env.bat
del env.bat
echo.

echo # Testing packages.
go test std -short -timeout=120s
if errorlevel 1 goto fail
echo.

echo # runtime -cpu=1,2,4
go test runtime -short -timeout=120s -cpu=1,2,4
if errorlevel 1 goto fail
echo.

echo # sync -cpu=10
go test sync -short -timeout=120s -cpu=10
if errorlevel 1 goto fail
echo.

echo # ..\misc\dashboard\builder ..\misc\goplay
go build ..\misc\dashboard\builder ..\misc\goplay
if errorlevel 1 goto fail
echo.

:: TODO(brainman): disabled, because it fails with: mkdir C:\Users\ADMINI~1\AppData\Local\Temp\2.....\test\bench\: The filename or extension is too long.
::echo # ..\test\bench\go1
::go test ..\test\bench\go1
::if errorlevel 1 goto fail
::echo.

:: cgo tests
if x%CGO_ENABLED% == x0 goto nocgo
echo # ..\misc\cgo\life
go run %GOROOT%\test\run.go - ..\misc\cgo\life
if errorlevel 1 goto fail
echo.

:: TODO ..\misc\cgo\stdio

:: TODO(brainman): disabled, because it fails with: mkdir C:\Users\ADMINI~1\AppData\Local\Temp\2.....\go\misc\cgo\: The filename or extension is too long.
::echo # ..\misc\cgo\test
::go test ..\misc\cgo\test
::if errorlevel 1 goto fail
::echo.
:nocgo

:: TODO: The other tests in run.bash.

echo # test
cd ..\test
set FAIL=0
go run run.go
if errorlevel 1 set FAIL=1
cd ..\src
echo.
if %FAIL%==1 goto fail

echo # Checking API compatibility.
go tool api -c ..\api\go1.txt -next ..\api\next.txt
if errorlevel 1 goto fail
echo.

echo ALL TESTS PASSED
goto end

:fail
set GOBUILDFAIL=1

:end
