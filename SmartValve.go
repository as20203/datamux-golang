package main

import (
	"encoding/base64"

	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type SmartValveData struct {
	Time             time.Time
	RemainingVoltage uint32
	//ValveIsConnected string

	VALVEPOS    string
	TAMPER      string
	CABLE       string
	DI_0        string
	DI_1        string
	LEAKAGE     string
	FRAUD       string
	Temperature float64
	Hygrometry  float64
}

type decodedsmartvalvedata struct {
	DeviceEui string           `json:"deviceEui"`
	Seqno     uint32           `json:"seqno"`
	Port      uint8            `json:"port"`
	AppEui    string           `json:"appEui"`
	Time      string           `json:"time"`
	DeviceTx  devicetx         `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx      `json:"gatewayRx,omitempty"`
	Data      []SmartValveData `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseSmartValveData(receivedtime time.Time, port uint8, receiveddata string) []SmartValveData {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)

	//33 30 32 34 35 23 5d 65 96 c9
	fmt.Println("databytes", databytes)
	dst := ByteToHex(databytes)
	fmt.Println("Hex=", dst)

	//for forword

	first := dst[2:4]
	first1 := dst[5:6]
	first2 := dst[7:8]

	//ValveDisconnected Status
	first3 := dst[8:10]

	//first4 := dst[10:12]
	first5 := dst[12:14]
	first6 := dst[14:16]
	first7 := dst[16:18]

	first8 := dst[18:20]

	fi, _ := strconv.Atoi(first + first1 + first2)
	//fi1,err:=strconv.Atoi(first1)
	//fi2,err:=strconv.Atoi(first2)
	fi3, _ := strconv.ParseInt(first3, 16, 64)
	//fi4,err:=strconv.Atoi(first4)
	//fi5,err:=strconv.Atoi(first5)
	//fi6,err:=strconv.Atoi(first6)
	//fi7,err:=strconv.Atoi(first7)
	//fi8,err:=strconv.Atoi(first8)

	fmt.Println("fi3=", fi3)

	bval := strconv.FormatInt((fi3), 2)
	fmt.Println("bval=", bval)
	isvalveOpen := string([]rune(bval)[5])
	istemper := string([]rune(bval)[4])
	iscable := string([]rune(bval)[3])
	isDI_0 := string([]rune(bval)[2])
	isDI_1 := string([]rune(bval)[1])
	isLEAKAGE := string([]rune(bval)[0])
	isFRAUD := string([]rune(bval)[0])
	// if err != nil {
	// 	// handle error
	// }
	//fmt.Println("first=",first)
	//fmt.Println("first1=",first1)
	//fmt.Println("first2=",first2)
	//fmt.Println("first3=",first3)
	//fmt.Println("first4=",first4)
	//fmt.Println("first5=",first5)
	//fmt.Println("first6=",first6)
	//fmt.Println("first7=",first7)
	//fmt.Println("first8=",first8)
	str := first5 + first6
	fmt.Println("str=", str)
	str1 := first7 + first8
	//fmt.Println("str1=",str1)

	var strAbs, strRev string
	var parsedvalues []SmartValveData

	//fmt.Println("len(databytes):", len(databytes))
	if len(databytes)%10 == 0 {
		// return nil
		//databytes = databytes[1:]
		capacity := len(databytes) / 10
		parsedvalues = make([]SmartValveData, capacity)
		for index := 0; index < capacity; index++ {

			strAbs = fmt.Sprint(str)
			//fmt.Println("strAbs=",strAbs)

			strRev = fmt.Sprint(str1)

			sat, _ := strconv.ParseInt(hexaNumberToInteger(strAbs), 16, 64)  //strconv.Atoi(strAbs)
			sat1, _ := strconv.ParseInt(hexaNumberToInteger(strRev), 16, 64) //strconv.Atoi(strRev)
			//fmt.Println("sat=",sat)
			var j, k float64
			//fmt.Println("sat1=",sat1)
			j = float64(sat)
			k = float64(sat1)

			a := (j / 65536) * 165
			//fmt.Println("a=",a)
			b := (k / 65536)
			//fmt.Println("b=",b)
			parsedvalues[index].Temperature = float64(a - 40)

			parsedvalues[index].Hygrometry = float64(b * 100)

			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

			// if err != nil {
			// 	// handle error
			// }

			parsedvalues[index].RemainingVoltage = uint32(fi)
			if isvalveOpen == "0" {
				parsedvalues[index].VALVEPOS = "Valve Closed"
			} else {
				parsedvalues[index].VALVEPOS = "Valve Opened"
			}
			if istemper == "0" {
				parsedvalues[index].TAMPER = "Enclosure Closed"
			} else {
				parsedvalues[index].TAMPER = "Enclosure Opened"
			}
			if iscable == "0" {
				parsedvalues[index].CABLE = "Cable Disconnected"
			} else {
				parsedvalues[index].CABLE = "Cable Connected"
			}
			if isDI_0 == "0" {
				parsedvalues[index].DI_0 = "Digital Input 0 is OFF"
			} else {
				parsedvalues[index].DI_0 = "Digital Input 0 is ON"
			}
			if isDI_1 == "0" {
				parsedvalues[index].DI_1 = "Digital Input 1 is OFF"
			} else {
				parsedvalues[index].DI_1 = "Digital Input 1 is ON"
			}
			if isLEAKAGE == "0" {
				parsedvalues[index].LEAKAGE = "Leak Detected"
			} else {
				parsedvalues[index].LEAKAGE = "Leak Not Detected"
			}
			if isFRAUD == "0" {
				parsedvalues[index].FRAUD = " Fraud Detected"
			} else {
				parsedvalues[index].FRAUD = "Fraud Not Detected"
			}

		}

	}

	//fmt.Println("parsedvalues: ", parsedvalues)
	return parsedvalues
}

func publishSmartValveData(dev device, entry loradata, parsedvalues []SmartValveData) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedsmartvalvedata
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
