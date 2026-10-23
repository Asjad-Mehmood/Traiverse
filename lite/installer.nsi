; Installer for Traverse & Level (Professional Traverse and Level Adjustment Software).
Unicode true
!ifndef VERSION
  !define VERSION "2.1.0"
!endif
Name "Traverse & Level"
OutFile "TraverseLevel-Setup.exe"
; Keep these names: every update installs in place over the same folder and keys. They are the
; ones Traverse Pro 1.x used, so a Traverse Pro install updates into Traverse & Level.
InstallDir "$LOCALAPPDATA\TraversePro"
InstallDirRegKey HKCU "Software\TraversePro" "InstallDir"
RequestExecutionLevel user
SetCompressor /SOLID lzma
ShowInstDetails show

!include "MUI2.nsh"
!define MUI_ICON "icon.ico"
!define MUI_UNICON "icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\TraverseLevel.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Open Traverse & Level now"

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "Traverse & Level"
VIAddVersionKey "FileDescription" "Traverse & Level installer"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "Copyright (c) 2026 Asjad Mehmood"

Section "Install"
  ; Windows will not overwrite a running program, so clear the way first
  ; (TraversePro.exe is the Traverse Pro 1.x program this replaces)
  StrCpy $R9 0
  retry_running:
    ClearErrors
    Delete "$INSTDIR\TraverseLevel.exe"
    Delete "$INSTDIR\TraversePro.exe"
    IfErrors 0 clear_running
    IntOp $R9 $R9 + 1
    IntCmp $R9 10 giveup 0 giveup
    MessageBox MB_ABORTRETRYIGNORE|MB_ICONEXCLAMATION \
      "Traverse & Level is still open, so it cannot be replaced.$\r$\n$\r$\nClose it and press Retry. Retry will also try to close it for you." \
      /SD IDRETRY IDRETRY kill_running IDIGNORE clear_running
  giveup:
    Abort "Installation stopped. Close Traverse & Level and run this again."
  kill_running:
    nsExec::Exec 'taskkill /IM TraverseLevel.exe /F'
    Pop $0
    nsExec::Exec 'taskkill /IM TraversePro.exe /F'
    Pop $0
    Sleep 900
    Goto retry_running
  clear_running:

  SetOutPath "$INSTDIR"
  File "TraverseLevel.exe"
  File "icon.ico"

  ; Traverse Pro 1.x shortcuts
  Delete "$SMPROGRAMS\Traverse Pro\Traverse Pro.lnk"
  RMDir "$SMPROGRAMS\Traverse Pro"
  Delete "$DESKTOP\Traverse Pro.lnk"

  CreateDirectory "$SMPROGRAMS\Traverse & Level"
  CreateShortcut "$SMPROGRAMS\Traverse & Level\Traverse & Level.lnk" \
                 "$INSTDIR\TraverseLevel.exe" "" "$INSTDIR\icon.ico" 0
  CreateShortcut "$DESKTOP\Traverse & Level.lnk" \
                 "$INSTDIR\TraverseLevel.exe" "" "$INSTDIR\icon.ico" 0

  WriteRegStr HKCU "Software\TraversePro" "InstallDir" "$INSTDIR"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayName" "Traverse & Level"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "Publisher" "SurveyX Digital"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "DisplayIcon" "$INSTDIR\icon.ico"
  WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro" \
              "NoModify" 1
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; silent run = started by the in-app updater, so reopen the program
  IfSilent 0 skip_launch
    Exec '"$INSTDIR\TraverseLevel.exe"'
  skip_launch:
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\TraverseLevel.exe"
  Delete "$INSTDIR\TraversePro.exe"
  Delete "$INSTDIR\icon.ico"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  Delete "$SMPROGRAMS\Traverse & Level\Traverse & Level.lnk"
  RMDir "$SMPROGRAMS\Traverse & Level"
  Delete "$DESKTOP\Traverse & Level.lnk"
  DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\TraversePro"
  DeleteRegKey HKCU "Software\TraversePro"
SectionEnd
