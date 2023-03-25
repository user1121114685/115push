go get -u github.com/deadblue/elevengo@master
set currenttime=%date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%
DEL 115push_*
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_windows_amd64.exe
SET GOOS=linux
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_linux_amd64
SET GOARCH=arm
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_linux_arm
SET GOARCH=arm64
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_linux_arm64
SET GOOS=darwin
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_Mac_arm64
SET GOARCH=amd64
go build -ldflags="-s -w -X 'main.Version=%currenttime%'" -o 115push_Mac_amd64
D:\7-Zip\7z.exe a 115push_all.zip 115push_*