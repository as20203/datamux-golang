package main

import (
   // "bytes"
    // "strconv"
// "strings"	
 "encoding/base64"
	"encoding/json"
	"fmt"
	"time"
        "math"
"strconv"
)

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1600","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

type OY1600V1Data struct {
        Temperature float64
		RelativeHumidity    float64
        AbsoluteHumidity       string
	Time        time.Time
	ResistanceValue string

}


type decoded1600v1data struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []OY1600V1Data  `json:"data,omitempty"`
}

func parseOY1600V1Data(receivedtime time.Time, port uint8, receiveddata string) []OY1600V1Data {

	//Input Validation
	//Length should be a multiple of 2
         fmt.Println("dataval", receiveddata)
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
        var expval, VS,F1,F2,F3,F4,F5,hamza,hadi float64 
dst := ByteToHex(databytes)
fmt.Println("Hex=",dst)
dec:=hex2int(dst)
fmt.Println("Dec=",dec)

fmt.Println("message", len(databytes))

	var parsedvalues []OY1600V1Data
	switch port {
	case 1: //status
                
	case 2: //periodic single measurement
	
	
	       
	
	
	
	
		if len(databytes)%5 != 0 {
	        fmt.Println("databytesadil", databytes)	
        	return nil
		}
		fmt.Println("databytes", databytes)
               first2 := dst[6:8]
                last2  := dst[len(dst)-2:len(dst)]
                tt:=hex2int(first2)
                exp:=math.Floor(math.Mod(tt,128)/4)
		capacity := len(databytes) / 5
		parsedvalues = make([]OY1600V1Data, capacity)
		for index := 0; index < capacity; index++ {
            parsedvalues[index].Temperature = toFixed(float64((int32(databytes[index*5])<<4|(int32(databytes[(index*5)+2])&0xF0)>>4)-800) / 10.0,2)
			parsedvalues[index].RelativeHumidity = toFixed(float64((int32(databytes[(index*5)+1])<<4|(int32(databytes[(index*5)+2])&0x0F))-250) / 10.0,2)
			
		if parsedvalues[index].Temperature >= -50 &&  parsedvalues[index].Temperature <= -0.1 {
                hamza=2* math.Pow(37.230718,2)

				hadi= math.Pow((108.19749 - parsedvalues[index].Temperature),2) * -1
				fmt.Println("hadi",hadi )
				expval= hadi/hamza
                    fmt.Println("expval",expval ) 
               fmt.Println("math.Pow((108.19749 - parsedvalues[index].Temperature),2)",math.Pow((108.19749 - parsedvalues[index].Temperature),2) )					
				VS=330.67796* math.Exp(expval)
				parsedvalues[index].AbsoluteHumidity=fmt.Sprintf("%.2f", VS)
				fmt.Println("parsedvalues[index].absoluteHumidity",parsedvalues[index].AbsoluteHumidity )
                            

 
			}
			if  parsedvalues[index].Temperature >= 0.0 && parsedvalues[index].Temperature <=100.0 {
					F1=0.33229003 *  parsedvalues[index].Temperature
					F2=0.010508257*math.Pow(parsedvalues[index].Temperature,2)
					F3=0.00015035187*math.Pow(parsedvalues[index].Temperature,3)
					F4=0.0000021798571* math.Pow(parsedvalues[index].Temperature,4)
					F5=0.000000008613191 * math.Pow(parsedvalues[index].Temperature,5)
				
					 VS=4.8559296+F1+F2+F3+F4+F5
					V1:=(VS*parsedvalues[index].RelativeHumidity) / 100
					fmt.Println("parsedvalues[index].RelativeHumidity",parsedvalues[index].RelativeHumidity )
					parsedvalues[index].AbsoluteHumidity=fmt.Sprintf("%.2f", V1)
					fmt.Println("parsedvalues[index].absoluteHumidity",parsedvalues[index].AbsoluteHumidity )
			}
			
			//var floatBin = (int32(databytes[(index*5)+3]) << 8 | int32(databytes[(index*5)+4]) << 0);
			
                 fmt.Println("exp",exp)
               A2:=(1+(math.Mod(hex2int(first2),4)* 256 +hex2int(last2))/1024)
                fmt.Println("A2",A2)
               BB:=float64(math.Pow(2, exp - 15))
               fmt.Println("BB",BB)
              JJ:=(1-2 * math.Floor(float64(hex2int(first2)/128)))
 fmt.Println("JJ",JJ)            

 result:=JJ * BB * A2*1000
            
var result1=int(result)
                        parsedvalues[index].ResistanceValue = strconv.Itoa(result1)

                                            			
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}
	
	default:
		return nil
	}

	return parsedvalues
}
 // func hex2int(hexStr string) float64 {
 	// // remove 0x suffix if found in the input string
 	// cleaned := strings.Replace(hexStr, "0x", "", -1)

 	// // base 16 for hexadecimal
 	// result, _ := strconv.ParseUint(cleaned, 16, 64)
 	// return float64(result)
 // }

// func ByteToHex(data []byte) string {
    // buffer := new(bytes.Buffer)
    // for _, b := range data {

        // s := strconv.FormatInt(int64(b&0xff), 16)
        // if len(s) == 1 {
            // buffer.WriteString("0")
        // }
        // buffer.WriteString(s)
    // }

    // return buffer.String()
// }





func publishOY1600V1Data(dev device, entry loradata, parsedvalues []OY1600V1Data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decoded1600v1data
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
