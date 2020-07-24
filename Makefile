build:
	env GOOS=windows GOARCH=amd64 go build -o bin/wpp_broadcast_message_windows64x.exe ./main.go
	env GOOS=windows GOARCH=386 go build -o bin/wpp_broadcast_message_windows86x.exe ./main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/wpp_broadcast_message_linux_amd64 ./main.go
	env GOOS=linux GOARCH=386 go build -o bin/wpp_broadcast_message_linux_386 ./main.go
	env GOOS=linux GOARCH=arm go build -o bin/wpp_broadcast_message_linux_arm ./main.go

run_windows_64:
	.\bin\wpp_broadcast_message_windows64x