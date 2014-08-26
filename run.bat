tasklist|find /i "JustCms.exe" || goto lock
taskkill /f /t /im JustCms.exe
:lock
start /b run.vbs
