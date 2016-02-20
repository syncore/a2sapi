cd ../../src/db
go test
cd ../../src/steam
go test
cd ../../src/web
go test
rm -rf ../../bin/test_temp
cd ../../build/nix
