@ECHO OFF
setx I2CP_HOME "%ProgramFiles%\RTradeLtd\go-anonvpn"
setx GO_I2CP_CONF \i2cp.conf
md "%AppData%\RTradeLtd\"
cd "%AppData%\RTradeLtd\"
if exist "%ProgramFiles%\RTradeLtd\go-anonvpn\zero\router\i2p-zero.exe" (
    start /b "%ProgramFiles%\RTradeLtd\go-anonvpn\zero\router\i2p-zero.exe"
    start /b "%ProgramFiles%\RTradeLtd\go-anonvpn\tun2socks.bat"
)
"%ProgramFiles%\RTradeLtd\go-anonvpn\anonvpn.exe" -file "%ProgramFiles%\RTradeLtd\go-anonvpn\rtrade-testserver.ini" -directory "%AppData%\RTradeLtd\"
pause
