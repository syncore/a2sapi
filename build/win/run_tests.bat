@echo off
cd %cd%\..\..\src\db
go test
cd %cd%\..\..\src\steam
go test
cd %cd%\..\..\src\web
go test
rmdir /S /Q %cd%\..\..\bin\test_temp
cd %cd%\..\..\build\win