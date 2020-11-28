mkdir build
set GOOS=windows
go build -o ./build/AutoUpdater.exe ./updater/main.go
go build -o ./build/SelfUpdater.exe ./selfupdate/main.go
strip ./build/AutoUpdater.exe
strip ./build/SelfUpdater.exe
upx -9 ./build/SelfUpdater.exe
upx -9 ./build/AutoUpdater.exe