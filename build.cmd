@echo off
echo:

if [%1]==[] goto :inv
if [%1]==[/help] goto :help
if [%1]==[/?] goto :help

if [%1]==[-b] goto :build
goto :inv

:build
cd src
echo running application go build (%2)
if [%2]==[stripped] (
    go build -o ../build/im-next.exe -ldflags '-s'
    goto :done
)
go build -o ../build/im-next.exe
goto :done

:inv
echo Invalid Arguments! See Help Page Below 

:help
echo ===========================HELP============================
echo build.cmd [-b build configuration (full/stripped)] 
goto :void

:done
echo done building

:void
echo: