set dir=C:\TOOL\
if not exist %dir% (
  mkdir %dir%
)
copy lan-group-server.vbs "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp"
copy lan-group-windows-amd64.exe %dir%
