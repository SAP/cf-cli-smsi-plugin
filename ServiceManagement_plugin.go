package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/buger/jsonparser"
)

type ServiceManagementPlugin struct {
	serviceOfferingName *string
	servicePlanName     *string
	showCredentials     *bool
	outputFormat        *string
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type Containers struct {
	ContainerID string
	TenantID    string
}

func (c *ServiceManagementPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	// flags
	flags := flag.NewFlagSet("service-manager-service-instances", flag.ExitOnError)
	serviceOfferingName := flags.String("offering", "hana", "Service offering")
	servicePlanName := flags.String("plan", "hdi-shared", "Service plan")
	showCredentials := flags.Bool("credentials", false, "Show credentials")
	includeMeta := flags.Bool("meta", false, "Include Meta containers")
	includeOwner := flags.Bool("owner", false, "Include Owner credentials")
	outputFormat := flags.String("o", "Txt", "Show as JSON | SQLTools | Txt)")
	modifySettings := flags.Bool("m", false, "Modify settings.json")
	forceUpdates := flags.Bool("f", false, "Force updates (requires -m)")
	offerAll := flags.Bool("a", false, "Offer All Containers option")
	err := flags.Parse(args[1:])
	handleError(err)

	serviceNameGiven := false

	if args[0] == "service-manager-service-instances" {
		/*
			if len(args) < 2 {
				fmt.Println("Please specify an instance of service manager")
				return
			}
		*/

		// https://github.com/cloudfoundry/cli/tree/master/plugin/plugin_examples
		// https://github.com/cloudfoundry/cli/blob/master/plugin/plugin_examples/DOC.md

		// org := plugin_models.Organization{}
		// org, err = cliConnection.GetCurrentOrg()
		// handleError(err)
		// fmt.Println("org = " + org.OrganizationFields.Name)

		serviceManagerName := "Unknown"
		//fmt.Println("args[0] = " + args[0])
		//fmt.Println("args[1] = " + args[1])

		if len(args) > 1 {
			if args[1][0] == '-' {
				//fmt.Println("no sm in args")
				err = flags.Parse(args[1:])
				handleError(err)
			} else {
				serviceNameGiven = true
				serviceManagerName = args[1]
				err = flags.Parse(args[2:])
				handleError(err)
			}
		}

		// return

		if !serviceNameGiven {

			svcs := []plugin_models.GetServices_Model{}

			svcs, err = cliConnection.GetServices()
			handleError(err)

			foundSvcs := []plugin_models.GetServices_Model{}

			for i := 0; i < len(svcs); i++ {
				//fmt.Println("Service Name: " + svcs[i].Name)
				if svcs[i].Service.Name == "service-manager" {
					//fmt.Println("Service Type: " + svcs[i].Service.Name)
					if svcs[i].ServicePlan.Name == "container" {
						//fmt.Println("Service Plan: " + svcs[i].ServicePlan.Name)
						foundSvcs = append(foundSvcs, svcs[i])
					}
				}
			}

			if len(foundSvcs) >= 1 {
				if len(foundSvcs) == 1 {
					serviceManagerName = foundSvcs[0].Name
				} else {
					for i := 0; i < len(foundSvcs); i++ {
						fmt.Println(fmt.Sprintf("%d :", i) + foundSvcs[i].Name)
					}
					fmt.Print("Which service-manager?: ")
					var input string
					fmt.Scanln(&input)
					//fmt.Print(input)
					smidx, _ := strconv.Atoi(input)
					serviceManagerName = foundSvcs[smidx].Name
				}
			} else {
				fmt.Println("Please create at least one instance of service-manager with plan type container.")
				return
			}
		}

		fmt.Println("Using service manager = " + serviceManagerName)

		serviceOfferingName := strings.ToLower(*serviceOfferingName)
		servicePlanName := strings.ToLower(*servicePlanName)

		// validate output format
		outputFormat := strings.ToLower(*outputFormat)
		switch outputFormat {
		case "json", "sqltools", "txt":
		default:
			fmt.Println("Output format must be JSON, SQLTools or Txt")
			return
		}

		// check instance exists
		_, err := cliConnection.GetService(serviceManagerName)
		handleError(err)

		// create service key
		serviceKeyName := "sk-" + args[0]
		_, err = cliConnection.CliCommandWithoutTerminalOutput("create-service-key", serviceManagerName, serviceKeyName)
		handleError(err)

		// get service key
		serviceKey, err := cliConnection.CliCommandWithoutTerminalOutput("service-key", serviceManagerName, serviceKeyName)
		handleError(err)

		// cleanup headers to make parsable as JSON
		serviceKey[0] = ""
		serviceKey[1] = ""

		// authenticate to service manager REST API
		cli := &http.Client{}
		url1, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "url")
		handleError(err)
		req1, err := http.NewRequest("POST", url1+"/oauth/token?grant_type=client_credentials", nil)
		handleError(err)
		req1.Header.Set("Content-Type", "application/json")
		clientid, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "clientid")
		handleError(err)
		clientsecret, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "clientsecret")
		handleError(err)
		req1.SetBasicAuth(clientid, clientsecret)
		res1, err := cli.Do(req1)
		handleError(err)
		defer res1.Body.Close()
		body1Bytes, err := ioutil.ReadAll(res1.Body)
		handleError(err)
		accessToken, err := jsonparser.GetString(body1Bytes, "access_token")

		// get id of service offering
		url2, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "sm_url")
		handleError(err)
		req2, err := http.NewRequest("GET", url2+"/v1/service_offerings", nil)
		handleError(err)
		q2 := req2.URL.Query()
		q2.Add("fieldQuery", "catalog_name eq '"+serviceOfferingName+"'")
		req2.URL.RawQuery = q2.Encode()
		req2.Header.Set("Authorization", "Bearer "+accessToken)
		res2, err := cli.Do(req2)
		handleError(err)
		defer res2.Body.Close()
		body2Bytes, err := ioutil.ReadAll(res2.Body)
		handleError(err)
		numItems, err := jsonparser.GetInt(body2Bytes, "num_items")
		handleError(err)
		if numItems < 1 {
			fmt.Printf("Service offering not found: %s\n", serviceOfferingName)
		} else {
			// get id of service plan for offering
			serviceOfferingId, err := jsonparser.GetString(body2Bytes, "items", "[0]", "id")
			url3, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "sm_url")
			handleError(err)
			req3, err := http.NewRequest("GET", url3+"/v1/service_plans", nil)
			handleError(err)
			q3 := req3.URL.Query()
			q3.Add("fieldQuery", "catalog_name eq '"+servicePlanName+"' and service_offering_id eq '"+serviceOfferingId+"'")
			req3.URL.RawQuery = q3.Encode()
			req3.Header.Set("Authorization", "Bearer "+accessToken)
			res3, err := cli.Do(req3)
			handleError(err)
			defer res3.Body.Close()
			body3Bytes, err := ioutil.ReadAll(res3.Body)
			handleError(err)
			numItems, err = jsonparser.GetInt(body3Bytes, "num_items")
			handleError(err)
			if numItems < 1 {
				fmt.Printf("Service plan not found: %s\n", servicePlanName)
			} else {
				servicePlanId, err := jsonparser.GetString(body3Bytes, "items", "[0]", "id")

				// get service instances for service plan
				url4, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "sm_url")
				handleError(err)
				req4, err := http.NewRequest("GET", url4+"/v1/service_instances", nil)
				handleError(err)
				q4 := req4.URL.Query()
				q4.Add("fieldQuery", "service_plan_id eq '"+servicePlanId+"'")
				req4.URL.RawQuery = q4.Encode()
				req4.Header.Set("Authorization", "Bearer "+accessToken)
				res4, err := cli.Do(req4)
				handleError(err)
				defer res4.Body.Close()
				body4Bytes, err := ioutil.ReadAll(res4.Body)
				handleError(err)
				numItems, err = jsonparser.GetInt(body4Bytes, "num_items")
				handleError(err)

				foundContainers := []Containers{}
				var addConn = `{`

				// for each item
				var item = 0
				var isMeta = false
				jsonparser.ArrayEach(body4Bytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					isMeta = false
					id, _ := jsonparser.GetString(value, "id")

					// get service binding
					url5, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "sm_url")
					handleError(err)
					req5, err := http.NewRequest("GET", url5+"/v1/service_bindings", nil)
					handleError(err)
					q5 := req5.URL.Query()
					q5.Add("fieldQuery", "service_instance_id eq '"+id+"'")
					req5.URL.RawQuery = q5.Encode()
					req5.Header.Set("Authorization", "Bearer "+accessToken)
					res5, err := cli.Do(req5)
					handleError(err)
					defer res5.Body.Close()
					body5Bytes, err := ioutil.ReadAll(res5.Body)
					handleError(err)

					tenantID, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "labels", "tenant_id", "[0]")

					//spaceName, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "context", "space_name")

					var splits = strings.Split(tenantID, "-")
					if splits[0] == "TENANT" {
						isMeta = true
					}

					if !isMeta || (isMeta && *includeMeta) {
						//fmt.Printf("%d: %s \n", item, tenantID)
						container := Containers{ContainerID: id, TenantID: tenantID}
						foundContainers = append(foundContainers, container)
						item = item + 1
					}
				}, "items")

				whichID := "ALL"

				if len(foundContainers) > 1 {
					if *offerAll {
						fmt.Printf("%d: %s \n", 0, "Include All")
					}
					for i := 0; i < len(foundContainers); i++ {
						fmt.Printf("%d. %s \n", i+1, foundContainers[i].TenantID)
					}

					fmt.Print("Container> ")
					var input string
					fmt.Scanln(&input)
					cidx, _ := strconv.Atoi(input)
					if cidx == 0 && *offerAll {
						fmt.Printf("Using: %s \n", "All Containers")
					} else {
						whichContainer := foundContainers[cidx-1].TenantID
						fmt.Printf("Using: %s \n", whichContainer)
						whichID = foundContainers[cidx-1].ContainerID
						item = 1
					}
				} else {
					whichID = foundContainers[0].ContainerID
					item = 1
				}

				switch outputFormat {
				case "json":
					fmt.Printf("{\n\"service_offering\": \"%s\", \n\"service_plan\": \"%s\", \n\"num_items\": %d, \n\"items\": \n [\n", serviceOfferingName, servicePlanName, item)
				case "sqltools":
					fmt.Printf(`{"sqltools.connections": [`)
				case "txt":
					if !*modifySettings {
						fmt.Printf("%d items found for service offering %s and service plan %s.\n", numItems, serviceOfferingName, servicePlanName)
					}
				}

				// for each item
				item = 0
				isMeta = false
				jsonparser.ArrayEach(body4Bytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					//item = item + 1
					isMeta = false
					id, _ := jsonparser.GetString(value, "id")

					name, _ := jsonparser.GetString(value, "name")

					createdAt, _ := jsonparser.GetString(value, "created_at")
					updatedAt, _ := jsonparser.GetString(value, "updated_at")
					ready, _ := jsonparser.GetBoolean(value, "ready")
					usable, _ := jsonparser.GetBoolean(value, "usable")

					if (whichID == id) || (whichID == "ALL") {

						// get service binding
						url5, err := jsonparser.GetString([]byte(strings.Join(serviceKey, "")), "sm_url")
						handleError(err)
						req5, err := http.NewRequest("GET", url5+"/v1/service_bindings", nil)
						handleError(err)
						q5 := req5.URL.Query()
						q5.Add("fieldQuery", "service_instance_id eq '"+id+"'")
						req5.URL.RawQuery = q5.Encode()
						req5.Header.Set("Authorization", "Bearer "+accessToken)
						res5, err := cli.Do(req5)
						handleError(err)
						defer res5.Body.Close()
						body5Bytes, err := ioutil.ReadAll(res5.Body)
						handleError(err)
						host, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "host")

						port, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "port")

						driver, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "driver")
						schema, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "schema")

						certificate, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "certificate")
						re := regexp.MustCompile(`\n`)
						certificate = re.ReplaceAllString(certificate, "")
						url, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "url")
						user, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "user")
						password, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "password")
						var hdiuser = ""
						var hdipassword = ""
						if servicePlanName == "hdi-shared" {
							hdiuser, _ = jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "hdi_user")
							hdipassword, _ = jsonparser.GetString(body5Bytes, "items", "[0]", "credentials", "hdi_password")
						}

						tenantID, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "labels", "tenant_id", "[0]")

						spaceName, _ := jsonparser.GetString(body5Bytes, "items", "[0]", "context", "space_name")

						var splits = strings.Split(tenantID, "-")
						if splits[0] == "TENANT" {
							isMeta = true
						}

						//name = serviceManagerName + ":" + tenantID
						//name = tenantID

						// Need to use the SAPCP API to get the subdomain from the subaccount GUID which is the tenantID
						// sapcp get accounts/subaccount b44f32d4-6e31-4d95-b17f-6c6fcdb37e1f

						if !isMeta || (isMeta && *includeMeta) {
							item = item + 1
							if outputFormat == "json" {
								if item > 1 {
									fmt.Printf(",\n")
								}
								fmt.Printf("  {\n  \"name\": \"%s\", \n  \"id\": \"%s\", \n  \"tenant\": \"%s\", \n  \"created_at\": \"%s\", \n  \"updated_at\": \"%s\", \n  \"ready\": %t, \n  \"usable\": %t, \n  \"schema\": \"%s\", \n  \"host\": \"%s\", \n  \"port\": \"%s\", \n  \"url\": \"%s\", \n  \"driver\": \"%s\"", name, id, tenantID, createdAt, updatedAt, ready, usable, schema, host, port, url, driver)
								if *showCredentials {
									fmt.Printf(", \n  \"user\": \"%s\", \n  \"password\": \"%s\", \n  \"certificate\": \"%s\"", user, password, certificate)
									if servicePlanName == "hdi-shared" && *includeOwner {
										fmt.Printf(",\n  \"hdi_user\": \"%s\", \n  \"hdi_password\": \"%s\"", hdiuser, hdipassword)
									}
								}
								fmt.Printf("\n  }")
							} else if outputFormat == "sqltools" {
								if item > 1 {
									fmt.Printf(`,`)
								}
								fmt.Printf(`{"name": "%s", "group": "SMSI", "dialect": "SAPHana", "driver": "SAPHana", "server": "%s", "port": %s, "database": "%s", "username": "%s", "password": "%s", "connectionTimeout": 30, "hanaOptions": {"encrypt": true, "sslValidateCertificate": true, "sslCryptoProvider": "openssl", "sslTrustStore": "%s"}}`, serviceManagerName+":"+tenantID, host, port, schema, user, password, certificate)
								if servicePlanName == "hdi-shared" && *includeOwner {
									fmt.Printf(`,{"name": "%s", "group": "SMSI", "dialect": "SAPHana", "driver": "SAPHana", "server": "%s", "port": %s, "database": "%s", "username": "%s", "password": "%s", "connectionTimeout": 30, "hanaOptions": {"encrypt": true, "sslValidateCertificate": true, "sslCryptoProvider": "openssl", "sslTrustStore": "%s"}}`, serviceManagerName+":"+spaceName+":"+tenantID+":OWNER", host, port, schema, hdiuser, hdipassword, certificate)
								}
							} else {
								//txt
								fmt.Printf("\nName: %s \nId: %s \nTenant: %s \nCreatedAt: %s \nUpdatedAt: %s \nReady: %t \nUsable: %t \nSchema: %s \nHost: %s \nPort: %s \nURL: %s \nDriver: %s\n", name, id, tenantID, createdAt, updatedAt, ready, usable, schema, host, port, url, driver)
								if *showCredentials {
									fmt.Printf("User: %s \nPassword: %s \nCertificate: %s \n", user, password, certificate)
									if servicePlanName == "hdi-shared" && *includeOwner {
										fmt.Printf("HDIUser: %s \nHDIPassword: %s \n", hdiuser, hdipassword)
									}
								}
								// fmt.Printf("TenantID: %s \n", tenantID)
							}
							if item > 1 {
								addConn += `},{`
							}
							// Put all the addConn stuff here
							addConn += "\n" + `"name": "` + serviceManagerName + `:` + spaceName + `:` + tenantID + `",` + "\n"
							addConn += `"group": "` + serviceManagerName + `:` + spaceName + `",` + "\n"
							addConn += `"driver": "` + `SAPHana` + `",` + "\n"
							addConn += `"dialect": "` + `SAPHana` + `",` + "\n"

							addConn += `"server": "` + host + `",` + "\n"
							addConn += `"port": ` + port + `,` + "\n"

							addConn += `"database": "` + schema + `",` + "\n"

							if *includeOwner {
								addConn += `"username": "` + hdiuser + `",` + "\n"
								addConn += `"password": "` + hdipassword + `",` + "\n"
							} else {
								addConn += `"username": "` + user + `",` + "\n"
								addConn += `"password": "` + password + `",` + "\n"
							}

							addConn += `"previewLimit": ` + `50` + `,` + "\n"
							addConn += `"connectionTimeout": ` + `30` + `,` + "\n"
							addConn += `"hanaOptions": ` + `{` + `` + "\n"
							addConn += `     "encrypt": ` + `true` + `,` + "\n"
							addConn += `     "sslValidateCertificate": ` + `true` + `,` + "\n"
							addConn += `     "sslCryptoProvider": ` + `"openssl"` + `,` + "\n"
							addConn += `     "sslTrustStore": "` + certificate + `"` + "\n"

							addConn += `` + `}` + "\n"

						}
					}
				}, "items")

				switch outputFormat {
				case "json":
					fmt.Println("\n ]\n}\n")
				case "sqltools":
					fmt.Println(`]}`)
				}

				addConn += `}`

				// =====================================================================
				// =====================================================================
				// =====================================================================
				// modifySettings = mod_settings.go
				// =====================================================================
				// =====================================================================
				// =====================================================================

				if *modifySettings {
					fmt.Println("")
					fmt.Println("modifySettings: " + "true")
					if *forceUpdates {
						fmt.Println("forceUpdates: " + "true")
					} else {
						fmt.Println("forceUpdates: " + "false")
					}

					fmt.Println("")

					// fmt.Println("addConn: " + addConn)

					//fmt.Println(runtime.GOOS)
					//fmt.Println(runtime.GOARCH)

					user, err := user.Current()
					if err != nil {
						log.Fatalf(err.Error())
					}
					homeDirectory := user.HomeDir
					fmt.Printf("Home Directory: %s\n", homeDirectory)

					var inSettings = false
					var isBAS = false
					// Scan for *.theia-workspace files in BAS ??
					var defaultsFile = "Unknown"
					//var defaultsExists = false

					var settingsFile = "Unknown"
					var settingsExists = false

					var skipping = false

					switch runtime.GOOS {
					case "darwin":
						fmt.Println("On Mac:")
						//
						// The current code-workspace file can be found by looking here.
						// cat $HOME/Library/Application\ Support/Code/storage.json | grep -A 3 lastActiveWindow
						// User (GLOBAL) settings file
						// If this(User) file has a sqltools.connections object but the current code-workspace files doesn't
						// then this is used.  Otherwise it's ignored as soon as the code-workspace settings->sqltools.connections exists
						// Currently SQLTools won't allow writing settings into the User file, but will display them if they already exist.
						settingsFile = homeDirectory + "/Library/Application Support/Code/User/settings.json"

						defaultsFile = homeDirectory + "/Library/Application Support/Code/storage.json"
						byteValue, err := ioutil.ReadFile(defaultsFile)
						if err == nil {
							configURIPath, err := jsonparser.GetString(byteValue, "windowsState", "lastActiveWindow", "workspaceIdentifier", "configURIPath")
							if err == nil {
								fmt.Println("configURIPath: " + configURIPath)
								settingsFile = "/" + strings.TrimLeft(configURIPath, "file:/")
								inSettings = true // File has sqltools.connections at the top-level
							}
						}

					case "linux":
						fmt.Println("On Linux:")

						// Check to see if BAS
						settingsFile = homeDirectory + "/.theia/settings.json"
						if _, err := os.Stat(settingsFile); err == nil {
							// path/to/whatever exists
							fmt.Println("We are in BAS since " + settingsFile + " Exists!")
							inSettings = false
							isBAS = true
						}

						if inSettings {
							settingsFile = "~/Code/User/"
						} else { //User(Global) Settings
							if !isBAS {
								settingsFile = homeDirectory + "/.config/Code/User/settings.json"
							}
						}

					case "windows":
						fmt.Println("On Windoz:")

						appData := os.Getenv("APPDATA")
						fmt.Printf("appData: %s\n", appData)

						//APPDATA=C:\Users\I830671\AppData\Roaming

						defaultsFile = appData + "/Code/storage.json"
						fmt.Println("defaultsFile: " + defaultsFile)

						byteValue, err := ioutil.ReadFile(defaultsFile)
						if err == nil {
							configURIPath, err := jsonparser.GetString(byteValue, "windowsState", "lastActiveWindow", "workspaceIdentifier", "configURIPath")
							if err == nil {
								fmt.Println("configURIPath: " + configURIPath)
								settingsFile = strings.TrimLeft(configURIPath, "file:/")
								settingsFile = strings.Replace(settingsFile, "%3A", ":", -1)
								//fmt.Println("settingsFile: " + settingsFile)
								inSettings = true // File has sqltools.connections at the top-level
							}
						}
					}

					fmt.Println("settingsFile: " + settingsFile)
					if inSettings {
						fmt.Println("Look in settings...")
					} else {
						fmt.Println("Look at top-level..")
					}

					if _, err := os.Stat(settingsFile); err == nil {
						// path/to/whatever exists
						fmt.Println("settingsFile: " + settingsFile + " Exists!")
						settingsExists = true

					} else if os.IsNotExist(err) {
						// path/to/whatever does *not* exist
						fmt.Println("settingsFile: " + settingsFile + " Does NOT Exist!")

					} else {
						// Schrodinger: file may or may not exist. See err for details.

						// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
						fmt.Println("settingsFile: " + settingsFile + " Existence Unknown!")
						settingsExists = true
					}

					fmt.Println("")

					// var newConn = `{
					// 	"name": "CAPMT_SMC:subAcct",
					// 	"group": "SMSI",
					// 	"driver": "SAPHana",
					// 	"dialect": "SAPHana",
					// 	"server": "833726c5-cca3-4dce-a325-4385426009e7.hana.trial-us10.hanacloud.ondemand.com",
					// 	"port": 443,
					// 	"database": "D53EE042B6AD4E8093FF0A24F931586B",
					// 	"username": "D53EE042B6AD4E8093FF0A24F931586B_B5IBO9PWMQ841D52POXNE26XN_RT",
					// 	"password": "Mw9h7H.5r6CBidD2vtq.vxmzisxLAMx2_UJ9YrjZim2Yop-kUOcBII-g6VHYZMDpPzjT0PCQua.8i-V2f8MrjDqkGG6hRZAct2a2YIL7PFrlzeSDhO5qBOl6ni-VRF3t",
					// 	"connectionTimeout": 30,
					// 	"hanaOptions": {
					// 		"encrypt": true,
					// 		"sslValidateCertificate": true,
					// 		"sslCryptoProvider": "openssl",
					// 		"sslTrustStore": "-----BEGIN CERTIFICATE-----MIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0wNjExMTAwMDAwMDBaFw0zMTExMTAwMDAwMDBaMGExCzAJBgNVBAYTAlVTMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5jb20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jvhEXLeqKTTo1eqUKKPC3eQyaKl7hLOllsBCSDMAZOnTjC3U/dDxGkAV53ijSLdhwZAAIEJzs4bg7/fzTtxRuLWZscFs3YnFo97nh6Vfe63SKMI2tavegw5BmV/Sl0fvBf4q77uKNd0f3p4mVmFaG5cIzJLv07A6Fpt43C/dxC//AH2hdmoRBBYMql1GNXRor5H4idq9Joz+EkIYIvUX7Q6hL+hqkpMfT7PT19sdl6gSzeRntwi5m3OFBqOasv+zbMUZBfHWymeMr/y7vrTC0LUq7dBMtoM1O/4gdW7jVg/tRvoSSiicNoxBN33shbyTApOB6jtSj1etX+jkMOvJwIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUA95QNVbRTLtm8KPiGxvDl7I90VUwHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUwDQYJKoZIhvcNAQEFBQADggEBAMucN6pIExIK+t1EnE9SsPTfrgT1eXkIoyQY/EsrhMAtudXH/vTBH1jLuG2cenTnmCmrEbXjcKChzUyImZOMkXDiqw8cvpOp/2PV5Adg06O/nVsJ8dWO41P0jmP6P6fbtGbfYmbW0W5BjfIttep3Sp+dWOIrWcBAI+0tKIJFPnlUkiaY4IBIqDfv8NZ5YBberOgOzW6sRBc4L0na4UU+Krk2U886UAb3LujEV0lsYSEY1QSteDwsOoBrp+uvFRTp2InBuThs4pFsiv9kuXclVzDAGySj4dzp30d8tbQkCAUw7C29C79Fv1C5qfPrmAESrciIxpg0X40KPMbp1ZWVbd4=-----END CERTIFICATE-----"
					// 		}
					// 	}`

					//connName, _ := jsonparser.GetString([]byte(newConn), "name")
					connName, _ := jsonparser.GetString([]byte(addConn), "name")

					var foundIdx int = -1

					if settingsExists {
						// read file
						byteValue, err := ioutil.ReadFile(settingsFile)
						if err != nil {
							fmt.Print(err)
						} else {
							//err := jsonparser.GetString(data, "items", "[0]", "id")
							//colorTheme, err := jsonparser.GetString(byteValue, "workbench.colorTheme")
							//handleError(err)
							//fmt.Println("colorTheme: " + colorTheme)

							// var newValue []byte
							// var newType jsonparser.ValueType
							// var newOffset int = 0

							var dataValue []byte
							var dataType jsonparser.ValueType
							var dataOffset int = 0
							var needsSettings = false

							if inSettings {
								dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "settings")
								if err != nil {
									fmt.Println("settings" + " Key path not found")
									needsSettings = true
								} else {
									dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "settings", "sqltools.connections")
									if err != nil {
										fmt.Println("sqltools.connections" + " Key path not found")
										// We can go ahead and add it.
									}
								}
							} else {
								dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "sqltools.connections")
							}

							if err != nil {
								fmt.Println("sqltools.connections" + " Key path not found")
								// We can go ahead and add it.
							}

							// fmt.Println("dataValue: " + string(dataValue))
							// fmt.Println("offset: ", dataOffset)

							if dataType == jsonparser.NotExist {
								fmt.Println("offset: ", dataOffset)
								fmt.Println("sqltools.connections" + " is NotExist")
								// IF this is the case then we can safely create a new sqltools.connections array and append it to settings

								var newSQLToolsConn string
								newSQLToolsConn = string(byteValue)
								newSQLToolsConn = strings.TrimSpace(newSQLToolsConn)
								newSQLToolsConn = strings.TrimRight(newSQLToolsConn, "}")
								newSQLToolsConn = strings.TrimSpace(newSQLToolsConn)
								newSQLToolsConn += ",\n"
								if needsSettings {
									fmt.Println("adding settings {} ")
									newSQLToolsConn += "\"settings\": \n { "
								}
								fmt.Println("adding sqltools.connections [] ")
								newSQLToolsConn += "\"sqltools.connections\": [ \n"
								//newSQLToolsConn += newConn + "] }"
								newSQLToolsConn += addConn
								newSQLToolsConn += "\n]\n"
								if needsSettings {
									newSQLToolsConn += "} \n"
								}
								newSQLToolsConn += "} \n"

								// write file
								err = ioutil.WriteFile(settingsFile, []byte(newSQLToolsConn), 0644)
								handleError(err)

							} else if dataType == jsonparser.Array {
								// fmt.Println("sqltools.connections" + " is an Array")

								var scidx int = 0
								jsonparser.ArrayEach(dataValue, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
									name, _ := jsonparser.GetString(value, "name")
									// fmt.Println("name: " + name)
									if connName != name {
										fmt.Println("keeping: " + name)
									} else {
										if *modifySettings {
											if *forceUpdates {
												fmt.Println("replacing: " + name)
											} else {
												fmt.Println("duplicate: " + name)
											}
										} else {
											fmt.Println("skipping: " + name)
										}
										foundIdx = scidx
										skipping = true
									}
									scidx = scidx + 1
								})
								// https://github.com/buger/jsonparser#set

								if !skipping {
									fmt.Println("adding:  " + connName + "")

									var newSQLToolsConn string

									newSQLToolsConn = string(dataValue)
									newSQLToolsConn2 := strings.TrimRight(newSQLToolsConn, "]")
									newSQLToolsConn = newSQLToolsConn2
									if scidx > 0 {
										newSQLToolsConn += ","
									}
									//newSQLToolsConn += newConn + "]"
									newSQLToolsConn += addConn + "]"

									var setValue []byte

									// fmt.Println("attempt set: ")

									if inSettings {
										setValue, err = jsonparser.Set(byteValue, []byte(newSQLToolsConn), "settings", "sqltools.connections")
									} else {
										setValue, err = jsonparser.Set(byteValue, []byte(newSQLToolsConn), "sqltools.connections")
									}
									handleError(err)

									//fmt.Println("after set: ")
									// jsonparser.ArrayEach(setValue, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
									// 	name, _ := jsonparser.GetString(value, "name")
									// 	fmt.Println("name: " + name)
									// })

									// fmt.Println("newConn: " + newConn)
									// fmt.Println("setValue: " + string(setValue))

									// write file
									err = ioutil.WriteFile(settingsFile, setValue, 0644)
									handleError(err)
								} else {
									if *modifySettings {
										if !*forceUpdates {
											fmt.Println("Connection with name " + connName + " already exists!  -f to force replacement.")
										}
										idxStr := "[" + strconv.Itoa(foundIdx) + "]"
										// idxStr := strconv.Itoa(foundIdx)
										// fmt.Println("idxStr:" + idxStr)
										var setValue []byte
										if inSettings {
											// dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "settings", "sqltools.connections", idxStr)
											// setValue, err = jsonparser.Set(byteValue, []byte(newConn), "settings", "sqltools.connections", idxStr)
											if *modifySettings && *forceUpdates {
												setValue, err = jsonparser.Set(byteValue, []byte(addConn), "settings", "sqltools.connections", idxStr)
											}
										} else {
											// dataValue, dataType, dataOffset, err = jsonparser.Get(byteValue, "sqltools.connections", idxStr)
											// setValue, err = jsonparser.Set(byteValue, []byte(newConn), "settings", "sqltools.connections", idxStr)
											if *modifySettings && *forceUpdates {
												setValue, err = jsonparser.Set(byteValue, []byte(addConn), "sqltools.connections", idxStr)
											}
										}
										handleError(err)

										//fmt.Println("setValue: " + string(setValue))
										//fmt.Println("offset: ", dataOffset)

										// fmt.Println("setValue: " + string(setValue))

										if *modifySettings && *forceUpdates {
											// write file
											err = ioutil.WriteFile(settingsFile, setValue, 0644)
											handleError(err)
										}

									} else {
										fmt.Println("Connection with name " + connName + " already exists!  Delete it first and rerun.")
									}
								}

							} else if dataType == jsonparser.Object {
								fmt.Println("sqltools.connections" + " is Object")

							} else if dataType == jsonparser.Null {
								fmt.Println("sqltools.connections" + " is Null")

							} else {
								fmt.Println("sqltools.connections" + " is unexpected")

							}

						}
					}

				} else {
					fmt.Println("modifySettings: " + "false")
				}
			}

		}

		// delete service key
		_, err = cliConnection.CliCommandWithoutTerminalOutput("delete-service-key", serviceManagerName, serviceKeyName, "-f")
		handleError(err)

		fmt.Println("")

	}
}

func (c *ServiceManagementPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "ServiceManagement",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 12,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "service-manager-service-instances",
				Alias:    "smsi",
				HelpText: "Show service manager service instances for a service offering and plan.",
				UsageDetails: plugin.Usage{
					Usage: "cf service-manager-service-instances [SERVICE_MANAGER_INSTANCE] [--offering <SERVICE_OFFERING>] [--plan <SERVICE_PLAN>] [--credentials] [--meta] [--owner] [-o JSON | SQLTools | Txt] [-m [-f]] [-a]",
					Options: map[string]string{
						"credentials": "Show credentials",
						"meta":        "Include Meta containers",
						"owner":       "Include Owner credentials",
						"o":           "Show as JSON | SQLTools | Txt (default 'Txt')",
						"offering":    "Service offering (default 'hana')",
						"plan":        "Service plan (default 'hdi-shared')",
						"m":           "Modify settings.json",
						"f":           "Force updates (requires -m)",
						"a":           "Offer All Containers option",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(ServiceManagementPlugin))
}
