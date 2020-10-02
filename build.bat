mkdir build
set GOOS=windows
go build -o ./build/AutoUpdater.exe
set GOOS=darwin
go build -o ./build/AutoUpdater
strip ./build/AutoUpdater.exe
upx ./build/AutoUpdater.exe
upx ./build/AutoUpdater