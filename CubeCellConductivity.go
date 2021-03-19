package main

import (
	
	"time"
	 "encoding/base64"
	"encoding/json"
	"fmt"
		
	
)

type CubeCellConductivitydata struct {
	Time        time.Time
	Conductivity string
	Temperature    string

	
}

type decodedCubeCellConductivitydata struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []CubeCellConductivitydata  `json:"data,omitempty"`
}

func parsedCubeCellConductivityData(receivedtime time.Time, port uint8, receiveddata string) []CubeCellConductivitydata {

	
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
fmt.Println("databytes", databytes)
fmt.Println("len(databytes)", len(databytes))


	if len(databytes)%4 != 0 {
		return nil
	}


	
	capacity := len(databytes) / 4
	parsedvalues := make([]CubeCellConductivitydata, capacity)
	for index := 0; index < capacity; index++ {
	fmt.Println("receiveddata", receiveddata)
	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
	
		
		first1 := dst[0:4]
		fmt.Println("first1", first1)
		first2 := dst[4:8]
		fmt.Println("first2", first2)
		
		
		strConductivity :=  toFixed(hex2int(first1),3)  //strconv.Atoi(first1)
	    strTemperature :=   toFixed(hex2int(first2),3)
	 
	  
	
		
		parsedvalues[index].Conductivity =  fmt.Sprintf("%f",toFixed((float64(strConductivity)),3))
		parsedvalues[index].Temperature=fmt.Sprintf("%f",toFixed((float64(strTemperature)*0.01),3))
	
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}
		
		
		

	return parsedvalues

}




func publishCubeCellConductivityData(dev device, entry loradata, parsedvalues []CubeCellConductivitydata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedCubeCellConductivitydata
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

