# ComputerClubApp
______
# Table of Contents
* [About](#about)
* [Clone repo](#clone-repo)
* [Technologies](#technologies)
* [Quick start](#quick-start)

## About
This app is made for simulation of computer club system, which proccess events in file that is set at startup.

## Clone repo
You can clone this repo to observe source code and run tests localy using this command (golang version 1.22.2 or higher):
```
git clone https://github.com/Razzle131/ComputerClub.git
```

## Technologies
Project is created with:
* Golang version: 1.22.2

## Quick start
### Docker
This app could be builded and executed using docker:
* Clone repo
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* use this command to build app docker image  
```
docker build -t IMAGE-NAME .
```
replace the IMAGE-NAME with the name you want, memorize it to execute builded image later

* to execute programm use this programm, you will need to use volume to access file with source data
```
docker run -v /path/to/test/files/folder:/mnt/data IMAGE-NAME /mnt/data/testFilename
```
Example: `docker run -v /home/myuser/myTestFiles/:/mnt/data foo /mnt/data/test_file`  
Assuming the example we connect folder /home/myuser/myTestFiles/ to our image, named "foo", and giving it as argument path inside volume with filename test_file  
______
### Build from source
#### Without Makefile
* Ensure that Go is installed on your machine and it`s version is equal or higther than 1.22.2
```
go version
```
* Clone repo
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* Run code with example data
```
go run main.go test_file.txt
```
______
#### Using Makefile
* Ensure that Go is installed on your machine and it`s version is equal or higther than 1.22.2
```
go version
```
* Clone repo
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* Now you can use commands `make <command>`
  * `make build` - builds source code into binary in binaries folder with the name set on bin by default
  * `make build name=foo` - builds source code into binary with name setted foo
  * `make run` - does "go run main.go"
  * `make test` - runs all unit tests
  * `docker-build` - builds docker image with name helloworldclub by default
  * `docker-build image=foo` - builds docker image with name "foo"
  * `make docker-run` - runs compiled image, sets image=helloworldclub and file=test_file.txt by default, you can change them adding after "...run" tags image=<yourimagename> and file=<yourtestfile>
