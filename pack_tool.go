//go:build linux || windows

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	licenses := map[string]WebLicense{}
	f, err := os.OpenFile("users.csv", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open users.csv: %v\n", err)
		return
	}
	csv_data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Printf("Failed to read users.csv: %v\n", err)
		return
	}
	for _, row := range csv_data {
		l := WebLicense{
			LicenseId:  row[0],
			DeviceId:   row[1][0:7],
			ExpireTime: row[2],
			Products:   strings.Split(row[3], "|"),
			Timestamp:  row[6],
		}
		activeCode := row[4]
		licenses[activeCode] = l
	}

	//create output files
	//if licenses/ not exists, create it
	os.Mkdir("web_licenses", 0755)
	for activeCode, licenseInfo := range licenses {

		f, err := os.OpenFile(fmt.Sprintf("web_licenses/%s.dat", HexHash256([]byte(activeCode))), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to create %s.md: %v\n", HexHash256([]byte(activeCode)), err)
			return
		}
		defer f.Close()

		data, err := MarshalLicense(licenseInfo)
		if err != nil {
			fmt.Printf("Failed to marshal licenseInfos: %v\n", err)
			return
		}
		_, err = f.Write(data)
		if err != nil {
			fmt.Printf("Failed to write licenseInfos: %v\n", err)
			return
		}
		f.Close()
	}
	fmt.Printf("convert %d licenses\n", len(licenses))
}
