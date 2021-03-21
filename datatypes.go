package main

import (
	"time"
)

type device struct {
	Deviceeui     string    `bson:"deviceeui,omitempty"`
	Devicetype    string    `bson:"devicetype,omitempty"`
	Endpointtype  string    `bson:"endpointtype,omitempty"`
	Endpointdest  string    `bson:"endpointdest,omitempty"`
	AccessToken   string    `bson:"access_token,omitempty"`
	Customer      string    `bson:"customer,omitempty"`
	InclRadio     bool      `bson:"incl_radio,omitempty"`
	RawData       bool      `bson:"raw_data,omitempty"`
	LastUpdatedOn time.Time `bson:"last_updated_on,omitempty"`
}

/*{"deviceEui":"77-77-77-77-77-77-77-77","appEui":"70-B3-D5-D7-2F-F8-16-00","seqno":21,"port":2,"data":"PyaJ","time":"2019-02-05T02:24:03.543Z","deviceTx":{"sf":10,"bw":125,"freq":868.1,"adr":false},"gatewayRx":[{"gatewayEui":"CC-81-17-1D-98-38-E0-16","time":"2019-02-05T02:24:03.543Z","isTimeFromGateway":true,"chan":6,"rssi":-113,"snr":7.0}]}
 */

type devicetx struct {
	Freq float32 `json:"freq"`
	Sf   uint8   `json:"sf"`
	Bw   uint8   `json:"bw"`
	Adr  bool    `json:"adr"`
}

type gatewayrx struct {
	GatewayEui string  `json:"gatewayEui"`
	Chan       uint8   `json:"chan"`
	Rssi       int8    `json:"rssi"`
	Snr        float32 `json:"snr"`
}

type loradata struct {
	DeviceEui string      `json:"deviceEui"`
	Seqno     uint32      `json:"seqno"`
	Port      uint8       `json:"port"`
	AppEui    string      `json:"appEui"`
	Time      string      `json:"time"`
	Data      string      `json:"data"`
	DeviceTx  devicetx    `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx `json:"gatewayRx,omitempty"`
}
type downlinkdata struct {
	DeviceEui   string `json:"deviceEui"`
	Command     string `json:"command"`
	value       uint64 `json:"value"`
	AccessToken string `json:"access_token"`
	Devicetype  string `json:"devicetype"`
}
