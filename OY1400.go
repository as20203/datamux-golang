package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

/*
"payload": {"applicationID": "20","applicationName": "TalkpoolOY","deviceName": "OY1400","devEUI": "20-56-31-55-54-33-57-14","rxInfo": [{"mac": "C0-EE-40-FF-FF-29-3D-F8","rssi": -33,"loRaSNR": 7.2,"name": "GW name","latitude": "","longitude": "","altitude": ""}],"txInfo": {"frequency": "X","dataRate": {"modulation": "LORA","bandwidth": "X","spreadFactor": "X"},"adr": true,"codeRate": "X"},"fCnt": 25,"fPort": 1,"data": [{"time": 1547008071054,"temp": 23.7,"hum": 62.1},{"time": 1547005891054,"temp": 23.7,"hum": 62},{"time": 1547010251054,"temp": 23.6,"hum": 60.8}]}
*/

type oy1400data struct {
	Time       time.Time
	AnalogCh1  float64
	AnalogCh2  float64
	DigitalCh1 float64
	DigitalCh2 float64
}

type decoded1400data struct {
	DeviceEui string       `json:"deviceEui"`
	Seqno     uint32       `json:"seqno"`
	Port      uint8        `json:"port"`
	AppEui    string       `json:"appEui"`
	Time      string       `json:"time"`
	DeviceTx  devicetx     `json:"deviceTx,omitempty"`
	GatewayRx []gatewayrx  `json:"gatewayRx,omitempty"`
	Data      []oy1400data `json:"data,omitempty"`
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func parseOY1400Data(receivedtime time.Time, port uint8, receiveddata string) []oy1400data {

	//Input Validation
	//Length should be a multiple of 6
	//var expval, VS,F1,F2,F3,F4,F5 float64
	var val uint16
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)
	var parsedvalues []oy1400data
	switch port {
	case 1: //status
		//01 21 47 90 17 38

		if len(databytes)%6 == 0 {

			capacity := len(databytes) / 6
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)
				val = uint16(databytes[(index*6)])<<8 + uint16(databytes[(index*6)+1])

				fmt.Println("val", val)
				if val == 289 {

					parsedvalues[index].AnalogCh1 = toFixed(float64(uint16(databytes[(index*1)+2])<<8+uint16(databytes[(index*1)+3]))*0.5125, 1)

					parsedvalues[index].AnalogCh2 = toFixed(float64(uint16(databytes[(index*1)+4])<<8+uint16(databytes[(index*1)+5]))*0.5125, 1)

					parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

				}

			}
		}

		//01 21 1c 18 00 00 1b e0 00 00 1b e0 00 08

		if len(databytes)%14 == 0 {

			capacity := 3
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)

				parsedvalues[index].AnalogCh1 = toFixed(float64(uint16(databytes[(index*4)+2])<<8+uint16(databytes[(index*4)+3]))*0.5125, 1)

				parsedvalues[index].AnalogCh2 = toFixed(float64(uint16(databytes[(index*4)+4])<<8+uint16(databytes[(index*4)+5]))*0.5125, 1)

				parsedvalues[index].Time = receivedtime.Add(time.Duration((-5)*index) * time.Minute)

			}

		}

		// 01 22 01 17 38 00 21 88
		if len(databytes)%8 == 0 {

			capacity := len(databytes) / 4
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)
				val = uint16(databytes[(index*8)])<<8 + uint16(databytes[(index*8)+1])
				//val=(int32(databytes[index*8])<<4|(int32(databytes[(index*8)+1])&0xF0)>>4)
				fmt.Println("val", val)
				if val == 290 {
					parsedvalues[index].DigitalCh1 = toFixed(float64(uint16(databytes[(index*4)+2])<<8+uint16(databytes[(index*4)+3]))*0.5125, 1)
					fmt.Println("DigitalCh1", parsedvalues[index].DigitalCh1)
					parsedvalues[index].AnalogCh2 = toFixed(float64(uint16(databytes[(index*4)+4])<<8+uint16(databytes[(index*4)+6]))*0.5125, 1)
					fmt.Println("AnalogCh2", parsedvalues[index].AnalogCh2)
					parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

				}

			}

		}
		//01 23 47 90 00
		if len(databytes)%5 == 0 {

			capacity := len(databytes) / 5
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)
				val = uint16(databytes[(index*5)])<<8 + uint16(databytes[(index*5)+1])
				//val=(int32(databytes[index*5])<<4|(int32(databytes[(index*5)+1])&0xF0)>>4)
				fmt.Println("val", val)
				if val == 291 {
					fmt.Println(" came here")
					parsedvalues[index].AnalogCh1 = toFixed(float64(uint16(databytes[(index*5)+2])<<8+uint16(databytes[(index*5)+3]))*0.5125, 1)
					fmt.Println("AnalogCh1", parsedvalues[index].AnalogCh1)
					parsedvalues[index].DigitalCh2 = toFixed(float64(uint16(databytes[(index*5)+4])<<8+uint16(databytes[(index*5)+6]))*0.5125, 1)
					fmt.Println("DigitalCh2", parsedvalues[index].DigitalCh2)
					parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

				}

			}

		}

		// 01 24 01 00

		if len(databytes)%4 == 0 {

			capacity := len(databytes) / 4
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)
				val = uint16(databytes[(index*4)])<<8 + uint16(databytes[(index*4)+1])
				//val=(int32(databytes[index*4])<<4|(int32(databytes[(index*4)+1])&0xF0)>>4)
				fmt.Println("val", val)
				if val == 292 {
					parsedvalues[index].DigitalCh1 = toFixed(float64(uint16(databytes[(index*4)+2])), 0)
					fmt.Println("DigitalCh1", parsedvalues[index].DigitalCh1)
					parsedvalues[index].DigitalCh2 = toFixed(float64(uint16(databytes[(index*4)+3])), 0)
					// fmt.Println("DigitalCh2", parsedvalues[index].DigitalCh2)
					parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

				}

			}

		}

		//01 24 01 01 00 01 00 00 00 01

		if len(databytes)%10 == 0 {

			capacity := len(databytes) - 1
			parsedvalues = make([]oy1400data, capacity)
			for index := 0; index < capacity; index++ {
				fmt.Println("databytes", databytes)
				val = uint16(databytes[(index*2)])<<8 + uint16(databytes[(index*2)+1])
				//val=(int32(databytes[index*2])<<4|(int32(databytes[(index*2)+1])&0xF0)>>4)
				fmt.Println("val", val)
				if val == 292 {
					parsedvalues[index].DigitalCh1 = toFixed(float64(uint16(databytes[(index*2)+2])<<8+uint16(databytes[(index*2)+3]))*0.5125, 1)
					fmt.Println("DigitalCh1", parsedvalues[index].DigitalCh1)
					parsedvalues[index].DigitalCh2 = toFixed(float64(uint16(databytes[(index*2)+4])<<8+uint16(databytes[(index*2)+6]))*0.5125, 1)
					fmt.Println("DigitalCh2", parsedvalues[index].DigitalCh2)
					parsedvalues[index].Time = receivedtime.Add(time.Duration((-15)*index) * time.Minute)

				}

			}

		}
	default:
		return nil
	}
	return parsedvalues
}

func publishOY1400Data(dev device, entry loradata, parsedvalues []oy1400data) {
	loradatabytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Failed to encode message", err)
		return
	}

	if !dev.RawData {
		var decodeddata decoded1400data
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
