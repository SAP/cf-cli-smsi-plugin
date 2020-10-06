<!--
SPDX-FileCopyrightText: 2020 Andrew Lunde <andrew.lunde@sap.com>

SPDX-License-Identifier: Apache-2.0
-->
[![REUSE status](https://api.reuse.software/badge/github.com/SAP-samples/cloud-sfsf-benefits-ext)](https://api.reuse.software/info/github.com/SAP-samples/cloud-sfsf-benefits-ext)

# cf-cli-smsi-plugin CF CLI Service Management Service Instance Plugin

This is a Cloud Foundry CLI plugin designed to make it easier when using the [Service Management](https://help.sap.com/viewer/product/SERVICEMANAGEMENT/Cloud/en-US) service in Cloud Foundry. It uses the service manager REST API to get details of service instances that the service management service has created.

# Requirements
Installed CloudFoundry CLI - ensure that CloudFoundry CLI is installed and working. For more information about installation of CloudFoundry CLI, please visit the official CloudFoundry [documentation](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).

If you are running from within VSCode, you need to create a Workspace before using the -m (Modify settings.json) option.

# Download and Installation

Check whether you have a previous version installed, using the command: `cf plugins`. If the ServiceManagement plugin is already installed, you need to uninstall it first and then to install the new version. You can uninstall the plugin using command `cf uninstall-plugin ServiceManagement`.

## CF Community Plugin Repository

Hopefully the ServiceManagement CF CLI Plugin will soon be available on the CF Community Repository. To install the latest available version of the ServiceManagement CLI Plugin execute the following:

`cf install-plugin ServiceManagement`

If you do not have the community repository in your CF CLI you can add it first by executing:

`cf add-plugin-repo CF-Community https://plugins.cloudfoundry.org`

## Manual Installation

Alternatively you can install any version of the plugin by manually downloading it from the releases page and installing the binaries for your specific operating system.

### Download
The latest version of the plugin can also be downloaded from the project's [releases](https://github.com/andrewlunde/ServiceManagement/releases/latest). Download the plugin for your platform (Darwin, Linux, Windows). The name for the correct plugin for each platform can be found in the table below.

Mac OS X 64 bit | Windows 64 bit | Windows 32 bit | Linux 64 bit | Linux 32 bit
--- | --- | --- | --- | ---
ServiceManagement.osx | ServiceManagement.win64 | ServiceManagement.win32 | ServiceManagement.linux64 | ServiceManagement.linux32

### Installation
Install the plugin, using the following command:
```
cf install-plugin <path-to-the-plugin> -f
```
Note: if you are running on a Unix-based system, you need to make the plugin executable before installing it. In order to achieve this, execute the following commad `chmod +x <path-to-the-plugin>`

## Usage
The ServiceManagement CF plugin supports the following commands:

Command Name | Command Description
--- | ---
`service-manager-service-instances` `smsi` | Show service manager service instances for a service offering and plan. The default service offering is `hana` and the default service plan is `hdi-shared`. Defaults can be overriden with the `-offering`and/or `-plan`flags. Use the `-credentials` flag to show login credentials and certificates. Use the `-o` flag to show results in JSON, SQLTools or Txt format. It's possible to pipe the information to a file as shown in the examples below. When using `-o SQLTools` the `-credentials` flag will be set automatically. If `-m` is specified, an attempt to find the appropriate settings.json file will be made and ,if found, modify it.  If `-f` is specified in addition to `-m`, connections that already exist in the settings.json file will be forceably updated.

Usage:

```cf service-manager-service-instances [SERVICE_MANAGER_INSTANCE] [-offering <SERVICE_OFFERING>] [-plan <SERVICE_PLAN>] [--credentials] [--meta] [--owner] [-o JSON | SQLTools | Txt] [-m [-f]]```

Examples:

```cf service-manager-service-instances my-sm```

```cf service-manager-service-instances my-sm -o SQLTools```

```cf service-manager-service-instances my-sm -credentials```

```cf smsi```

```cf smsi my-sm -credentials -o JSON```

```cf smsi my-sm -credentials -o JSON > my-results.json```

```cf smsi my-sm -o SQLTools > settings.json```

```cf smsi my-sm -credentials > my-results.txt```

```cf smsi my-sm -credentials -o SQLTools -offering hanatrial -plan schema```

For more information, see the command help output available via `cf [command] --help` or `cf help [command]`.

# Further Information
Tutorials:
- [SAP Cloud Platform Developer Onboarding](https://www.youtube.com/playlist?list=PLkzo92owKnVw3l4fqcLoQalyFi9K4-UdY)
- [SAP HANA Academy](https://www.youtube.com/saphanaacademy)

# Changes

Refer to the [CHANGELOG](CHANGELOG.md) for revision details.

# Support

If you encounter an issue, you can [create a ticket](issues/new/choose).

# Contributing

## Fork & create a branch
If this is something you think you can fix, then fork this repo and create a branch with a descriptive name.

## Make a Pull Request
Go to GitHub and make a [Pull Request](pulls)

In the subject of the pull request, use *feat:* to denote an enhancement, **fix:** to denote a bug fix, ***chore:*** for small configuration updates or ***docs:*** for documentation updates and briefly describe the bug fix or enhancement you are contributing.

# Versioning
The Service Management plugin follows Semantic Versioning. These components strictly adhere to the [MAJOR].[MINOR].[PATCH] numbering system (also known as [BREAKING].[FEATURE].[FIX]).

# License

This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the [LICENSE](LICENSES/Apache-2.0.txt) file.
