package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type SmartMC0101Doordata struct {
	Time           time.Time
	Battery        int64
	Temperature    int64
	Sending_Reason int64
	Input_State    int64
	Device_Time    string
	//Volume    float64

}

type decodedSmartMC0101Doordata struct {
	DeviceEui string                `json:"deviceEui"`
	Seqno     uint32                `json:"seqno"`
	Port      uint8                 `json:"port"`
	AppEui    string                `json:"appEui"`
	Time      string                `json:"time"`
	DeviceTx  devicetx              `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx           `json:"gatewayRx,omitempty"`
	Data      []SmartMC0101Doordata `json:"data,omitempty"`
}

func parsedSmartMC0101DoorData(receivedtime time.Time, port uint8, receiveddata string) []SmartMC0101Doordata {

	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	fmt.Println("databytes", databytes)
	fmt.Println("len(databytes)", len(databytes))

	fmt.Println("databytes", databytes)
	parsedvalues := make([]SmartMC0101Doordata, 1)
	fmt.Println("khalid=", "khalid")
	switch port {
	case 0:
	case 1: //status
	case 2:

		fmt.Println("Adil=", "Adil")
		if len(databytes)%11 == 0 {

			capacity := len(databytes) / 11

			for index := 0; index < capacity; index++ {
				fmt.Println("receiveddata", receiveddata)

				dst := ByteToHex(databytes)
				fmt.Println("Dec=", dst)
				first1 := dst[2:4]
				fmt.Println("first1", first1)
				first2 := dst[6:8]
				fmt.Println("first2", first2)

				first3 := dst[8:10]
				fmt.Println("first3", first3)

				first4 := dst[10:12]
				fmt.Println("first4", first4)
				first5 := dst[12:14]
				fmt.Println("first5", first5)

				//For Unix time

				first6 := dst[14:16]
				fmt.Println("first6", first6)
				first7 := dst[16:18]
				fmt.Println("first7", first7)
				first8 := dst[18:20]
				fmt.Println("first8", first8)
				first9 := dst[20:22]
				fmt.Println("first5", first9)

				swap1 := first9 + first8 + first7 + first6
				strtime := hex2dec(swap1)

				t := time.Unix(strtime, 0)
				strDate := t.Format(time.RFC3339)
				strbattery := hex2dec(first1)
				swap2 := first3 + first2
				strTemprature := hex2dec(swap2)
				strreason := hex2dec(first4)
				strinputstate := hex2dec(first5)
				fmt.Println("strTemprature", strTemprature)

				parsedvalues[index].Battery = strbattery
				parsedvalues[index].Temperature = strTemprature / 10
				parsedvalues[index].Sending_Reason = strreason
				parsedvalues[index].Input_State = strinputstate
				parsedvalues[index].Device_Time = strDate
				parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
			}

		}

	case 3:

	case 4:

		if len(databytes)%5 == 0 {

			capacity := len(databytes) / 5
			parsedvalues := make([]SmartMC0101Doordata, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("receiveddata", receiveddata)

				dst := ByteToHex(databytes)
				fmt.Println("Dec=", dst)
				first1 := dst[2:4]
				fmt.Println("first1", first1)
				first2 := dst[4:6]
				fmt.Println("first2", first2)

				first3 := dst[6:8]
				fmt.Println("first3", first3)

				first4 := dst[8:10]
				fmt.Println("first4", first4)

				swap1 := first4 + first3 + first2 + first1
				strtime := hex2dec(swap1)
				t := time.Unix(strtime, 0)
				strDate := t.Format(time.RFC3339)

				parsedvalues[index].Device_Time = strDate
				parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
			}

		}

	default:
		return nil
	}

	return parsedvalues

}

func hex2dec(hex string) int64 {
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		fmt.Println("Failed to convert message", err) //This error is ok as the format of data is different
	}

	return value
}

func publishSmartMC0101DoorData(dev device, entry loradata, parsedvalues []SmartMC0101Doordata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedSmartMC0101Doordata
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
