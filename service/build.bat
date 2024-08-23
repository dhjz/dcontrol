:: 后台运行程序exe
::go build -ldflags "-s -w -H=windowsgui"
:: go build -o ddns.exe
::cp -Force -r ../web/dist/* ./webapp/
:: rmdir /S /Q .\webapp
:: xcopy /E /Y "..\web\dist\*" ".\webapp\"

:: go env -w GOOS=linux
:: go build -ldflags "-s -w" -o ./dist/

go env -w GOOS=windows
go build -ldflags "-s -w -H=windowsgui" -o ./dist/
:: go build -o ddns.exe

::cp -r ./webapp ./dist/

pause