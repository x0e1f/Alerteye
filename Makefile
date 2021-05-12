all:
	go build -o ./bin/alerteye ./cmd/alerteye/main.go 
rpi:	
	env CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o ./bin/alerteye_rpi ./cmd/alerteye/main.go
clean:
	rm -rf ./bin