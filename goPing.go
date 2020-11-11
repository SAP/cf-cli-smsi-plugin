// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"

	"github.com/SAP/go-hdb/driver"

	// Register hdb driver.
	_ "github.com/SAP/go-hdb/driver"
)

const (
	driverName = "hdb"
	hdbDsn     = "hdb://2E403288791D4A56829B833EFA7A1529_3RGDF00CGM08B711H2QKDPA3F_RT:Ms07DCg7kMf41AY7IvJP7dO3O4f-7sYrgVhpIdWHgG8_gsxyYFYG-iInearKTU95P.Pzvnm7m5snMtTe0NDVRlyY2oG2zixS3uCodF.QJ92H7fq1dq1_tfsYkYiksktr@a26dc540-b231-47b5-9a86-c019d0bce528.hana.prod-us10.hanacloud.ondemand.com:443?encrypt=true"
	//hdbDsn = "hdb://SYSTEM:Plak8484@192.168.1.103:30015"
)

const (
	HOST     = "a26dc540-b231-47b5-9a86-c019d0bce528.hana.prod-us10.hanacloud.ondemand.com"
	PORT     = ":443"
	USERNAME = "2E403288791D4A56829B833EFA7A1529_3RGDF00CGM08B711H2QKDPA3F_RT"
	PASSWORD = "Ms07DCg7kMf41AY7IvJP7dO3O4f-7sYrgVhpIdWHgG8_gsxyYFYG-iInearKTU95P.Pzvnm7m5snMtTe0NDVRlyY2oG2zixS3uCodF.QJ92H7fq1dq1_tfsYkYiksktr"
	SCHEMA   = "2E403288791D4A56829B833EFA7A1529"
)

func main() {

	c := driver.NewBasicAuthConnector(
		HOST+PORT,
		USERNAME,
		PASSWORD)

	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
		ServerName:         HOST,
	}

	c.SetTLSConfig(&tlsConfig)

	//db, err := sql.Open(driverName, hdbDsn)
	db := sql.OpenDB(c)
	defer db.Close()

	var s = "Unassigned"

	// s = db.Conn.DefaultSchema()

	// fmt.Print("DefaultSchema: " + s + "\n")

	if err := db.QueryRow(fmt.Sprintf("SELECT NOW() FROM DUMMY")).Scan(&s); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Server Time: " + s + "\n")

	if err := db.QueryRow(fmt.Sprintf("SELECT VIEW_NAME FROM VIEWS WHERE VIEW_TYPE='CALC' LIMIT 1")).Scan(&s); err != nil {
		log.Fatal(err)
	}

	fmt.Print("First CalcView Found: " + s + "\n")

	_, err := db.Exec("SET SCHEMA " + SCHEMA)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.QueryRow(fmt.Sprintf("SELECT TITLE FROM MY_BOOKSHOP_BOOKS LIMIT 1")).Scan(&s); err != nil {
		log.Fatal(err)
	}

	fmt.Print(s + "\n")

	/*
	   SET SCHEMA SYS
	   SELECT * FROM M_TABLES WHERE SCHEMA_NAME='2E403288791D4A56829B833EFA7A1529'
	   SELECT * FROM VIEWS WHERE SCHEMA_NAME='2E403288791D4A56829B833EFA7A1529' AND NOT VIEW_TYPE='CALC'
	   SET SCHEMA 2E403288791D4A56829B833EFA7A1529
	   SELECT COLUMN_NAME, DATA_TYPE_NAME FROM TABLE_COLUMNS  WHERE SCHEMA_NAME='2E403288791D4A56829B833EFA7A1529' AND TABLE_NAME='MY_BOOKSHOP_BOOKS' AND INDEX_TYPE='NONE'
	   SELECT * FROM VIEW_COLUMNS  WHERE SCHEMA_NAME='2E403288791D4A56829B833EFA7A1529' AND VIEW_NAME='ORDERS_VIEW' AND NOT COLUMN_NAME='ID'
	   SELECT COLUMN_NAME, DATA_TYPE_NAME FROM VIEW_COLUMNS  WHERE SCHEMA_NAME='2E403288791D4A56829B833EFA7A1529' AND VIEW_NAME='ORDERS_VIEW' AND NOT COLUMN_NAME='ID'

	*/

	// http://go-database-sql.org/retrieving.html
	var (
		name     string
		col_type string
	)

	rows, err := db.Query("SELECT TABLE_NAME FROM M_TABLES WHERE SCHEMA_NAME='" + SCHEMA + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("TABLE: " + name + "\n")

		cols, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE_NAME FROM TABLE_COLUMNS  WHERE SCHEMA_NAME='" + SCHEMA + "' AND TABLE_NAME='" + name + "' AND INDEX_TYPE='NONE'")
		if err != nil {
			log.Fatal(err)
		}
		defer cols.Close()
		for cols.Next() {
			err := cols.Scan(&name, &col_type)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(" " + name + " : " + col_type + "\n")

		}
		err = cols.Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("" + "\n")
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("SELECT VIEW_NAME FROM VIEWS WHERE SCHEMA_NAME='" + SCHEMA + "' AND NOT VIEW_TYPE='CALC'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("VIEW: " + name + "\n")
		cols, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE_NAME FROM VIEW_COLUMNS  WHERE SCHEMA_NAME='" + SCHEMA + "' AND VIEW_NAME='" + name + "' AND NOT COLUMN_NAME='ID'")
		if err != nil {
			log.Fatal(err)
		}
		defer cols.Close()
		for cols.Next() {
			err := cols.Scan(&name, &col_type)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(" " + name + " : " + col_type + "\n")

		}
		err = cols.Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("" + "\n")
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
