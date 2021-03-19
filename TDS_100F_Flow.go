package main

import (
	
	"time"
	 "encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"unsafe"
	
)

type TDS_100F_Flowdata struct {
	Time        time.Time
	Temperature1 float64
	Temperature2    float64
	NetFlow float64
	NetEnergy float64
	
}

type decodedtds_100f_flowdata struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []TDS_100F_Flowdata  `json:"data,omitempty"`
}

func parsedTDS_100F_FlowData(receivedtime time.Time, port uint8, receiveddata string) []TDS_100F_Flowdata {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)

	/*if len(databytes)%29 != 0 {
		return nil
	}*/

	fmt.Println("len(databytes)", len(databytes))
	
	parsedvalues := make([]TDS_100F_Flowdata, len(databytes))
	
	
	if len(databytes)%29 == 0 {
	capacity := len(databytes) / 29
	parsedvalues := make([]TDS_100F_Flowdata, capacity)
	for index := 0; index < capacity; index++ {

//	fmt.Println("receiveddata", receiveddata)
	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
	fmt.Println("capacity=",capacity)
		
		first1 := dst[10:12]
	//	fmt.Println("first1", first1)
		first2 := dst[12:14]
		
//		fmt.Println("first2", first2)
		first3 := dst[14:16]
//		fmt.Println("first3", first3)
		
		first4 := dst[16:18]
	//	fmt.Println("first4", first4)
		swap1:=first3+first4+first1+first2
	//	fmt.Println("swap1", swap1)
		first5 := dst[18:20]
//		fmt.Println("first5", first5)
		first6 := dst[20:22]
//		fmt.Println("first6", first6)
		
		first7 := dst[22:24]
//		fmt.Println("first7", first7)
		first8 := dst[24:26]
//		fmt.Println("first8", first8)
		swap2:=first7+first8+first5+first6
//		fmt.Println("swap2", swap2)
		first9 := dst[26:28]
//		fmt.Println("first9", first9)
		first10 := dst[28:30]
//		fmt.Println("first10", first10)
		
		first11 := dst[30:32]
		
//		fmt.Println("first11", first11)
		first12 := dst[32:34]
//		fmt.Println("first12", first12)
		swap3:=first11+first12+first9+first10
//		fmt.Println("swap3", swap3)
		first13 := dst[34:36]
//		fmt.Println("first13", first13)
		first14 := dst[36:38]
//		fmt.Println("first14", first14)
	
		first15 := dst[38:40]
//		fmt.Println("first15", first15)
		first16 := dst[40:42]
//		fmt.Println("first16", first16)
		swap4:=first15+first16+first13+first14
//		fmt.Println("swap4", swap4)
		first17 := dst[42:44]
//		fmt.Println("first17", first17)
		first18 := dst[44:46]
//		fmt.Println("first18", first18)
		
		first19 := dst[46:48]
//		fmt.Println("first19", first19)
		first20 := dst[48:50]
//		fmt.Println("first20", first20)
		swap5:=first19+first20+first17+first18
//		fmt.Println("swap5", swap5)
		first21 := dst[50:52]
//		fmt.Println("first21", first21)
		first22 := dst[52:54]
//		fmt.Println("first22", first22)
		
		first23 := dst[54:56]
//		fmt.Println("first23", first23)
	
		first24 := dst[56:58]
		fmt.Println("first24", first24)
		swap6:=first23+first24+first21+first22
	    fmt.Println("swap6", swap6)
		fmt.Println("toFixed(hex2float(swap1),1)", toFixed(hex2float(swap1),1))
		fmt.Println("toFixed(hex2float(swap2),1)", toFixed(hex2float(swap2),1))
		fmt.Println("toFixed(hex2float(swap2),1)", toFixed(hex2float(swap2),1))
		fmt.Println("hex2int(swap1)=", hex2int(swap1))
	
		parsedvalues[index].Temperature1 = toFixed(hex2float(swap1),1) // uint32(uint32(databytes[(index*4)+1]) << 4 uint32(databytes[(index*4)+2]) << 4 uint32(databytes[(index*4)+3]))
		parsedvalues[index].Temperature2 = toFixed(hex2float(swap2),1)
		parsedvalues[index].NetFlow = toFixed(hex2int(swap3)+hex2float(swap4),3)
		parsedvalues[index].NetEnergy = toFixed(hex2int(swap5)+hex2float(swap6),3)
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}
	return parsedvalues
}	


if len(databytes)%13==0 {
		
		
	capacity := len(databytes) / 13
	
	
	for index := 0; index < capacity; index++ {

//	fmt.Println("receiveddata", receiveddata)
	
	dst := ByteToHex(databytes)
	fmt.Println("Dec=",dst)
		
	firstreject := dst[0:2]
	fmt.Println("firstreject=",firstreject)
		if  firstreject == "01" {
		return nil
	}
		first1 := dst[10:12]
	//	fmt.Println("first1", first1)
		first2 := dst[12:14]
		
//		fmt.Println("first2", first2)
		first3 := dst[14:16]
//		fmt.Println("first3", first3)
		
		first4 := dst[16:18]
	//	fmt.Println("first4", first4)
		swap1:=first3+first4+first1+first2
		fmt.Println("swap1", swap1)
		first5 := dst[18:20]
//		fmt.Println("first5", first5)
		first6 := dst[20:22]
//		fmt.Println("first6", first6)
		
		first7 := dst[22:24]
//		fmt.Println("first7", first7)
		first8 := dst[24:26]
//		fmt.Println("first8", first8)
		swap2:=first7+first8+first5+first6
		fmt.Println("swap2", swap2)
		
		 parsedvalues[index].Temperature1 = 0.0 // uint32(uint32(databytes[(index*4)+1]) << 4 uint32(databytes[(index*4)+2]) << 4 uint32(databytes[(index*4)+3]))
		 parsedvalues[index].Temperature2 = 0.0
		 
		parsedvalues[index].NetFlow = toFixed(hex2int(swap1)+hex2float(swap2),3)
		parsedvalues[index].NetEnergy = 	0.0
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}
		
		return parsedvalues
		
	}



	
return nil

}


 func hex2float(hexStr string) float64 {
 	n, err := strconv.ParseUint(hexStr, 16, 32)
if err != nil {
    panic(err)
}
var n2 = uint32(n)
f := *(*float32)(unsafe.Pointer(&n2))
 	return float64(f)
 }


func publishTDS_100F_FlowData(dev device, entry loradata, parsedvalues []TDS_100F_Flowdata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedtds_100f_flowdata
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

