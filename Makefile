send: weatherpi
	scp ./weatherpi pi@192.168.1.8:~/weatherpi
build: main.go
	rm weatherpi
	GOOS=linux GOARCH=arm GOARM=5 go build
