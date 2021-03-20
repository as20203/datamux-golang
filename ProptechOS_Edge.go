package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func publishDataToProptechOS_Edge(buf []byte, dev1 device, sDevicetype string, sEndpointdest string, jsondata string) {

	var databuffer []byte
	var err error

	fmt.Println("RawData=", dev1.RawData)

	if !dev1.RawData {
		fmt.Println("RawData False aa Gaya")

		fmt.Println("dev1.Devicetype ", strings.ToLower(dev1.Devicetype))
		switch strings.ToLower(dev1.Devicetype) {
		case "oy1100":
		case "oy1110":
			var decodeddata decoded1100data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			}
			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1200":
			var decodeddata decoded1200data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}

		case "oy1210":
			var decodeddata decoded1210data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			fmt.Println("Ready1")
			databuffer, err = json.Marshal(decodeddata)
			fmt.Println("Ready2")
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
			fmt.Println("Ready")
			res1 := strings.Split(sEndpointdest, ";")

			fmt.Println("res1[0]", res1[0])
			fmt.Println("res1[1]", res1[1])
			fmt.Println("res1[2]", res1[2])
			res4 := strings.Split(res1[0], "=")
			res2 := strings.Split(res1[1], "=")
			res3 := strings.Split(res1[2], "=")
			fmt.Println("res2[0]", res2[0])
			fmt.Println("res2[1]", res2[1])
			fmt.Println("res3[1]", res3[1])
			mystr := res4[1] + "/devices/" + res2[1]
			fmt.Println("mystr=", mystr)
			completestr := "{\"format\":\"rec3.1.1\",\"deviceId\":\"" + fmt.Sprint(res2[1]) + "\" ,\"observations\": ["
			token := GenerateSasToken(mystr, res3[1]+"==", "360", "")
			sensors := strings.Split(dev1.AccessToken, ",")

			observationsstr := "{\"observationTime\":\"" + decodeddata.Time + "\",\"value\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Temperature) + ",\"quantityKind\":\"Temperature\",\"sensorId\":\"" + sensors[0] + "\"},{\"observationTime\":\"" + decodeddata.Time + "\",\"value\":" + fmt.Sprintf("%.2f", decodeddata.Data[0].Humidity) + ",\"quantityKind\":\"Humidity\",\"sensorId\":\"" + sensors[1] + "\"},{\"observationTime\":\"" + decodeddata.Time + "\",\"value\":" + fmt.Sprint(decodeddata.Data[0].Co2) + ",\"quantityKind\":\"CO2\",\"sensorId\":\"" + sensors[2] + "\"}]}"
			fmt.Println("token=", token)
			body := completestr + observationsstr
			fmt.Println("completestr=", body)
			var jsonStr1 = []byte(body)
			strURL := "https://" + mystr + "/messages/events?api-version=2018-06-30"
			fmt.Println("strURL:", strURL)
			req2, _ := http.NewRequest("POST", strURL, bytes.NewBuffer(jsonStr1))
			req2.Header.Set("Content-Type", "application/json")
			req2.Header.Set("Authorization", token)
			client := &http.Client{}
			resp2, _ := client.Do(req2)
			defer resp2.Body.Close()
			fmt.Println("response Status:", resp2.Status)

			body1, err := ioutil.ReadAll(resp2.Body)
			fmt.Println("body1:", body1)
			if err != nil {
				fmt.Println("Error reading body. ", err)
			}

		//defer resp2.Body.Close()

		case "oy1320":
			var decodeddata decoded1320data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1320v1":
			var decodeddata decoded1320v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1400":
			var decodeddata decoded1400data
			fmt.Println("Hamza 0")
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
			fmt.Println("Hamza1")
		case "oy1410":
			var decodeddata decoded1410data

			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}
			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1600":
			var decodeddata decoded1600data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1600v1":
			var decodeddata decoded1600v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}

		case "oy1700":
			var decodeddata decoded1700data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "oy1700v1":
			var decodeddata decoded1700v1data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "tetraedrembus":
			var decodeddata decodedtetraedrembusdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "wateriwmlr3":
			var decodeddata decodedwateriwmlr3data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}

		case "tds_100f_flow":
			var decodeddata decodedtds_100f_flowdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "honeywell_ew773":
			var decodeddata decodedHoneywell_EW773data
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}
		case "landisgyr":
			var decodeddata decodedlandisgyrdata
			if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
				//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
				return
			}

			databuffer, err = json.Marshal(decodeddata)
			if err != nil {
				fmt.Println("Failed to encode message", err)
				return
			}

		default:
			return
		}
		fmt.Println("databuffer ", databuffer)

	}

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler

	/*
			var jsonStr1 = []byte(body)
			strURL="https://"+mystr+"/messages/events?api-version=2018-06-30"
			 req2, err := http.NewRequest("POST", url2 , bytes.NewBuffer(jsonStr1))
		     req2.Header.Set("Content-Type", "application/json")
			  req2.Header.Set("Authorization", token)
		      client := &http.Client{}
			   resp2, err := client.Do(req2)
			  if err != nil {
				fmt.Println("Error reading response. ", err)
			}

		   defer resp2.Body.Close()
	*/
	/*


				parameter := make(map[string]interface{})
				parameter["campaign_id"] = "test_notify"
				parameter["content"] = map[string]string{"template_id": "xxxxxxxx"}

				data := make(map[string]interface{})

				data["address"] = "xxxx@xxxxx.com"
				data["substitution_data"] = map[string]string{
				  "address1":"xxxx@xxxxx.com",
				  "address2": "xxxx@xxxxx.com"}
				recipients := make([]map[string]interface{}, 0)
		recipients = append(recipients, data)
		parameter["recipients"] = recipients

				parameter["recipients"] = data
				fmt.Println(data)
				fmt.Println(parameter)

				mapC, _ := json.Marshal(parameter)
				fmt.Println(string(mapC))




			//data:="{observationTime":"2020-09-27T09:12:51.6874058Z","value":24.221417634196257,"quantityKind":"Temperature","sensorId":"07be0271-c8af-4a8c-a2e5-5922b4720c91"}"
			//recipients = append(recipients, data)
			fmt.Println("recipients=",recipients)

			//jsonData := map[string]string{"format": "rec3.1.1", "deviceId": "d4d36470-6fd2-4db2-99a1-13fed62ddb2d","observations":[{"observationTime":"2020-09-27T09:12:51.6874058Z","value":24.221417634196257,"quantityKind":"Temperature","sensorId":"07be0271-c8af-4a8c-a2e5-5922b4720c91"}]}

			HostString :="idun-bwprod-iothub-01.azure-devices.net/devices/d4d36470-6fd2-4db2-99a1-13fed62ddb2d/messages/events?api-version=2018-06-30"
			fmt.Println("HostString=",HostString)

			//SendJsonToIOTHUB(res2[1],"gK7/7sTDh4MIsZyowXatmw==" ,HostString,recipients )


	*/
}

