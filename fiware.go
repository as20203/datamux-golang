package main

import (
	"encoding/json"
	"fmt"
	"strings"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type JSONdata struct {
	Deviceeui   string  `json:"Deviceeui"`
	Temperature float32 `json:"Temperature"`
	Humidity    float32 `json:"Humidity"`
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func publishDataToFiware(buf []byte, dev1 device,sDevicetype string,sEndpointdest string) {

	var databuffer []byte
	var err error

	switch strings.ToLower(sDevicetype) {
	case "oy1700":
		var decodeddata decoded1700data
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}

		databuffer, err = json.Marshal(decodeddata.Data[0])
		if err != nil {
			fmt.Println("Failed to encode message", err)
			return
		}

        case "oy1600":
                var decodeddata decoded1600data
                if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
                        //fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
                        return
                }

                databuffer, err = json.Marshal(decodeddata.Data[0])
                if err != nil {
                        fmt.Println("Failed to encode message", err)
                        return
                }


	case "oy1210":
		var decodeddata decoded1210data
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}
		databuffer, err = json.Marshal(decodeddata.Data[0])
		if err != nil {
			fmt.Println("Failed to encode message", err)
			return
		}
	case "oy1100":
	case "oy1110":
		var decodeddata decoded1100data
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}
		databuffer, err = json.Marshal(decodeddata.Data[0])
		if err != nil {
			fmt.Println("Failed to encode message", err)
			return
		}
	default:
		return
	}

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(sEndpointdest)
	opts.SetClientID("LoraDataPublisher")
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername("talkpool")
	opts.SetPassword("talkpool")

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	defer c.Disconnect(300)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Publishing")
	//Publish messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	postdata := string(databuffer)
	//postdata := "{ \"d\" :" + jsonStr + "}"
	fmt.Println("Data: ", postdata)
	topic := "/lora/" + dev1.Deviceeui + "/igress"
	token := c.Publish(topic, 0, false, postdata)
	fmt.Println("Token --> ", token)
	token.Wait()
	//time.Sleep(3 * time.Second)

	c.Disconnect(250)
}
