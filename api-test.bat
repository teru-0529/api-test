@ECHO OFF

set time2=%time: =0%
SET ts=%date:~0,4%%date:~5,2%%date:~8,2%-%time2:~0,2%%time2:~3,2%%time2:~6,2%
ECHO %ts%

go-testbinary.exe -test.v -test.run ^TestApi$ > plain.out
@REM go-testbinary.exe -test.v > plain.out

ECHO test executed...
