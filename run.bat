tasklist|find /i "JustCms.exe" || goto lock
taskkill /f /t /im JustCms.exe
:lock
go build
start /b run.vbs
