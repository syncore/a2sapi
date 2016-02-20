mkdir -p ../../bin/
rm -rf ../../bin/a2sapi
go get -u github.com/fatih/color
go get -u github.com/gorilla/mux
go get -u github.com/mattn/go-sqlite3
go get -u github.com/oschwald/maxminddb-golang
go get -u github.com/stretchr/testify/assert
go build -i ../../src/a2sapi.go
mv a2sapi ../../bin/
cd ../../bin/
