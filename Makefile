build:
	go build -o wfp-cli .

run:
	go run main.go

compile:
	GOOS=linux GOARCH=amd64 go build -o wfp-cli .