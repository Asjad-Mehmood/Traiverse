; Installer for Traverse Pro.
Unicode true
!ifndef VERSION
  !define VERSION "1.0.0"
!endif
Name "Traverse Pro"
OutFile "TraversePro-Setup.exe"
; keep these names: every update installs in place over the same folder and keys
InstallDir "$LOCALAPPDATA\TraversePro"
InstallDirRegKey HKCU "Software\TraversePro" "InstallDir"
RequestExecutionLevel user
SetCompressor /SOLID lzma
ShowInstDetails show

!include "MUI2.nsh"
!define MUI_ICON "icon.ico"
!define MUI_UNICON "icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\TraversePro.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Open Traverse Pro now"

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "Traverse Pro"
VIAddVersionKey "FileDescription" "Traverse Pro installer"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "Copyright (c) 2026 Asjad Mehmood"

Section "Install"
  ; Windows will not overwrite a running program, so clear the way first
  StrCpy $R9 0
  retry_running:
    ClearErrors
    Delete "$INSTDIR\TraversePro.exe"
    IfErrors 0 clear_running
    IntOp $R9 $R9 + 1
    IntCmp $R9 10 giveup 0 giveup
    MessageBox MB_ABORTRETRYIGNORE|MB_ICONEXCLAMATION \
      "Traverse Pro is still open, so it cannot be replaced.$\r$\n$\r$\nClose it and press Retry. Retry will also try to close it for you." \
      /SD IDRETRY IDRETRY kill_running IDIGNORE clear_running
  giveup:
    Abort "Installation stopped. Close Traverse Pro and run this again."
  kill_running:
    nsExec::Exec 'taskkill /IM TraversePro.exe /F'
    Pop $0
    Sleep 900
    Goto retry_running
  clear_running:

  SetOutPath "$INSTDIR"
  File "TraversePro.exe"
  File "icon.ico"

  CreateDirectory "$SMPROGRAMS\Traverse Pro"
  CreateShortcut "$SMPROGRAMS\Traverse Pro\Traverse Pro.lnk" \
                 "$INSTDIR\TraversePro.exe" "" "$INSTDIR\icon.ico" 0
  CreateShortcut "$DESKTOP\Traverse Pro.lnk" \
                 "$INSTDIR\TraversePro.exe" "" "$INSTDIR\icon.ico" 0

  WriteRegStr HKCU "Software\TraversePro" "InstallDir" "$INSTDIR"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayName" "Traverse Pro"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "Publisher" "Asjad Mehmood"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayIcon" "$INSTDIR\icon.ico"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "NoModify" 1
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; silent run = started by the in-app updater, so reopen the program
  IfSilent 0 skip_launch
    Exec '"$INSTDIR\TraversePro.exe"'
  skip_launch:
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\TraversePro.exe"
  Delete "$INSTDIR\icon.ico"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  Delete "$SMPROGRAMS\Traverse Pro\Traverse Pro.lnk"
  RMDir "$SMPROGRAMS\Traverse Pro"
  Delete "$DESKTOP\Traverse Pro.lnk"
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro"
  DeleteRegKey HKCU "Software\TraversePro"
SectionEnd
