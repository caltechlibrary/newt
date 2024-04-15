rem
rem This is a simple batch file to build newt on windows.
rem
jq .version codemeta.json >NEWT_VERSION
SET /p NEWT_VERSION=<NEWT_VERSION
date /T >NEWT_RELEASE_DATE
SET /p NEWT_RELEASE_DATE=<NEWT_RELEASE_DATE
git log --pretty=format:'%h' -n 1 >>NEWT_RELEASE_HASH
SET /p NEWT_RELEASE_HASH=<NEWT_RELEASE_HASH
del /Q NEWT_*

@echo off
echo package newt >version.go
echo. >>version.go
echo import ( >>version.go
echo     "strings" >>version.go
echo ) >>version.go
echo. >>version.go
echo const ( >>version.go
echo    // Version of newt package >>version.go
echo    Version = %NEWT_VERSION% >>version.go 
echo    // ReleaseDate, the date version.go was generated >>version.go
echo    ReleaseDate = "%NEWT_RELEASE_DATE%" >>version.go
echo. >>version.go
echo    // ReleaseHash, the Git hash when version.go was generated >>version.go
echo    ReleaseHash = "%NEWT_RELEASE_HASH%" >>version.go
echo. >>version.go
echo    LicenseText = `>>version.go
echo. >>version.go
type LICENSE >>version.go
echo. >>version.go
echo `>>version.go
echo ) >>version.go
echo. >>version.go
echo. >>version.go
echo // FmtHelp lets you process a text block with simple curly brace markup. >>version.go
echo func FmtHelp(src string, appName string, version string, releaseDate string, releaseHash string) string { >>version.go
echo	m := map[string]string { >>version.go
echo		"{app_name}": appName, >>version.go
echo		"{version}": version, >>version.go
echo		"{release_date}": releaseDate, >>version.go
echo		"{release_hash}": releaseHash, >>version.go
echo	} >>version.go
echo	for k, v := range m { >>version.go
echo		if strings.Contains(src, k) { >>version.go
echo			src = strings.ReplaceAll(src, k, v) >>version.go
echo		} >>version.go
echo 	} >>version.go
echo	return src >>version.go
echo } >>version.go
@echo on
go fmt version.go

echo.
echo Buidling bin\newt.exe
echo.
mkdir bin
go build -o bin\newt.exe cmd\newt\newt.go
go build -o bin\newtmustache.exe cmd\newtmustache\newtmustache.go
go build -o bin\newtrouter.exe cmd\newtwrouter\newtrouter.go
go build -o bin\newtgenerator.exe cmd\newtgenerator\newtgenerator.go
go build -o bin\ws.exe cmd\ws\ws.go
go build -o bin\mustache.exe cmd\mustache\mustache.go
