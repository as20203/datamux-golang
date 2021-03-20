package main

import (
	"encoding/base64"

	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type LR210Data struct {
	Time        time.Time
	Relay1      string
	Relay2      string
	Temperature float64
}

type decodedlr210data struct {
	DeviceEui string      `json:"deviceEui"`
	Seqno     uint32      `json:"seqno"`
	Port      uint8       `json:"port"`
	AppEui    string      `json:"appEui"`
	Time      string      `json:"time"`
	DeviceTx  devicetx    `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx `json:"gatewayRx,omitempty"`
	Data      []LR210Data `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseLR210Data(receivedtime time.Time, port uint8, receiveddata string) []LR210Data {

	//Input Validation
	//Length should be a multiple of 6

	var parsedvalues []LR210Data
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	capacity := len(databytes) / 4

	parsedvalues = make([]LR210Data, capacity)

	//
	fmt.Println("databytes", databytes)
	dst := ByteToHex(databytes)
	fmt.Println("Hex=", dst)

	switch port {
	case 1: //status
	case 2: //periodic single measurement

		//00 01 04 4c
		//first1 := dst[0:2]
		first2 := dst[2:4]
		first3 := dst[4:6]
		first4 := dst[6:8]

		//fmt.Println("first1=",first1)
		//fmt.Println("first2=",first2)
		//fmt.Println("first3=",first3)
		//fmt.Println("first4=",first4)
		var strRev string
		var k float64
		str1 := first3 + first4
		strRev = fmt.Sprint(str1)

		sat1, _ := strconv.ParseInt(hexaNumberToInteger(strRev), 16, 64) //strconv.Atoi(strRev)
		k = float64(sat1)
		b := (k / 10)
		//relay1,err:=strconv.ParseInt(first1 ,16,64)

		relay2, _ := strconv.ParseInt(first2, 16, 64)
		//fmt.Println("relay1=",relay1)

		for index := 0; index < capacity; index++ {

			parsedvalues[index].Temperature = float64(b - 80)
			switch relay2 {
			case 0:
				parsedvalues[index].Relay1 = "Deactivated"
				parsedvalues[index].Relay2 = "Deactivated"

			case 1:

				parsedvalues[index].Relay1 = "Activated"
				parsedvalues[index].Relay2 = "Deactivated"

			case 2:

				parsedvalues[index].Relay1 = "Deactivated"
				parsedvalues[index].Relay2 = "Activated"
			case 3:

				parsedvalues[index].Relay1 = "Activated"
				parsedvalues[index].Relay2 = "Activated"

			default:
				parsedvalues[index].Relay1 = "Deactivated"
				parsedvalues[index].Relay2 = "Deactivated"

			}

			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

			// if err != nil {
			// 	// handle error
			// }
		}

	default:
		return nil
	}

	fmt.Println("parsedvalues: ", parsedvalues)
	return parsedvalues
}

func publishLR210Data(dev device, entry loradata, parsedvalues []LR210Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedlr210data
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
