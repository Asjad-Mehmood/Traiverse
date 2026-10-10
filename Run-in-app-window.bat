@echo off
rem Quick launcher - opens Traverse Pro in its own window using Edge or Chrome.
rem No installation needed. Works even before you build the .exe.
set "APP=%~dp0app\index.html"
set "PROF=%LOCALAPPDATA%\TraversePro\browserprofile"

set "EDGE=%ProgramFiles(x86)%\Microsoft\Edge\Application\msedge.exe"
if exist "%EDGE%" goto run
set "EDGE=%ProgramFiles%\Microsoft\Edge\Application\msedge.exe"
if exist "%EDGE%" goto run
set "EDGE=%ProgramFiles%\Google\Chrome\Application\chrome.exe"
if exist "%EDGE%" goto run
set "EDGE=%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe"
if exist "%EDGE%" goto run

start "" "%APP%"
exit /b 0

:run
start "" "%EDGE%" --app="file:///%APP:\=/%" --user-data-dir="%PROF%" --window-size=1240,840
exit /b 0
