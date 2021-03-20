package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	//	"strconv"
)

type TetraedraEDCBAData struct {
	Time         time.Time
	MeterReading string
}

type decodedTetraedraEDCBAdata struct {
	DeviceEui string               `json:"deviceEui"`
	Seqno     uint32               `json:"seqno"`
	Port      uint8                `json:"port"`
	AppEui    string               `json:"appEui"`
	Time      string               `json:"time"`
	DeviceTx  devicetx             `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx          `json:"gatewayRx,omitempty"`
	Data      []TetraedraEDCBAData `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseTetraedraEDCBADataData(receivedtime time.Time, port uint8, receiveddata string) []TetraedraEDCBAData {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	fmt.Println("databytes", databytes)
	dst := ByteToHex(databytes)
	fmt.Println("Hex=", dst)
	//first := dst[24:26]
	first1 := dst[26:27]
	first2 := dst[28:30]
	first3 := dst[30:32]
	first4 := dst[32:34]
	first5 := dst[34:36]

	fmt.Println("first1=", first1)
	fmt.Println("first2=", first2)
	fmt.Println("first3=", first3)
	fmt.Println("first4=", first4)
	fmt.Println("first5=", first5)

	str := first5 + first4 + first3 + first2 + first1
	/*
	   str1, err := strconv.ParseInt(str, 10, 64)
	   	if err == nil {
	   		fmt.Println(str1)
	   	}

	   fmt.Println("str=",str)
	   fmt.Println(str1)
	*/
	var parsedvalues []TetraedraEDCBAData

	fmt.Println("len(databytes):", len(databytes))
	if len(databytes)%18 == 0 {
		// return nil
		//databytes = databytes[1:]
		capacity := len(databytes) / 18
		parsedvalues = make([]TetraedraEDCBAData, capacity)
		for index := 0; index < capacity; index++ {

			parsedvalues[index].MeterReading = str
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}

	}

	fmt.Println("parsedvalues: ", parsedvalues)
	return parsedvalues
}

func publishTetraedraEDCBAData(dev device, entry loradata, parsedvalues []TetraedraEDCBAData) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedTetraedraEDCBAdata
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
