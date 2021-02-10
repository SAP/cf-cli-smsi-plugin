```
go get -v github.com/buger/jsonparser
go get -v code.cloudfoundry.org/cli/plugin

cf uninstall-plugin ServiceManagement

GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go
chmod 755 ServiceManagement.osx
cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.osx -f
cf plugins | grep ServiceManage


GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go ; chmod 755 ServiceManagement.osx ; cf uninstall-plugin service-management ; cf install-plugin ServiceManagement.osx -f ; cf plugins | grep ServiceManage

```

In BAS
```
cd ~

cf api https://api.cf.us10.hana.ondemand.com ; cf login -u andrew.lunde@sap.com -o ae67provider -s dev

curl -LJO https://github.com/SAP/cf-cli-smsi-plugin/releases/download/latest/ServiceManagement.linux64 ; chmod +x ServiceManagement.linux64 ; f uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.linux64 -f

```
For WIP Drops
```
GOOS=windows GOARCH=amd64 go build -o ServiceManagement.win64 ServiceManagement_plugin.go ; cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.linux64 -f

curl -LJO https://github.com/SAP/cf-cli-smsi-plugin/releases/download/latest/ServiceManagement.win64 ; cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.win64 -f

cf plugins | grep ServiceManage

```

For Release:
The CF cli supports 5 combinations:

linux/386 (known as linux32)
linux/amd64 (known as linux64)
windows/386 (known as win32)
windows/amd64 (known as win64)
darwin /amd64 (known as osx)

```
GOOS=linux GOARCH=amd64 go build -o ServiceManagement.linux64 ServiceManagement_plugin.go
shasum -a 1 ServiceManagement.linux64
GOOS=linux GOARCH=386 go build -o ServiceManagement.linux32 ServiceManagement_plugin.go
shasum -a 1 ServiceManagement.linux32
GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go
shasum -a 1 ServiceManagement.osx
GOOS=windows GOARCH=386 go build -o ServiceManagement.win32 ServiceManagement_plugin.go
shasum -a 1 ServiceManagement.win32
GOOS=windows GOARCH=amd64 go build -o ServiceManagement.win64 ServiceManagement_plugin.go
shasum -a 1 ServiceManagement.win64


```
