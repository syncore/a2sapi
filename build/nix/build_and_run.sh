mkdir ../../bin/
rm -rf ../../bin/a2sapi
go get github.com/fatih/color
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
go get github.com/oschwald/maxminddb-golang
go build -i ../../src/a2sapi.go
mv a2sapi ../../bin/
cd ../../bin/
./a2sapi