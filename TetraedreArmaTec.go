package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TetraedreArmaTecdata struct {
	Time   time.Time
	Energy string
	//Volume    float64

}

type decodedTetraedreArmaTecdata struct {
	DeviceEui string                 `json:"deviceEui"`
	Seqno     uint32                 `json:"seqno"`
	Port      uint8                  `json:"port"`
	AppEui    string                 `json:"appEui"`
	Time      string                 `json:"time"`
	DeviceTx  devicetx               `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx            `json:"gatewayRx,omitempty"`
	Data      []TetraedreArmaTecdata `json:"data,omitempty"`
}

func parsedTetraedreArmaTecData(receivedtime time.Time, port uint8, receiveddata string) []TetraedreArmaTecdata {

	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	fmt.Println("databytes", databytes)
	fmt.Println("len(databytes)", len(databytes))

	if len(databytes)%18 != 0 {
		return nil
	}

	fmt.Println("databytes", databytes)

	capacity := len(databytes) / 18
	parsedvalues := make([]TetraedreArmaTecdata, capacity)
	for index := 0; index < capacity; index++ {
		fmt.Println("receiveddata", receiveddata)

		dst := ByteToHex(databytes)
		fmt.Println("Dec=", dst)

		first1 := dst[28:30]
		fmt.Println("first1", first1)
		first2 := dst[30:32]
		fmt.Println("first2", first2)

		first3 := dst[32:34]
		fmt.Println("first3", first3)

		first4 := dst[34:36]
		fmt.Println("first4", first4)

		swap1 := first4 + first3 + first2 + first1
		strval, err := strconv.Atoi(swap1)

		fmt.Println("swap1", strval)

		if err == nil {
			fmt.Println(err)
		}

		parsedvalues[index].Energy = fmt.Sprintf("%f", toFixed((float64(strval)), 3))

		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}

	return parsedvalues

}

func publishTetraedreArmaTecData(dev device, entry loradata, parsedvalues []TetraedreArmaTecdata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedTetraedreArmaTecdata
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
