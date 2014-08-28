tasklist|find /i "justcms.exe" || goto lock
taskkill /f /t /im justcms.exe
:lock
start /b run.vbs
