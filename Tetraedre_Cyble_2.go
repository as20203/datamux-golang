package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	//"strconv"
)

type Tetraedre_Cyble_2data struct {
	Time         time.Time
	MeterReading string
	//Volume    float64

}

type decodedTetraedre_Cyble_2data struct {
	DeviceEui string                  `json:"deviceEui"`
	Seqno     uint32                  `json:"seqno"`
	Port      uint8                   `json:"port"`
	AppEui    string                  `json:"appEui"`
	Time      string                  `json:"time"`
	DeviceTx  devicetx                `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx             `json:"gatewayRx,omitempty"`
	Data      []Tetraedre_Cyble_2data `json:"data,omitempty"`
}

func parsedTetraedre_Cyble_2Data(receivedtime time.Time, port uint8, receiveddata string) []Tetraedre_Cyble_2data {

	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	fmt.Println("databytes", databytes)
	fmt.Println("len(databytes)", len(databytes))

	if len(databytes)%18 != 0 {
		return nil
	}

	fmt.Println("databytes", databytes)

	capacity := len(databytes) / 18
	parsedvalues := make([]Tetraedre_Cyble_2data, capacity)
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
		//strval, err := strconv.ParseInt(swap1, 10, 64)//strconv.Atoi(swap1)
		strval := hex2int(swap1)
		fmt.Println("Osman=", strval)

		parsedvalues[index].MeterReading = fmt.Sprintf("%f", toFixed((float64(strval)*0.001), 3))

		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}

	return parsedvalues

}

func publishTetraedre_Cyble_2Data(dev device, entry loradata, parsedvalues []Tetraedre_Cyble_2data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decodedTetraedre_Cyble_2data
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
