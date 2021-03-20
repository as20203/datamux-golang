package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type OY1320V1Data struct {
	Time         time.Time
	MeterReading float64
	Status       string
}

type decoded1320v1data struct {
	DeviceEui string         `json:"deviceEui"`
	Seqno     uint32         `json:"seqno"`
	Port      uint8          `json:"port"`
	AppEui    string         `json:"appEui"`
	Time      string         `json:"time"`
	DeviceTx  devicetx       `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx    `json:"gatewayRx,omitempty"`
	Data      []OY1320V1Data `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseOY1320V1Data(receivedtime time.Time, port uint8, receiveddata string) []OY1320V1Data {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	fmt.Println("databytes", databytes)
	dst := ByteToHex(databytes)
	fmt.Println("Hex=", dst)

	first := dst[4:12]
	fmt.Println("MeterHex=", first)
	strMR, err := strconv.ParseInt(hexaNumberToInteger(first), 16, 64)

	if err == nil {
		fmt.Println(err)
	}
	var parsedvalues []OY1320V1Data

	fmt.Println("len(databytes):", len(databytes))
	if len(databytes)%6 == 0 {

		capacity := len(databytes) / 6
		parsedvalues = make([]OY1320V1Data, capacity)
		for index := 0; index < capacity; index++ {
			fmt.Println("MeterIntegerValue=", strMR)
			parsedvalues[index].MeterReading = toFixed(float64(strMR)*0.001, 4) //uint32(databytes[(index*6)+2])<<8 + uint32(databytes[(index*6)+3])<<8 + uint32(databytes[(index*6)+4])<<8 + uint32(databytes[(index*6)+5])<<8
			parsedvalues[index].Status = "0"
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}

	}
	fmt.Println("len(databytes)%9:", len(databytes)%9)

	if len(databytes)%9 == 0 {
		// return nil
		//databytes = databytes[1:]
		capacity := len(databytes) / 9
		parsedvalues = make([]OY1320V1Data, capacity)
		for index := 0; index < capacity; index++ {
			fmt.Println("MeterIntegerValue=", strMR)
			parsedvalues[index].MeterReading = toFixed(float64(strMR)*0.001, 4) //uint32(databytes[(index*9)+2])<<8 + uint32(databytes[(index*9)+3])<<8 + uint32(databytes[(index*9)+4])<<8 + uint32(databytes[(index*9)+5])<<8
			parsedvalues[index].Status = fmt.Sprint(databytes[8])
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}

	}

	fmt.Println("parsedvalues: ", parsedvalues)
	return parsedvalues
}

func publishOY1320V1Data(dev device, entry loradata, parsedvalues []OY1320V1Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decoded1320v1data
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
