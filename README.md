# GetUsersIPFromEvtx
通过域控日志定位个人PC的IP地址


用户在域内个人PC使用域账号登陆成功后，会在域控日志留下4624、4768、4769三条日志记录，4624对应登录成功，4768、4769对应Kerberos身份认证。


.\main.exe C:\Windows\System32\winevt\Logs\Security.evtx
2021/10/28 13:12:46 [+] Already start...
2021/10/28 13:12:47 [+] The result has been saved in the current folder...
