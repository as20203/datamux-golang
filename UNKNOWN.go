package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type UNKNOWNData struct {
	Time        time.Time
	CO2Raw      uint16	
//	CO2Filtered    uint16
//	Temperature   float32
//	Humidity       float32
	
}

type decodedunknowndata struct {
	DeviceEui string       `json:"deviceEui"`
	Seqno     uint32       `json:"seqno"`
	Port      uint8        `json:"port"`
	AppEui    string       `json:"appEui"`
	Time      string       `json:"time"`
	DeviceTx  devicetx     `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx  `json:"gatewayRx,omitempty"`
	Data      []UNKNOWNData `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseUNKNOWNData(receivedtime time.Time, port uint8, receiveddata string) []UNKNOWNData {

	//Input Validation
	//Length should be a multiple of 6
//        fmt.Println("receiveddata",receiveddata)
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	var parsedvalues []UNKNOWNData
	
	fmt.Println("port", port)
//        fmt.Println("length", len(databytes))

		fmt.Println("databytes", databytes)
		
                parsedvalues = make([]UNKNOWNData, 1)
		
parsedvalues[0].CO2Raw = uint16(databytes[2])

//parsedvalues[0].CO2Filtered = uint16(databytes[4])<<8 + uint16(databytes[5])

//parsedvalues[0].Temperature = float32(uint16(databytes[6])<<8 + uint16(databytes[7]))/100.0
//parsedvalues[0].Humidity = float32(uint16(databytes[8])<<8 + uint16(databytes[9]))/100.0
			
			parsedvalues[0].Time = receivedtime.Add(time.Duration((-15)*0) * time.Minute)
		

	

	return parsedvalues
}

func publishUNKNOWNData(dev device, entry loradata, parsedvalues []UNKNOWNData) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedunknowndata
		if err := json.Unmarshal([]byte(loradatabytes), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
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
	fmt.Println("Data sent: ", string(loradatabytes))
		transferDatatoEndPoint(loradatabytes, dev)
	}

}
