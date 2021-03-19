package main

import (
	
	"time"
	 "encoding/base64"
	"encoding/json"
	"fmt"
	//"strconv"
)

type digimondodata struct {
	Time        time.Time
	MeterStatus string
	MeterReading    float64
   
}

type decodeddigimondodata struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []digimondodata  `json:"data,omitempty"`
}

func parsedDigimondoMeterData(receivedtime time.Time, port uint8, receiveddata string) []digimondodata {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)

	if len(databytes)%4 != 0 {
		return nil
	}

	//fmt.Println("databytes", databytes)
	capacity := len(databytes) / 4
	parsedvalues := make([]digimondodata, capacity)
	for index := 0; index < capacity; index++ {
	if int32(databytes[index*3])==3{
	parsedvalues[index].MeterStatus = "OK"
	
} else {
	parsedvalues[index].MeterStatus = "Not OK"
	}
	fmt.Println("receiveddata", receiveddata)
	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
		var sat float64
		var final string
		
		first2 := dst[2:8]
		fmt.Println("first2", first2)
		
		 final=first2
		
		fmt.Println("final", final)
		sat = hex2int(final)
		// strAbs=fmt.Sprint(uint32(databytes[(index*4)+1]))+fmt.Sprint(uint32(databytes[(index*4)+2]))+fmt.Sprint(uint32(databytes[(index*4)+3]))
		// sat, err :=strconv.Atoi(strAbs)
		 // if err != nil {
      // // handle error
   // }
		parsedvalues[index].MeterReading = sat // uint32(uint32(databytes[(index*4)+1]) << 4 uint32(databytes[(index*4)+2]) << 4 uint32(databytes[(index*4)+3]))
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}

	return parsedvalues

}


func publishDigimondoMeterData(dev device, entry loradata, parsedvalues []digimondodata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodeddigimondodata
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

