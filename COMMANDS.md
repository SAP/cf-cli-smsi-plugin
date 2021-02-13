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
  - checksum: 436433204d70570802033364bfdb926567ef72b9
    platform: linux64
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.4/ServiceManagement.linux64
  - checksum: 6ba3cb8477334254971cdcc857aae617f0eb73f3
    platform: linux32
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.4/ServiceManagement.linux32
  - checksum: 2621df2a2a656742774dcaf462df76342e68c914
    platform: osx
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.4/ServiceManagement.osx
  - checksum: 4e74c2a58193fb301e3179cc8a27ac6fe6f2321f
    platform: win32
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.4/ServiceManagement.win32
  - checksum: 2b71e062e6460a9f6052d36616f4fa9d6f48f695
    platform: win64
    url: https://github.com/SAP/cf-cli-smsi-plugin/releases/download/v1.2.4/ServiceManagement.win64
  company: SAP
  created: 2020-10-07T00:00:00Z
  description: Plugin that makes it easier to use the Service Management service in
    Cloud Foundry.
  homepage: https://github.com/SAP/cf-cli-smsi-plugin
  name: service-management
  updated: 2021-02-13T00:00:00Z
  version: 1.2.4
```
