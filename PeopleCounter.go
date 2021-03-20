package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	// "encoding/binary"
	// "encoding/hex"
)

type PeopleCounterData struct {
	Time            time.Time
	DeviceStatus    uint16
	BatteryVoltage  uint16
	Counter_A       uint16
	Counter_B       uint16
	SensorStatus    uint16
	TotalCounter_A  uint16
	TotalCounter_B  uint16
	Payload_Counter uint16
}

type decodedpeoplecounterdata struct {
	DeviceEui string              `json:"deviceEui"`
	Seqno     uint32              `json:"seqno"`
	Port      uint8               `json:"port"`
	AppEui    string              `json:"appEui"`
	Time      string              `json:"time"`
	DeviceTx  devicetx            `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx         `json:"gatewayRx,omitempty"`
	Data      []PeopleCounterData `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parsePeopleCounterData(receivedtime time.Time, port uint8, receiveddata string) []PeopleCounterData {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	//	dst := ByteToHex(databytes)
	//	b, err := hex.DecodeString(dst)
	//if err != nil {
	// TODO: handle error
	//}

	//fmt.Println("Hex=",dst)
	var parsedvalues []PeopleCounterData
	switch port {
	//case 1: //status
	case 1: //periodic single measurement

		if len(databytes) != 23 {
			return nil
		}
		//fmt.Println("databytes", databytes)
		capacity := len(databytes) / 23
		parsedvalues = make([]PeopleCounterData, capacity)
		for index := 0; index < capacity; index++ {
			fmt.Println("Adil1")

			// parsedvalues[index].DeviceStatus = float32((int32(databytes[index*23])<<4|(int32(databytes[(index*23)+2])&0xF0)>>4)-800) / 10.0
			//parsedvalues[index].RelativeHumidity = float32((int32(databytes[(index*23)+1])<<4|(int32(databytes[(index*23)+2])&0x0F))-250) / 10.0
			parsedvalues[index].DeviceStatus = uint16(databytes[(index*23)+10])
			parsedvalues[index].BatteryVoltage = uint16(databytes[(index*23)+11])<<8 + uint16(databytes[(index*23)+12])
			parsedvalues[index].Counter_A = uint16(databytes[(index*23)+13])<<8 + uint16(databytes[(index*23)+14]) //binary.BigEndian.Uint16(b[:13])
			parsedvalues[index].Counter_B = uint16(databytes[(index*23)+15])<<8 + uint16(databytes[(index*23)+16]) //binary.BigEndian.Uint16(b[:15])
			parsedvalues[index].SensorStatus = uint16(databytes[(index*23)+17])
			parsedvalues[index].TotalCounter_A = uint16(databytes[(index*23)+18])<<8 + uint16(databytes[(index*23)+19]) //binary.BigEndian.Uint16(b[:18])
			parsedvalues[index].TotalCounter_B = uint16(databytes[(index*23)+20])<<8 + uint16(databytes[(index*23)+21])
			parsedvalues[index].Payload_Counter = uint16(databytes[(index*23)+22])
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
			fmt.Println("Adil2")

		}

	default:
		return nil
	}

	return parsedvalues
}

func publishPeopleCounterData(dev device, entry loradata, parsedvalues []PeopleCounterData) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedpeoplecounterdata
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
