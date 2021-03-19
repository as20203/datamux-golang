package main

import (
	"encoding/base64"
	"time"
)

type onyielddata struct {
	Time        time.Time
	Temperature float64
	RelativeHumidity    float64
        AbsoluteHumidity       string
}

func parseOnYieldData(receivedtime time.Time, receiveddata string) []onyielddata {

	//Input Validation
	//Length should be a multiple of 6
	databytes, _ := base64.StdEncoding.DecodeString(receiveddata)

	if len(databytes)%3 != 0 {
		return nil
	}

	//fmt.Println("databytes", databytes)
	capacity := len(databytes) / 3
	parsedvalues := make([]onyielddata, capacity)
	for index := 0; index < capacity; index++ {
		parsedvalues[index].Temperature = float64(int32(databytes[index*3])<<4|(int32(databytes[(index*3)+2])&0xF0)>>4) / 10.0
		parsedvalues[index].RelativeHumidity = float64(int32(databytes[(index*3)+1])<<4|(int32(databytes[(index*3)+2])&0x0F)) / 10.0
		parsedvalues[index].Time = receivedtime.Add(time.Duration((-2)*index) * time.Hour)
	}

	return parsedvalues
}
