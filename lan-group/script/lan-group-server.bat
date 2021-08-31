@echo off
set hostname=
cls
if not "%hostname%"=="" ( goto main
) else (
  echo set hostname
  pause
  exit
)
:main
if "%1" == "h" goto begin
mshta vbscript:createobject("wscript.shell").run("""%~nx0"" h",0)(window.close)&&exit
:begin
C:\TOOL\lan-group-windows-amd64.exe -s -h %hostname%
