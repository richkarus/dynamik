package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

func client() (*cloudflare.API, context.Context) {
	godotenv.Load()
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	return api, ctx

}

func PrettyStruct(data interface{}) string {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return string(val)
}

func GetZoneID(name string) string {
	api, _ := client()

	zoneID, err := api.ZoneIDByName(name)
	if err != nil {
		log.Fatal(err)
	}

	return zoneID
}

func GetZoneDetails(zoneId string) string {
	api, ctx := client()
	zone, err := api.ZoneDetails(ctx, zoneId)
	if err != nil {
		log.Fatal(err)
	}

	return PrettyStruct(zone)
}

func GetZoneDnsRecords(zoneId string) []cloudflare.DNSRecord {
	api, ctx := client()

	records, err := api.DNSRecords(ctx, zoneId, cloudflare.DNSRecord{})

	if err != nil {
		log.Fatal(err)
	}

	return records
}

func ParseForDynamicRecord(zoneId string, name string) (string, string) {
	for _, val := range GetZoneDnsRecords(zoneId) {
		if val.Name == name {
			return string(val.ID), string(val.Content)
		}
	}
	return "", ""
}

func GetCurrentIp() string {
	resp, err := http.Get("https://ifconfig.io")
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func CheckIpMatches() bool {
	godotenv.Load()
	zoneName := os.Getenv("CLOUDFLARE_ZONE_NAME")
	dynamicRecord := os.Getenv("CLOUDFLARE_DYNAMIC_RECORD_NAME")
	zoneId := GetZoneID(zoneName)
	recordId, recordIp := ParseForDynamicRecord(zoneId, dynamicRecord)
	currentIp := strings.TrimSpace(GetCurrentIp())

	if currentIp == recordIp {
		return true
	} else {
		UpdateDynamicDnsRecord(zoneId, recordId, cloudflare.DNSRecord{Content: currentIp})
		return false
	}
}

func UpdateDynamicDnsRecord(zoneId string, recordId string, value cloudflare.DNSRecord) error {
	api, ctx := client()
	res := api.UpdateDNSRecord(ctx, zoneId, recordId, value)

	return res
}

func GetRecordByName(zoneId string, recordName string) cloudflare.DNSRecord {
	records := GetZoneDnsRecords(zoneId)

	for _, r := range records {
		if r.Name == recordName {
			return r
		}
	}
	return cloudflare.DNSRecord{}
}
