tasklist|find /i "main.exe" || goto lock
taskkill /f /t /im main.exe
:lock
go build main.go
start cmd /c run.vbs
