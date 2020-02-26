
!include "MUI.nsh"
!include "LogicLib.nsh"
!include "x64.nsh"

!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.

!insertmacro MUI_PAGE_WELCOME # Welcome to the installer page.
!insertmacro MUI_PAGE_LICENSE "../LICENSES.txt"
!insertmacro MUI_PAGE_DIRECTORY # In which folder install page.
!insertmacro MUI_PAGE_INSTFILES # Installing page.
!insertmacro MUI_PAGE_FINISH # Finished installation page.

!insertmacro MUI_LANGUAGE "English"

Name "go-anonvpn"
OutFile "go-anonvpn-installer.exe"
InstallDir "$PROGRAMFILES64\RTradeLtd\go-anonvpn"
ShowInstDetails show
RequestExecutionLevel admin

!include "geti2p.nsi"

Section "AnonVPN"
  SetOutPath $INSTDIR

  ${If} ${RunningX64}
    File ../cmd/anonvpn/anonvpn.exe
  ${Else}
    File /oname=anonvpn.exe ../cmd/anonvpn/anonvpn-32.exe
  ${EndIf}

  File ../etc/anonvpn/anonvpn.ini
  File ../etc/anonvpn/reseed.zip
  File ../etc/anonvpn/i2cp.conf
  #File ../rtrade-testserver.ini
  #File ../davpn.ico
  File ../cmd/anonvpn/launch.bat
  File ../anonvpn.exe.manifest

  ${If} ${RunningX64}
    File wintun.msi
  ${Else}
    File /oname=wintun.msi wintun32.msi
  ${EndIf}

  File wintun.msm
  File wintap.exe
  ExecWait "msiexec.exe /i wintun.msi"
  ExecWait "wintap.exe"
  createDirectory "$SMPROGRAMS\RTradeLtd"
  createDirectory "$SMPROGRAMS\RTradeLtd\i2p-zero(installed by davpn)"
  CreateShortCut "$SMPROGRAMS\RTradeLtd\davpn.lnk" "C:\Windows\system32\cmd.exe" "/c $\"$INSTDIR\launch.bat$\"" "$INSTDIR\davpn.ico"
  CreateShortCut "$SMPROGRAMS\RTradeLtd\i2p-zero(installed by davpn)\i2p-zero.lnk" "$INSTDIR\zero\router\i2p-zero.exe"
SectionEnd


