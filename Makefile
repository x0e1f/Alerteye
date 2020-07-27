all:
	go build -i -o ./bin/alerteye main.go 
	upx ./bin/alerteye
rpi:	
	env CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -i -o ./bin/alerteye_arm main.go 
	upx ./bin/alerteye_arm	
clean:
	rm -rf ./bin