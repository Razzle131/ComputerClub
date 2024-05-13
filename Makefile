.SILENT:

name:=bin
build:
	go build -o ./binaries/$(name) .

run:
	go run main.go

test:
	go test ./...

image:=helloworldclub
docker-build:
	sudo docker build -t $(image) .

image:=helloworldclub
file:=test_file.txt
docker-run:
	sudo docker run -v $(dir $(realpath $(lastword $(MAKEFILE_LIST)))):/mnt/data/ $(image) /mnt/data/$(file)
