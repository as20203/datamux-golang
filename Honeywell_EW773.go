package main

import (
	
	"time"
	 "encoding/base64"
	"encoding/json"
	"fmt"
		"strconv"
	
)

type Honeywell_EW773data struct {
	Time        time.Time
	Energy int
	Volume    float64

	
}

type decodedHoneywell_EW773data struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []Honeywell_EW773data  `json:"data,omitempty"`
}

func parsedhoneywellew773Data(receivedtime time.Time, port uint8, receiveddata string) []Honeywell_EW773data {

	
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
fmt.Println("databytes", databytes)
fmt.Println("len(databytes)", len(databytes))


	if len(databytes)%22 != 0 {
		return nil
	}

	fmt.Println("databytes", databytes)
	
	capacity := len(databytes) / 22
	parsedvalues := make([]Honeywell_EW773data, capacity)
	for index := 0; index < capacity; index++ {
	fmt.Println("receiveddata", receiveddata)
	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
	
		
		first1 := dst[24:26]
		fmt.Println("first1", first1)
		first2 := dst[26:28]
		fmt.Println("first2", first2)
		
		first3 := dst[28:30]
		fmt.Println("first3", first3)
		
		first4 := dst[30:32]
		fmt.Println("first4", first4)
		check14:= dst[34:36]
		fmt.Println("check14", check14)
		first5 := dst[36:38]
		fmt.Println("first5", first5)
		first6 := dst[38:40]
		fmt.Println("first6", first6)
		first7 := dst[40:42]
		fmt.Println("first7", first7)
		first8 := dst[42:44]
		fmt.Println("first8", first8)
		swap1:=first4+first3+first2+first1
		swap2:=first8+first7+first6+first5
		//first2 := dst[35:44]
		//reverse(first1)
		//reverse(first2)
	    fmt.Println("swap1", swap1)
	    fmt.Println("swap2", swap2)
	
		
		
		// strAbs=fmt.Sprint(uint32(databytes[(index*4)+1]))+fmt.Sprint(uint32(databytes[(index*4)+2]))+fmt.Sprint(uint32(databytes[(index*4)+3]))
		// sat, err :=strconv.Atoi(strAbs)
		 // if err != nil {
      // // handle error
   // }
      strMR, err := strconv.Atoi(swap1)
	  strval, err := strconv.Atoi(swap2)
	  if err == nil {
		fmt.Println(err)
	}
		parsedvalues[index].Energy =  strMR
		if check14=="14" {
		parsedvalues[index].Volume =  toFixed((float64(strval)*0.01),2)
		} else {
		parsedvalues[index].Volume =  toFixed((float64(strval)*0.001),3)
		}
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}
		
		
		
	
	
	
	
	
	
	
	

	return parsedvalues

}




func publishHoneywell_EW773Data(dev device, entry loradata, parsedvalues []Honeywell_EW773data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedHoneywell_EW773data
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

