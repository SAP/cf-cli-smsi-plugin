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

```
/Users/i830671/git/cli-plugin-repo/repo-index.yml
- authors:
  - contact: andrew.lunde@sap.com
    homepage: https://github.com/andrewlunde
    name: Andrew Lunde
  binaries:
  - checksum: 207b1b0f1275894f65ef78040086c23558973989
    platform: linux64
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.2/ServiceManagement.linux64
  - checksum: 394c83276042d586683fdbea767610bc5742d77b
    platform: linux32
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.2/ServiceManagement.linux32
  - checksum: b8901228079598b29592948502089b95767744e5
    platform: osx
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.2/ServiceManagement.osx
  - checksum: 506ef402b807e077c87fdae7204be912692e3ac9
    platform: win32
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.2/ServiceManagement.win32
  - checksum: 7e5b7e00d1b4f6064803417f21554423d9c06edb
    platform: win64
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.2/ServiceManagement.win64
  company: SAP
  created: 2020-10-07T00:00:00Z
  description: Plugin that makes it easier to use the Service Management service in
    Cloud Foundry.
  homepage: https://github.com/SAP/cf-cli-smsi-plugin
  name: service-management
  updated: 2021-02-12T00:00:00Z
  version: 1.2.2
```
