package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/cloudflare/cloudflare-go"
	"io"
	"net/http"
	"strings"
)

type Authorization struct {
	Token string `env:"CLOUDFLARE_API_TOKEN,required"`
	Email string `env:"CLOUDFLARE_EMAIL,required"`
}

type Dynamik struct {
	Authorization Authorization
	Client        *cloudflare.API
	ZoneId        string
	ZoneName      string `env:"CLOUDFLARE_ZONE_NAME,required"`
	RecordName    string `env:"CLOUDFLARE_RECORD_NAME,required"`
}

func NewDynamikClient() (*Dynamik, error) {
	var err error
	var d Dynamik
	d, err = env.ParseAs[Dynamik]()
	if err != nil {
		log.Fatal("error parsing environment variables", "fatal", err)
	}
	if d.IsEmpty() {
		return nil, errors.New("missing required environment variables")
	}
	cf, err := cloudflare.New(d.Authorization.Token, d.Authorization.Email)

	d.Client = cf
	d.ZoneId = d.GetZoneID()

	if err != nil {
		return nil, err
	}
	return &d, nil

}

func (d *Dynamik) IsEmpty() bool {
	res := []bool{d.ZoneId == "", d.ZoneName == "", d.RecordName == "", d.Authorization.Token == "", d.Authorization.Email == ""}
	for i := range res {
		if !res[i] {
			return false
		}
	}
	return true
}

func (d *Dynamik) GetZoneID() string {
	zoneID, err := d.Client.ZoneIDByName(d.ZoneName)
	if err != nil {
		log.Fatal(err)
	}
	return zoneID
}

func (d *Dynamik) GetZoneDnsRecords(ctx context.Context) []cloudflare.DNSRecord {
	records, err := d.Client.DNSRecords(ctx, d.ZoneId, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func (d *Dynamik) ParseForDynamicRecord() (string, string) {
	ctx := context.Background()
	for _, val := range d.GetZoneDnsRecords(ctx) {
		if val.Name == d.RecordName {
			return val.ID, val.Content
		}
	}
	return "", ""
}

func (d *Dynamik) CurrentIp() string {
	resp, err := http.Get("https://ifconfig.io")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (d *Dynamik) CheckIpMatches() bool {
	recordId, recordIp := d.ParseForDynamicRecord()
	currentIp := strings.TrimSpace(d.CurrentIp())

	if currentIp == recordIp {
		return true
	} else {
		err := d.UpdateDynamicDnsRecord(recordId, cloudflare.DNSRecord{Content: currentIp})
		if err != nil {
			log.Fatal("error updating DDNS record", "fatal", err)
		}
		return false
	}
}

func (d *Dynamik) UpdateDynamicDnsRecord(recordId string, value cloudflare.DNSRecord) error {
	ctx := context.Background()
	res := d.Client.UpdateDNSRecord(ctx, d.ZoneId, recordId, value)
	return res
}
