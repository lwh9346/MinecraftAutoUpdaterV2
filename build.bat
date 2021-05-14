mkdir program
set GOOS=windows
go build -o ./program/AutoUpdater.exe ./updater/main.go
go build -o ./program/SelfUpdater.exe ./selfupdate/main.go
strip ./program/AutoUpdater.exe
strip ./program/SelfUpdater.exe
upx -9 ./program/SelfUpdater.exe
upx -9 ./program/AutoUpdater.exe
program\AutoUpdater.exe hash