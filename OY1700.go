package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type OY1700Data struct {
	Time             time.Time
	Temperature      float32
	RelativeHumidity float32
	PM10             uint16
	PM2_5            uint16
	PM1_0            uint16
}

type decoded1700data struct {
	DeviceEui string       `json:"deviceEui"`
	Seqno     uint32       `json:"seqno"`
	Port      uint8        `json:"port"`
	AppEui    string       `json:"appEui"`
	Time      string       `json:"time"`
	DeviceTx  devicetx     `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx  `json:"gatewayRx,omitempty"`
	Data      []OY1700Data `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseOY1700Data(receivedtime time.Time, port uint8, receiveddata string) []OY1700Data {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	var parsedvalues []OY1700Data
	switch port {
	case 1: //status
	case 2: //periodic single measurement

		if len(databytes) != 9 {
			return nil
		}
		//fmt.Println("databytes", databytes)
		capacity := len(databytes) / 9
		parsedvalues = make([]OY1700Data, capacity)
		for index := 0; index < capacity; index++ {
			fmt.Println("Adil1")
			parsedvalues[index].Temperature = float32((int32(databytes[index*9])<<4|(int32(databytes[(index*9)+2])&0xF0)>>4)-800) / 10.0
			parsedvalues[index].RelativeHumidity = float32((int32(databytes[(index*9)+1])<<4|(int32(databytes[(index*9)+2])&0x0F))-250) / 10.0
			parsedvalues[index].PM1_0 = uint16(databytes[(index*9)+3])<<8 + uint16(databytes[(index*9)+4])
			parsedvalues[index].PM2_5 = uint16(databytes[(index*9)+5])<<8 + uint16(databytes[(index*9)+6])
			parsedvalues[index].PM10 = uint16(databytes[(index*9)+7])<<8 + uint16(databytes[(index*9)+8])
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
			fmt.Println("Adil2")

		}
	case 3: //periodic group measurement
		if len(databytes)%9 != 1 {
			return nil
		}
		databytes = databytes[1:]
		capacity := len(databytes) / 7
		parsedvalues = make([]OY1700Data, capacity)
		for index := 0; index < capacity; index++ {

			parsedvalues[index].Temperature = float32((int32(databytes[index*9])<<4|(int32(databytes[(index*9)+2])&0xF0)>>4)-800) / 10.0
			parsedvalues[index].RelativeHumidity = float32((int32(databytes[(index*9)+1])<<4|(int32(databytes[(index*9)+2])&0x0F))-250) / 10.0
			parsedvalues[index].PM1_0 = uint16(databytes[(index*9)+3])<<8 + uint16(databytes[(index*9)+4])
			parsedvalues[index].PM2_5 = uint16(databytes[(index*9)+5])<<8 + uint16(databytes[(index*9)+6])
			parsedvalues[index].PM10 = uint16(databytes[(index*9)+7])<<8 + uint16(databytes[(index*9)+8])
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

		}
	default:
		return nil
	}

	return parsedvalues
}

func publishOY1700Data(dev device, entry loradata, parsedvalues []OY1700Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decoded1700data
		if err := json.Unmarshal([]byte(loradatabytes), &decodeddata); err != nil {
			fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
		}
		decodeddata.Data = append(decodeddata.Data, parsedvalues...)

		loradecodeddatabytes, err := json.Marshal(decodeddata)
		if err != nil {
			fmt.Println("Failed to encode message", err)
			return
		}

		fmt.Println("Data sent: ", string(loradecodeddatabytes))
		transferDatatoEndPoint(loradecodeddatabytes, dev)
	} else {
		transferDatatoEndPoint(loradatabytes, dev)
	}

}
