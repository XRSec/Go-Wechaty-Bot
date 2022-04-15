@chcp 65001
@cls
@set WECHATY_PUPPET=wechaty-puppet-xp
@set WECHATY_TOKEN=5f3029c0-0f46-4436-bdc6-02efcbad3309
@REM export WECHATY_TOKEN=$(curl -s https://www.uuidgenerator.net/api/version4)

@set WECHATY_PUPPET_SERVICE_TOKEN=insecure_34bf8353-0874-4b29-851d-e8a2502fc747
@REM set WECHATY_PUPPET_SERVICE_TOKEN="insecure_$(curl -s https://www.uuidgenerator.net/api/version4)"

@set WECHATY_PUPPET_SERVER_PORT=25000
@set WECHATY_LOG="verbose"
@set WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER=true

@echo WECHATY_TOKEN = %WECHATY_TOKEN%
@echo WECHATY_PUPPET_SERVICE_TOKEN = %WECHATY_PUPPET_SERVICE_TOKEN%

wechaty gateway --puppet %WECHATY_PUPPET% --port %WECHATY_PUPPET_SERVER_PORT% --token %WECHATY_PUPPET_SERVICE_TOKEN% --puppet-token %WECHATY_TOKEN%