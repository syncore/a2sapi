@echo off
mkdir %cd%\..\..\bin
del %cd%\..\..\bin\a2sapi.exe
go get github.com/fatih/color
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
go get github.com/oschwald/maxminddb-golang
go get github.com/stretchr/testify/assert
go build -i %cd%\..\..\src\a2sapi.go
move /Y a2sapi.exe %cd%\..\..\bin\
cd %cd%\..\..\bin