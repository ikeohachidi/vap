build:
	GOOS=windows GOARCH=amd64 go build -o bin/vap.exe main.go
	go build -o bin/vap main.go