rem
rem This is a simple batch file to build newt on windows.
rem
@echo off
SET NEWT_PROJECT=newt
jq -j -r .version codemeta.json >NEWT_VERSION
SET /p NEWT_VERSION=<NEWT_VERSION
SET NEWT_RELEASE_DATE=%DATE%
git log --pretty=format:"%h" -n 1 >>NEWT_RELEASE_HASH
SET /p NEWT_RELEASE_HASH=<NEWT_RELEASE_HASH
del /Q NEWT_*

@echo on

echo | pandoc --from t2t --to plain ^
--metadata-file codemeta.json ^
--metadata package=%NEWT_PROJECT% ^
--metadata version=%NEWT_VERSION% ^
--metadata release_date="%NEWT_RELEASE_DATE%" ^
--metadata release_hash=%NEWT_RELEASE_HASH% ^
--template codemeta-version-go.tmpl ^
--output version.go ^
LICENSE

go fmt version.go

echo.
echo Buidling bin\newt*.exe
echo.
mkdir bin
go build -o bin\newt.exe cmd\newt\newt.go
go build -o bin\nte.exe cmd\ndr\ndr.go
go build -o bin\ntr.exe cmd\nte\nte.go
