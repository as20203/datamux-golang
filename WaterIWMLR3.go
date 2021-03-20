package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type WaterIWMLR3Data struct {
	Time               time.Time
	MeterReading       float64
	ReverseFlowCounter float64
	KFactor            uint32
	Medium             uint32
	VIF                uint32

	Alarms uint32
}

type decodedwateriwmlr3data struct {
	DeviceEui string            `json:"deviceEui"`
	Seqno     uint32            `json:"seqno"`
	Port      uint8             `json:"port"`
	AppEui    string            `json:"appEui"`
	Time      string            `json:"time"`
	DeviceTx  devicetx          `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx       `json:"gatewayRx,omitempty"`
	Data      []WaterIWMLR3Data `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseWaterIWMLR3Data(receivedtime time.Time, port uint8, receiveddata string) []WaterIWMLR3Data {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)

	//44991400000000000000011300
	fmt.Println("databytes", databytes)
	dst := ByteToHex(databytes)
	fmt.Println("Hex=", dst)
	first := dst[2:4]
	first1 := dst[4:6]
	first2 := dst[6:8]
	first3 := dst[8:10]
	first4 := dst[10:12]
	first5 := dst[12:14]
	first6 := dst[14:16]
	first7 := dst[16:18]
	//first8 := dst[18:20]
	first10 := dst[22:24]
	fmt.Println("first=", first)
	fmt.Println("first1=", first1)
	fmt.Println("first2=", first2)
	fmt.Println("first3=", first3)
	fmt.Println("first4=", first4)
	fi10, _ := strconv.Atoi(first10)
	str := first3 + first2 + first1 + first
	str1 := first7 + first6 + first5 + first4
	// if err != nil {
	// 	// handle error
	// }
	var strAbs, strRev string
	var parsedvalues []WaterIWMLR3Data

	fmt.Println("len(databytes):", len(databytes))
	if len(databytes)%13 == 0 {
		// return nil
		//databytes = databytes[1:]
		capacity := len(databytes) / 13
		parsedvalues = make([]WaterIWMLR3Data, capacity)
		for index := 0; index < capacity; index++ {

			strAbs = fmt.Sprint(str)
			fmt.Println("strAbs=", strAbs)

			strRev = fmt.Sprint(str1)

			sat, _ := strconv.ParseInt(hexaNumberToInteger(strAbs), 10, 64)  //strconv.Atoi(strAbs)
			sat1, _ := strconv.ParseInt(hexaNumberToInteger(strRev), 10, 64) //strconv.Atoi(strRev)
			if uint32(databytes[(index*13)+9]) == 0 {

				// if err != nil {
				// 	// handle error
				// }
				parsedvalues[index].KFactor = 1
				parsedvalues[index].MeterReading = toFixed(float64(sat)*0.001, 4)
				parsedvalues[index].ReverseFlowCounter = toFixed(float64(sat1)*0.001, 4)
			} else if uint32(databytes[(index*13)+9]) == 1 {
				parsedvalues[index].KFactor = 10
				parsedvalues[index].MeterReading = toFixed(10.0*float64(sat)*0.001, 4)
				parsedvalues[index].ReverseFlowCounter = toFixed(10.0*float64(sat1)*0.001, 4)
			} else if uint32(databytes[(index*13)+9]) == 2 {
				parsedvalues[index].KFactor = 100
				parsedvalues[index].MeterReading = toFixed(100.0*float64(sat)*0.001, 4)
				parsedvalues[index].ReverseFlowCounter = toFixed(100.0*float64(sat1)*0.001, 4)
			}
			parsedvalues[index].Medium = uint32(databytes[(index*13)+10])
			// if(uint32(databytes[(index*13)+11])==19){

			// }else {
			// parsedvalues[index].Alarms=uint32(databytes[(index*13)+11])
			// }
			parsedvalues[index].VIF = uint32(fi10)
			parsedvalues[index].Alarms = uint32(databytes[(index*13)+12])
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}

	}

	//	fmt.Println("sat=",sat)
	//       fmt.Println("sat1=",sat1)

	fmt.Println("parsedvalues: ", parsedvalues)
	return parsedvalues
}

func publishWaterIWMLR3Data(dev device, entry loradata, parsedvalues []WaterIWMLR3Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedwateriwmlr3data
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
