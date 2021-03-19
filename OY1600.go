package main

import (
   "bytes"
    "strconv"
"strings"	
 "encoding/base64"
	"encoding/json"
	"fmt"
	"time"
        "math"
//"strconv"
)

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1600","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

type OY1600Data struct {
	Time        time.Time
	ResistanceValue string

}


type decoded1600data struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []OY1600Data  `json:"data,omitempty"`
}

func parseOY1600Data(receivedtime time.Time, port uint8, receiveddata string) []OY1600Data {

	//Input Validation
	//Length should be a multiple of 2
         fmt.Println("dataval", receiveddata)
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
        
dst := ByteToHex(databytes)
fmt.Println("Hex=",dst)
dec:=hex2int(dst)
fmt.Println("Dec=",dec)

fmt.Println("message", len(databytes))

	var parsedvalues []OY1600Data
	switch port {
	case 1: //status
                
	case 2: //periodic single measurement
		if len(databytes)%2 != 0 {
	        fmt.Println("databytesadil", databytes)	
        	return nil
		}
		fmt.Println("databytes", databytes)
                first2 := dst[0:2]
                last2  := dst[len(dst)-2:]
                tt:=hex2int(first2)
           //     fmt.Println(tt)                
//f,err := strconv.ParseFloat(tt, 64)
              //  fmt.Println(err)
                exp:=math.Floor(math.Mod(tt,128)/4)
		capacity := len(databytes) / 2
		parsedvalues = make([]OY1600Data, capacity)
		for index := 0; index < capacity; index++ {
               //        floatBin := (databytes[0] << 8 | databytes[1] << 0)
                fmt.Println("exp",exp)
              A2:=(1+(math.Mod(hex2int(first2),4)* 256 +hex2int(last2))/1024)
               fmt.Println("A2",A2)
              BB:=float64(math.Pow(2, exp - 15))
              fmt.Println("BB",BB)
             JJ:=(1-2 * math.Floor(float64(hex2int(first2)/128)))
fmt.Println("JJ",JJ)            
//KK:=float64((1+A2)/1024.0)*1000.0
//fmt.Println("KK",KK)

result:=JJ * BB * A2*1000
            

                        parsedvalues[index].ResistanceValue =fmt.Sprintf("%.2f",float64(result))
                                            			
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}
	
	default:
		return nil
	}

	return parsedvalues
}
 func hex2int(hexStr string) float64 {
 	// remove 0x suffix if found in the input string
 	cleaned := strings.Replace(hexStr, "0x", "", -1)

 	// base 16 for hexadecimal
 	result, _ := strconv.ParseUint(cleaned, 16, 64)
 	return float64(result)
 }

func ByteToHex(data []byte) string {
    buffer := new(bytes.Buffer)
    for _, b := range data {

        s := strconv.FormatInt(int64(b&0xff), 16)
        if len(s) == 1 {
            buffer.WriteString("0")
        }
        buffer.WriteString(s)
    }

    return buffer.String()
}





func publishOY1600Data(dev device, entry loradata, parsedvalues []OY1600Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decoded1600data
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
