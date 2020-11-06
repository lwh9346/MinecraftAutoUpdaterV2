mkdir build
set GOOS=windows
go build -o ./build/AutoUpdater.exe ./updater/main.go