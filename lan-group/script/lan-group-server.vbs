hostname=""

If StrComp(hostname,"") = 0 then
  MsgBox( "Lan-Group-Server: Set hostname" )
Else
  Set ws = CreateObject("Wscript.Shell")
  ws.run "cmd /c C:\TOOL\lan-group-windows-amd64.exe -s -h "&hostname ,0
End If
