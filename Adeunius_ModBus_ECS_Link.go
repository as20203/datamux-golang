package main

import (
	
	"time"
	 "encoding/base64"
	"encoding/json"
	
	"fmt"
	"unsafe"
	"strconv"

	
)

type Adeunius_ModBus_ECS_Linkdata struct {
	Time        time.Time
	MeterReading string
	

	
}

type decodedAdeunius_ModBus_ECS_Linkdata struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []Adeunius_ModBus_ECS_Linkdata  `json:"data,omitempty"`
}

func parsedAdeunius_ModBus_ECS_LinkData(receivedtime time.Time, port uint8, receiveddata string) []Adeunius_ModBus_ECS_Linkdata {

	
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
fmt.Println("databytes", databytes)
fmt.Println("len(databytes)", len(databytes))


	if len(databytes)%6 != 0 {
		return nil
	}


	
	capacity := len(databytes) / 6
	parsedvalues := make([]Adeunius_ModBus_ECS_Linkdata, capacity)
	for index := 0; index < capacity; index++ {
	fmt.Println("receiveddata", receiveddata)

	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
	
		
		first1 := dst[4:12]
	//	data, err := hex.DecodeString(first1)
		
		
	n, err := strconv.ParseUint(first1, 16, 32)
	if err != nil {
		panic(err)
	}

	n2 := uint32(n)
	f := *(*float32)(unsafe.Pointer(&n2))
		

		parsedvalues[index].MeterReading =  fmt.Sprintf("%f",toFixed(float64(f),3))
		//parsedvalues[index].Temperature=fmt.Sprintf("%f",toFixed((float64(strTemperature)*0.01),3))
	
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}
		
		
		

	return parsedvalues

}




func publishAdeunius_ModBus_ECS_LinkData(dev device, entry loradata, parsedvalues []Adeunius_ModBus_ECS_Linkdata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedAdeunius_ModBus_ECS_Linkdata
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
		transferDatatoEndPoint(loradatabytes, dev)
	}

}