func ComputeHmac256(message string, secret string) string {
	data, _ := base64.StdEncoding.DecodeString(secret)
	key := []byte(data)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func GenerateSasToken(resourceUri string, signingKey string, expiresInMins string, policyName string) string {
	uri := template.URLQueryEscaper(resourceUri)

	duration, _ := strconv.Atoi(expiresInMins)
	expire := time.Now().Add(time.Duration(duration) * time.Minute)
	fmt.Println("expire", expire)
	secs := expire.Unix()
	signed := uri + "\n" + strconv.FormatInt(secs, 10)

	val := ComputeHmac256(signed, signingKey)
	encoded_val := template.URLQueryEscaper(val)

	token := "SharedAccessSignature sr=" + uri + "&sig=" + encoded_val + "&se=" + strconv.FormatInt(secs, 10)
	if len(policyName) > 0 {
		token += "&skn=" + policyName
	}

	return token
}
func SendJsonToIOTHUB(deviceEUI string, deviceKey string, HostString string, recipients string) error {

	jsonValue, _ := json.Marshal(recipients)
	response, err := http.Post(HostString, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	data1 := string(data)
	fmt.Println(data1)
	res1 := strings.Split(data1, ":")
	res2 := strings.Replace(res1[2], "\",\"role\"", "", -1)
	res3 := strings.Replace(res2, "\"", "", -1)
	fmt.Println(res3)

	//this accepts both strings and int and works
	// 2nd Request
	//	var jsonStr1 = []byte("{\"deviceEui\":\""+deviceEUI+"\",\"type\" : \"data\",\"port\" : \"" + port+ "\",\"confirmedFrame\" : 1,\"data\" : \""+ dataval + "\"}")

	//	fmt.Println(bytes.NewBuffer(jsonStr1))
	//	var url2="https://ns.talkpool.com/internal/device/"+deviceEUI+"/downlink/data"
	//	fmt.Println(url2)

	//req2, err := http.NewRequest("POST", url2 , bytes.NewBuffer(jsonStr1))
	//	req2.Header.Set("Content-Type", "application/json")
	//	req2.Header.Set("X-AUTH-TOKEN", res3)
	//	client := &http.Client{}

	//	 for _, cookie := range response.Cookies() {
	//		req2.Header.Set("Cookie", cookie.String())
	//		fmt.Println(cookie.String())
	//	}
	//	 if err != nil {
	//        fmt.Println("Error1:", err)
	//    }

	//  resp2, err := client.Do(req2)
	//	  if err != nil {
	//		fmt.Println("Error reading response. ", err)
	//	}

	//  defer resp2.Body.Close()
	//	fmt.Println("response Status adil:", resp2.Status)

	//	body, err := ioutil.ReadAll(resp2.Body)
	//	if err != nil {
	//		fmt.Println("Error reading body. ", err)
	//	}

	//	fmt.Printf("this is body %s\n", body)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}
	return json.NewDecoder(response.Body).Decode(&body)

}
