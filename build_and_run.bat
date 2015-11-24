@echo off
del steamtest.exe
cls
go build -i %cd%\src\cmd\steamtest.go
steamtest.exe