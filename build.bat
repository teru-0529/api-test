@ECHO OFF

go test -c -o go-testbinary.exe
go build -o test-summary.exe

ECHO build completed...
