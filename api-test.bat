@ECHO OFF

go-testbinary.exe -test.v -test.run ^TestApi$ > plain.out
test-summary.exe

DEL plain.out

ECHO test executed...
