mkdir program
set GOOS=windows
go build -o ./program/AutoUpdater.exe ./updater/main.go
go build -o ./program/SelfUpdater.exe ./selfupdate/main.go
program\AutoUpdater.exe hash
set GOOS=linux
go build -o ./program/proxy ./pserver/main.go