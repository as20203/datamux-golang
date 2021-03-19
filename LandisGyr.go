package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
)

type LandisGyrData struct {
	Time        time.Time
	MeterReading float64
	VolumeM2 float64
	Power float64
	Flow float64
	ForwardTemperature float64
	ReturnTemperature float64
	
	Status  uint32

}


type decodedlandisgyrdata struct {
	DeviceEui string       `json:"deviceEui"`
	Seqno     uint32       `json:"seqno"`
	Port      uint8        `json:"port"`
	AppEui    string       `json:"appEui"`
	Time      string       `json:"time"`
	DeviceTx  devicetx     `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx  `json:"gatewayRx,omitempty"`
	Data      []LandisGyrData `json:"data,omitempty"`
}

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1100","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

func parseLandisGyrData(receivedtime time.Time, port uint8, receiveddata string) []LandisGyrData {



//00 0c 06   00 00 00 00 0c14010000000b2d0000000b3b1000f00a5a94010a5e96010c787859107002fd170000


	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
fmt.Println("databytes", databytes)
dst := ByteToHex(databytes)
fmt.Println("Hex=",dst)
first := dst[6:8]
firsthalf:=dst[8:10]
firsthalf1:=dst[10:12]
firsthalf2:=dst[12:14]

strMRHex:=firsthalf2+firsthalf1+firsthalf+first



first1 := dst[18:20]
first1half:= dst[20:22]
first1half2:= dst[22:24]
first1half3:= dst[24:26]

strVolHax :=first1half3+first1half2+first1half+first1


first2 := dst[30:32]
first2half := dst[32:34]
first2half1 := dst[34:36]

strPowerHax:=first2half1+first2half+first2

first3 := dst[40:42]
first3half := dst[42:44]
first3half1 := dst[44:46]
strFlowHax:=first3half1+first3half+first3


first4 := dst[50:52]
first4half := dst[52:54]

strFwdTempHex :=first4half+first4

first5 := dst[58:60]
first5half := dst[60:62]

strRtTemphex:=first5half+first5



strMR, err := strconv.ParseInt(strMRHex, 10, 64)
strVol, err := strconv.ParseInt(strVolHax, 10, 64)
strPower, err := strconv.ParseInt(strPowerHax, 10, 64)
strFlow, err := strconv.ParseInt(strFlowHax, 10, 64)
strFwdTemp, err := strconv.ParseInt(strFwdTempHex, 10, 64)
strRtTemp, err := strconv.ParseInt(strRtTemphex, 10, 64)
	if err == nil {
		fmt.Println(err)
	}
	
	var parsedvalues []LandisGyrData
        
              fmt.Println("len(databytes):", len(databytes))
                if len(databytes)%42 == 0 {
                       // return nil
				//databytes = databytes[1:]
                capacity := len(databytes) / 42
                parsedvalues = make([]LandisGyrData, capacity)
                for index := 0; index < capacity; index++ {
				fmt.Println("uint32(databytes[(index*42)+9]):", uint32(databytes[(index*42)+9]))
				fmt.Println("uint32(databytes[(index*42)+8]):", uint32(databytes[(index*42)+8]))
				fmt.Println("uint32(databytes[(index*42)+10]):", uint32(databytes[(index*42)+10]))
	
                             parsedvalues[index].MeterReading = toFixed(float64(strMR)*0.001,4)//float64(sat)
							 parsedvalues[index].VolumeM2 = toFixed(float64(strVol)*0.01,4)
							 parsedvalues[index].Power =toFixed(float64(strPower),2)
							 parsedvalues[index].Flow =toFixed(float64(strFlow),2)
							 parsedvalues[index].ForwardTemperature =toFixed(float64(strFwdTemp)*0.1,2)
							 parsedvalues[index].ReturnTemperature =toFixed(float64(strRtTemp)*0.1,2)
					
						
						
						parsedvalues[index].Status = 0
                        parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)
                		
                }
					
				
				}
	
	                    return parsedvalues

				}
              
		


func publishLandisGyrData(dev device, entry loradata, parsedvalues []LandisGyrData) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if dev.RawData == false {
		var decodeddata decodedlandisgyrdata
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
