@echo off
rem Quick launcher - opens Traverse & Level in its own window using Edge or Chrome.
rem No installation needed: it opens the web app, so you log in with your Excel to KML account.
set "APP=https://www.exceltokml.com/traverse/app"
set "PROF=%LOCALAPPDATA%\TraverseLevel\browserprofile"

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
start "" "%EDGE%" --app="%APP%" --user-data-dir="%PROF%" --window-size=1440,920
exit /b 0
