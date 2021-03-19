
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"math"
)

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

//type onyielddata struct {
//        Time        time.Time
//        Temperature float64 
//        Humidity    float64 
//        V       float64 
       
//}


type decoded1100data struct {
	DeviceEui string        `json:"deviceEui"`
	Seqno     uint32        `json:"seqno"`
	Port      uint8         `json:"port"`
	AppEui    string        `json:"appEui"`
	Time      string        `json:"time"`
	DeviceTx  devicetx      `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx   `json:"gatewayRx,omitempty"`
	Data      []onyielddata `json:"data,omitempty"`
}

func parseOY1110Data(receivedtime time.Time, port uint8, receiveddata string) []onyielddata {

	//Input Validation
	//Length should be a multiple of 6
 var expval, VS,F1,F2,F3,F4,F5,hamza,hadi float64 
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	var parsedvalues []onyielddata
	switch port {
	case 1: //status
	case 2: //periodic single measurement
		if len(databytes)%3 != 0 {
			return nil
		}
		//fmt.Println("databytes", databytes)
		capacity := len(databytes) / 3
		parsedvalues = make([]onyielddata, capacity)
		for index := 0; index < capacity; index++ {
			parsedvalues[index].Temperature = float64((int32(databytes[index*3])<<4|(int32(databytes[(index*3)+2])&0xF0)>>4)-800) / 10.0
			parsedvalues[index].RelativeHumidity = float64((int32(databytes[(index*3)+1])<<4|(int32(databytes[(index*3)+2])&0x0F))-250) / 10.0
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
            

 
			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
               
		
}
	case 3: //periodic group measurement
		if len(databytes)%3 != 1 {
			return nil
		}
		databytes = databytes[1:]
		capacity := len(databytes) / 3
		parsedvalues = make([]onyielddata, capacity)
		for index := 0; index < capacity; index++ {
			parsedvalues[index].Temperature = float64((int32(databytes[index*3])<<4|(int32(databytes[(index*3)+2])&0xF0)>>4)-800) / 10.0
			parsedvalues[index].RelativeHumidity = float64((int32(databytes[(index*3)+1])<<4|(int32(databytes[(index*3)+2])&0x0F))-250) / 10.0
			if parsedvalues[index].Temperature >= -50 &&  parsedvalues[index].Temperature <= -0.1 {
                hamza=2* math.Pow(37.230718,2)
				expval= ((math.Pow((108.19749 - parsedvalues[index].Temperature),2)) * -1)/hamza
                        fmt.Println("expval",expval ) 
						fmt.Println("math.Pow((108.19749 - parsedvalues[index].Temperature),2)",math.Pow((108.19749 - parsedvalues[index].Temperature),2) )
						
				VS=330.67796 * math.Exp(expval)
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
             parsedvalues[index].AbsoluteHumidity=fmt.Sprintf("%.2f", V1)
			}
            
       			parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
		}
	default:
		return nil
	}

	return parsedvalues
}

func publishOY1110Data(dev device, entry loradata, parsedvalues []onyielddata) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decoded1100data
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
