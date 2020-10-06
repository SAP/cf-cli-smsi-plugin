```
go get -v github.com/buger/jsonparser
go get -v code.cloudfoundry.org/cli/plugin

cf uninstall-plugin ServiceManagement

GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go
chmod 755 ServiceManagement.osx
cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.osx -f
cf plugins | grep ServiceManage


GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go ; chmod 755 ServiceManagement.osx ; cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.osx -f ; cf plugins | grep ServiceManage

GOOS=linux GOARCH=amd64 go build -o ServiceManagement.linux64 ServiceManagement_plugin.go ; scp ServiceManagement.linux64 thedrop:/home/ec2-user/files

GOOS=windows GOARCH=amd64 go build -o ServiceManagement.win64 ServiceManagement_plugin.go ; scp ServiceManagement.win64 thedrop:/home/ec2-user/files


```

In BAS
```
cd ~

cf api https://api.cf.us10.hana.ondemand.com ; cf login -u andrew.lunde@sap.com -o ae67provider -s dev

curl -LJO https://github.com/andrewlunde/ServiceManagement/releases/download/latest/ServiceManagement.linux64 ; chmod +x ServiceManagement.linux64 ; f uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.linux64 -f

```
For WIP Drops
```
GOOS=linux GOARCH=amd64 go build -o ServiceManagement.linux64 ServiceManagement_plugin.go ; scp ServiceManagement.linux64 thedrop:/home/ec2-user/files

GOOS=windows GOARCH=amd64 go build -o ServiceManagement.win64 ServiceManagement_plugin.go ; cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.linux64 -f

curl -LJO https://github.com/andrewlunde/ServiceManagement/releases/download/latest/ServiceManagement.win64 ; cf uninstall-plugin ServiceManagement ; cf install-plugin ServiceManagement.win64 -f

cf plugins | grep ServiceManage

curl -LJO "Redirects"

curl -o get_smsi http://thedrop.sap-partner-eng.com/files/get_smsi
chmod 755 get_smsi
./get_smsi


curl -LJO http://thedrop.sap-partner-eng.com/files/mod_settings ; chmod 755 mod_settings ; ./mod_settings
```

Andrew Testing
```
cd projects
git clone git@github.com:SAP-samples/cloud-cap-multitenancy.git
git clone https://github.com/SAP-samples/cloud-cap-multitenancy.git
cd ~
ssh-keygen
cat ~/.ssh/id_rsa.pub
<Import into github SSH keys>
cf api https://api.cf.us10.hana.ondemand.com
cf login -u andrew.lunde@sap.com
3<. ae67provider>
cf smsi CAPMT_SMC -o SQLTools > smc.json

jq '.["sqltools.connections"]' smc.json

jq '.["sqltools.connections"] = "[]"' /home/user/.theia/settings.json

vim /home/user/.theia/settings.json smc.json

```

For Release:
The CF cli supports 5 combinations:

linux/386 (known as linux32)
linux/amd64 (known as linux64)
windows/386 (known as win32)
windows/amd64 (known as win64)
darwin /amd64 (known as osx)

```
GOOS=darwin GOARCH=amd64 go build -o ServiceManagement.osx ServiceManagement_plugin.go
GOOS=linux GOARCH=amd64 go build -o ServiceManagement.linux64 ServiceManagement_plugin.go
GOOS=linux GOARCH=386 go build -o ServiceManagement.linux32 ServiceManagement_plugin.go
GOOS=windows GOARCH=amd64 go build -o ServiceManagement.win64 ServiceManagement_plugin.go
GOOS=windows GOARCH=386 go build -o ServiceManagement.win32 ServiceManagement_plugin.go
```