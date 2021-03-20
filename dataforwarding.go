package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func dec2hex(dec uint16) string {

	return fmt.Sprintf("%02x", dec)
}

// func casetdevicetx(Foobar devicetx) string {
//     return fmt.Sprintf("%v", Foobar)
// }
// func casetgatewayrx(Foobar gatewayrx) string {
//     return fmt.Sprintf("%v", Foobar)
// }
func transferDatatoEndPoint(buf []byte, dev device) {

	fmt.Println("Send Data for device  ", dev, buf)

	res1 := strings.Count(dev.Endpointtype, "|")

	if res1 == 0 {

		pushdata(buf, dev.Endpointtype, dev.Endpointdest, dev)

	} else {
		sEndpointty := strings.Split(dev.Endpointtype, "|")
		sEndpointDy := strings.Split(dev.Endpointdest, "|")
		for ir := 0; ir <= res1; ir++ {
			fmt.Println("ir=", ir)
			fmt.Println("Before Calling Push Function sEndpointty[ir] and sEndpointDy[ir]", "sEndpointty[ir]="+sEndpointty[ir]+"    sEndpointDy[ir]="+sEndpointDy[ir])
			pushdata(buf, sEndpointty[ir], sEndpointDy[ir], dev)

		}

	}
}

func pushdata(buf []byte, sEndpointty string, sEndpointdest string, dev1 device) {

	fmt.Println("sEndpointty=", sEndpointty)
	fmt.Println("sEndpointdest=", sEndpointdest)
	switch strings.ToLower(sEndpointty) {
	case "http":
		req, _ := http.NewRequest("POST", sEndpointdest, bytes.NewBuffer(buf))
		//req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return

	case "sensepool":
		var body string
		switch strings.ToLower(dev1.Devicetype) {
		case "oy1210":
			var decodeddata decoded1210data

			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			}
			strCarbondioxide := dec2hex(decodeddata.Data[0].Co2)
			fmt.Println("strCarbondioxide=", strCarbondioxide)
			strTemprature := dec2hex(uint16(decodeddata.Data[0].Temperature * 100))
			fmt.Println("strTemprature=", strTemprature)
			strHumidity := dec2hex(uint16(decodeddata.Data[0].Humidity * 100))
			fmt.Println("strHumidity=", strHumidity)

			strdata := "0121" + "0" + strCarbondioxide + "0" + strCarbondioxide + "0" + strTemprature + "" + strHumidity + "0000"
			length := len(strdata)
			fmt.Println("length=", length)
			if length == 23 {
				strdata = "0121" + "0" + strCarbondioxide + "0" + strCarbondioxide + "0" + strTemprature + "0" + strHumidity + "0000"
			}

			fmt.Println("strdata=", strdata)
			decodedHex, err := hex.DecodeString(strdata)
			fmt.Println("DeviceTx=", decodeddata.DeviceTx)

			if err != nil {
				panic(err)
			}
			//fmt.Println("decodedHex:", string(decodedHex))
			encodedHex := hex.EncodeToString(decodedHex)
			fmt.Println("encodedHex:", encodedHex)

			fmt.Println()
			encodedBase64 := base64.StdEncoding.EncodeToString(decodedHex)
			fmt.Println("encodedBase64:", encodedBase64)

			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"appEui\":\"" + decodeddata.AppEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"data\":\"" + fmt.Sprint(encodedBase64) + "\",\"time\":\"" + decodeddata.Time + "\", \"deviceTx\":{\"sf\":7,\"bw\":125,\"freq\":867.9,\"adr\":true},\"gatewayRx\":[{\"gatewayEui\":\"64-7F-DA-FF-FE-00-7A-B3\",\"time\":\"2020-09-14T23:57:12.778Z\",\"isTimeFromGateway\":false,\"chan\":7,\"rssi\":-63,\"snr\":10.8}]} "
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		default:
			return

		}

	case "thingsboard":
		var body string
		switch strings.ToLower(dev1.Devicetype) {
		case "oy1100":
			var decodeddata decoded1100data

			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			}

			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].RelativeHumidity) + ",\"absoluteHumidity\":" + decodeddata.Data[0].AbsoluteHumidity + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "oy1110":
			var decodeddata decoded1100data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			}

			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].RelativeHumidity) + ",\"absoluteHumidity\":" + decodeddata.Data[0].AbsoluteHumidity + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "oy1200":
			var decodeddata decoded1200data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Humidity) + ",\"CO2Raw\":" + fmt.Sprint(decodeddata.Data[0].CO2Raw) + ",\"CO2Filtered\":" + fmt.Sprint(decodeddata.Data[0].CO2Filtered) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1210":
			var decodeddata decoded1210data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Humidity) + ",\"Co2\":" + fmt.Sprint(decodeddata.Data[0].Co2) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "ONFrXZ2jbRiq6p2LMcKQ")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "adeunius_modbus_ecs_link":
			var decodeddata decodedAdeunius_ModBus_ECS_Linkdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + "}"
			//body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) +",\"Humidity\":"+fmt.Sprintf("%.2f",  decodeddata.Data[0].Humidity) +",\"Co2\":"+fmt.Sprintf("%.2f",  decodeddata.Data[0].Co2) +"}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1320":
			var decodeddata decoded1320data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"errorStatus\":" + fmt.Sprint(decodeddata.Data[0].Status) + "}"
			//body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) +",\"Humidity\":"+fmt.Sprintf("%.2f",  decodeddata.Data[0].Humidity) +",\"Co2\":"+fmt.Sprintf("%.2f",  decodeddata.Data[0].Co2) +"}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1320v1":
			var decodeddata decoded1320v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"errorStatus\":" + fmt.Sprint(decodeddata.Data[0].Status) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedretillquistdiz":
			var decodeddata decodedTETRAEDRETillquistDIZdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"errorStatus\":" + fmt.Sprint(decodeddata.Data[0].ErrorStatus) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedregwfmtk3":
			var decodeddata decodedTetraedreGWFMTK3data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Volume\":" + fmt.Sprint(decodeddata.Data[0].Volume) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedrelogt550":
			var decodeddata decodedTetraedreLoGT550data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Energy\":" + fmt.Sprint(decodeddata.Data[0].Energy) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedrearmatec":
			var decodeddata decodedTetraedreArmaTecdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Energy\":" + fmt.Sprint(decodeddata.Data[0].Energy) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "cubecellconductivity":
			var decodeddata decodedCubeCellConductivitydata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Conductivity\":" + fmt.Sprint(decodeddata.Data[0].Conductivity) + ",\"Temperature\":" + fmt.Sprint(decodeddata.Data[0].Temperature) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedre_cyble_2":
			var decodeddata decodedTetraedre_Cyble_2data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedraedcba":
			var decodeddata decodedTetraedraEDCBAdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "tetraedre_abb_b24":
			var decodeddata decodedTETRAEDRE_ABB_B24data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"errorStatus\":" + fmt.Sprint(decodeddata.Data[0].ErrorStatus) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "oy1410":
			var decodeddata decoded1410data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Pulses\":" + fmt.Sprint(decodeddata.Data[0].Pulses) + ",\"Status\":" + fmt.Sprint(decodeddata.Data[0].Status) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1400":
			var decodeddata decoded1400data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			length := len(decodeddata.Data)
			fmt.Println("length", length)
			if length == 1 {
				body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"analogCh1\":" + fmt.Sprint(decodeddata.Data[0].AnalogCh1) + ",\"analogCh2\":" + fmt.Sprint(decodeddata.Data[0].AnalogCh2) + ",\"digitalCh1\":" + fmt.Sprint(decodeddata.Data[0].DigitalCh1) + ",\"digitalCh2\":" + fmt.Sprint(decodeddata.Data[0].DigitalCh2) + ", \"Time\":\"" + fmt.Sprint(decodeddata.Data[0].Time) + "\"}"
				req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
				if err != nil {
					// handl err
					fmt.Println("Error1")

				}
				fmt.Println("body:", body)

				req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					// handle err
					fmt.Println("Error 2")

				}
				fmt.Println("response Status:", resp.Status)
				fmt.Println("response Headers:", resp.Header)
				defer resp.Body.Close()
			}
			if length == 3 {
				body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"analogCh1\":" + fmt.Sprint(decodeddata.Data[0].AnalogCh1) + ",\"analogCh2\":" + fmt.Sprint(decodeddata.Data[0].AnalogCh2) + ",\"digitalCh1\":" + fmt.Sprint(decodeddata.Data[0].DigitalCh1) + ",\"digitalCh2\":" + fmt.Sprint(decodeddata.Data[0].DigitalCh2) + ", \"Time1\":\"" + fmt.Sprint(decodeddata.Data[0].Time) + "\",\"analogCh1.1\":" + fmt.Sprint(decodeddata.Data[1].AnalogCh1) + ",\"analogCh2.1\":" + fmt.Sprint(decodeddata.Data[1].AnalogCh2) + ",\"digitalCh1.1\":" + fmt.Sprint(decodeddata.Data[1].DigitalCh1) + ",\"digitalCh2.1\":" + fmt.Sprint(decodeddata.Data[1].DigitalCh2) + ", \"Time2\":\"" + fmt.Sprint(decodeddata.Data[1].Time) + "\",\"analogCh1.3\":" + fmt.Sprint(decodeddata.Data[2].AnalogCh1) + ",\"analogCh2.3\":" + fmt.Sprint(decodeddata.Data[2].AnalogCh2) + ",\"digitalCh1.3\":" + fmt.Sprint(decodeddata.Data[2].DigitalCh1) + ",\"digitalCh2.3\":" + fmt.Sprint(decodeddata.Data[2].DigitalCh2) + ", \"Time3\":\"" + fmt.Sprint(decodeddata.Data[2].Time) + "\"}"
				req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
				if err != nil {
					// handl err
					fmt.Println("Error1")

				}
				fmt.Println("body:", body)

				req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					// handle err
					fmt.Println("Error 2")

				}
				fmt.Println("response Status:", resp.Status)
				fmt.Println("response Headers:", resp.Header)
				defer resp.Body.Close()
			}

		case "oy1600v1":
			var decodeddata decoded1600v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"resistanceValue\":" + decodeddata.Data[0].ResistanceValue + ",\"absoluteHumidity\":" + decodeddata.Data[0].AbsoluteHumidity + ",\"relativeHumidity\":" + fmt.Sprintf("%f", decodeddata.Data[0].RelativeHumidity) + ",\"Temperature\":" + fmt.Sprintf("%f", decodeddata.Data[0].Temperature) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "smartvalve":
			var decodeddata decodedsmartvalvedata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"remainingVoltage\":" + fmt.Sprint(decodeddata.Data[0].RemainingVoltage) + ",\"ValvePOS\":\"" + fmt.Sprint(decodeddata.Data[0].VALVEPOS) + "\",\"Tamper\":\"" + fmt.Sprint(decodeddata.Data[0].TAMPER) + "\",\"Cable\":\"" + fmt.Sprint(decodeddata.Data[0].CABLE) + "\",\"DI_0\":\"" + fmt.Sprint(decodeddata.Data[0].DI_0) + "\",\"DI_1\":\"" + fmt.Sprint(decodeddata.Data[0].DI_1) + "\",\"Leakage\":\"" + fmt.Sprint(decodeddata.Data[0].LEAKAGE) + "\",\"Fraud\":\"" + fmt.Sprint(decodeddata.Data[0].FRAUD) + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"Hygrometry\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Hygrometry) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "lr210":
			var decodeddata decodedlr210data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Relay1\":\"" + fmt.Sprint(decodeddata.Data[0].Relay1) + "\",\"Relay2\":\"" + fmt.Sprint(decodeddata.Data[0].Relay2) + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "oy1600":
			var decodeddata decoded1600data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"ResistanceValue\":" + decodeddata.Data[0].ResistanceValue + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1700":
			var decodeddata decoded1700data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].RelativeHumidity) + ",\"PM10\":" + fmt.Sprint(decodeddata.Data[0].PM10) + ",\"PM2.5\":" + fmt.Sprint(decodeddata.Data[0].PM2_5) + ",\"PM1.0\":" + fmt.Sprint(decodeddata.Data[0].PM1_0) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "oy1700v1":
			var decodeddata decoded1700v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Temperature\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"relativeHumidity\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].RelativeHumidity) + ",\"PM10\":" + fmt.Sprint(decodeddata.Data[0].PM10) + ",\"PM2.5\":" + fmt.Sprint(decodeddata.Data[0].PM2_5) + ",\"PM1.0\":" + fmt.Sprint(decodeddata.Data[0].PM1_0) + ",\"PMCount0.3\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT0_3) + ",\"PMCount0.5\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT0_5) + ",\"PMCount1.0\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT1_0) + ",\"PMCount2.5\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT2_5) + ",\"PMCount10.0\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT10_0) + ",\"PMCount5.0\":" + fmt.Sprint(decodeddata.Data[0].PMCOUNT5_0) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "peoplecounter":
			var decodeddata decodedpeoplecounterdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"DeviceStatus\":" + fmt.Sprint(decodeddata.Data[0].DeviceStatus) + ",\"BatteryVoltage\":" + fmt.Sprint(decodeddata.Data[0].BatteryVoltage) + ",\"Counter_A\":" + fmt.Sprint(decodeddata.Data[0].Counter_A) + ",\"Counter_B\":" + fmt.Sprint(decodeddata.Data[0].Counter_B) + ",\"SensorStatus\":" + fmt.Sprint(decodeddata.Data[0].SensorStatus) + ",\"TotalCounter_A\":" + fmt.Sprint(decodeddata.Data[0].TotalCounter_A) + ",\"TotalCounter_B\":" + fmt.Sprint(decodeddata.Data[0].TotalCounter_B) + ",\"Payload_Counter\":" + fmt.Sprint(decodeddata.Data[0].Payload_Counter) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "wateriwmlr3":
			var decodeddata decodedwateriwmlr3data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"reverseFlowCounter\":" + fmt.Sprint(decodeddata.Data[0].ReverseFlowCounter) + ",\"Medium\":" + fmt.Sprint(decodeddata.Data[0].Medium) + ",\"VIF\":" + fmt.Sprint(decodeddata.Data[0].VIF) + ",\"KFactor\":" + fmt.Sprint(decodeddata.Data[0].KFactor) + ",\"Alarms\":" + fmt.Sprint(decodeddata.Data[0].Alarms) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "digimondometer":
			var decodeddata decodeddigimondodata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"Status\":" + fmt.Sprint(decodeddata.Data[0].MeterStatus) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		case "tds_100f_flow":
			var decodeddata decodedtds_100f_flowdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"DeviceTemperature1\":" + fmt.Sprint(decodeddata.Data[0].Temperature1) + ",\"DeviceTemperature2\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature2) + ",\"NetFlow\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].NetFlow) + ",\"NetEnergy\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].NetEnergy) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "honeywell_ew773":
			var decodeddata decodedHoneywell_EW773data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Energy\":" + fmt.Sprint(decodeddata.Data[0].Energy) + ",\"Volume\":" + fmt.Sprint(decodeddata.Data[0].Volume) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			fmt.Println("body:", body)
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "smartmc0101door":
			var decodeddata decodedSmartMC0101Doordata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"Battery\":" + fmt.Sprint(decodeddata.Data[0].Battery) + ",\"Temperature\":" + fmt.Sprint(decodeddata.Data[0].Temperature) + ",\"sendingReason\":" + fmt.Sprint(decodeddata.Data[0].Sending_Reason) + ",\"Input_State\":" + fmt.Sprint(decodeddata.Data[0].Input_State) + ",\"deviceTime\":\"" + fmt.Sprint(decodeddata.Data[0].Device_Time) + "\"}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()

		case "landisgyr":
			var decodeddata decodedlandisgyrdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "{\"deviceEui\":\"" + decodeddata.DeviceEui + "\",\"seqno\":" + fmt.Sprint(decodeddata.Seqno) + " ,\"port\":" + fmt.Sprint(decodeddata.Port) + ",\"appEui\":\"" + decodeddata.AppEui + "\",\"time\":\"" + decodeddata.Time + "\",\"meterReading\":" + fmt.Sprint(decodeddata.Data[0].MeterReading) + ",\"volumeM2\":" + fmt.Sprint(decodeddata.Data[0].VolumeM2) + ",\"Power\":" + fmt.Sprint(decodeddata.Data[0].Power) + ",\"Flow\":" + fmt.Sprint(decodeddata.Data[0].Flow) + ",\"forwardTemperature\":" + fmt.Sprint(decodeddata.Data[0].ForwardTemperature) + ",\"returnTemperature\":" + fmt.Sprint(decodeddata.Data[0].ReturnTemperature) + ",\"Status\":" + fmt.Sprint(decodeddata.Data[0].Status) + "}"
			req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
			if err != nil {
				// handl err
				fmt.Println("Error1")

			}
			fmt.Println("body:", body)

			req.SetBasicAuth("token", "qJv4XZ0W7h9ALBMb4Wi6")
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				fmt.Println("Error 2")

			}
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			defer resp.Body.Close()
		default:
			return

		}

	case "corlysis":
		var body string
		switch strings.ToLower(dev1.Devicetype) {
		case "oy1700":
			var decodeddata decoded1700data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "readings,device=" + decodeddata.DeviceEui + " Temperature=" + fmt.Sprintf("%f", decodeddata.Data[0].Temperature) +
				",relativeHumidity=" + fmt.Sprintf("%f", decodeddata.Data[0].RelativeHumidity) +
				",PM10=" + fmt.Sprintf("%d", decodeddata.Data[0].PM10) +
				",PM2.5=" + fmt.Sprintf("%d", decodeddata.Data[0].PM2_5) +
				",PM1.0=" + fmt.Sprintf("%d", decodeddata.Data[0].PM1_0)

		case "oy1210":
			var decodeddata decoded1210data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			body = "readings,device=" + decodeddata.DeviceEui + " Temperature=" + fmt.Sprintf("%f", decodeddata.Data[0].Temperature) +
				",Humidity=" + fmt.Sprintf("%f", decodeddata.Data[0].Humidity) +
				",Co2=" + fmt.Sprintf("%d", decodeddata.Data[0].Co2)

			fmt.Println("Adil" + body)
			/*case "OY1320":
			  var decodeddata decoded1320data
			  if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			          fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			          return
			  }
			  body = "readings,device=" + decodeddata.DeviceEui + " Temperature=" + fmt.Sprintf("%f", decodeddata.Data[0].Temperature) +
			          ",Humidity=" + fmt.Sprintf("%f", decodeddata.Data[0].Humidity) +
			          ",Co2=" + fmt.Sprintf("%d", decodeddata.Data[0].Co2) +
			          ",Time=" + fmt.Sprintf("%d", decodeddata.Data[0].Time)
			          fmt.Println("Adil"+body)*/

		default:
			return
		}
		fmt.Println("Nabeel" + body)

		req, err := http.NewRequest("POST", sEndpointdest, strings.NewReader(body))
		if err != nil {
			// handl err
			fmt.Println("Error1")

		}
		req.SetBasicAuth("token", "026656139cc88855b805e0934069392a")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
			fmt.Println("Error 2")

		}
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		defer resp.Body.Close()
		return
	case "proptechos_edge":

		publishDataToProptechOS_Edge(buf, dev1, sEndpointty, sEndpointdest, string(buf))
		return
	case "mqtt":
		publishDataToMQTT(buf, dev1, sEndpointty, sEndpointdest, string(buf))

		return

	case "fiware":
		publishDataToFiware(buf, dev1, sEndpointty, sEndpointdest)
		return

	default:
		fmt.Println("Unsupported Endpoint Type")
	}

}
