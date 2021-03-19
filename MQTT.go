package main

import (
	"encoding/json"
	"fmt"
	"strings"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)





//define a function for the default message handler
// var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	// fmt.Printf("TOPIC: %s\n", msg.Topic())
	// fmt.Printf("MSG: %s\n", msg.Payload())
// }

func publishDataToMQTT(buf []byte, dev1 device,sDevicetype string,sEndpointdest string,jsondata string) {

	var databuffer []byte
	var err error
	fmt.Println("dev1.RawData", dev1.RawData)
if  dev1.RawData==false {
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
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}
		 databuffer, err = json.Marshal(decodeddata)
		 if err != nil {
			 fmt.Println("Failed to encode message", err)
			 return
		 }
	 
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
	case "peoplecounter":
		var decodeddata decodedpeoplecounterdata
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}

		 databuffer, err = json.Marshal(decodeddata)
		 if err != nil {
			 fmt.Println("Failed to encode message", err)
			 return
		 }
	case "tetraedretillquistdiz":
		var decodeddata decodedTETRAEDRETillquistDIZdata
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}

		 databuffer, err = json.Marshal(decodeddata)
		 if err != nil {
			 fmt.Println("Failed to encode message", err)
			 return
		 }
	case "tetraedre_abb_b24":
		var decodeddata decodedTETRAEDRE_ABB_B24data
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
		 
    case "tetraedregwfmtk3":
		var decodeddata decodedTetraedreGWFMTK3data
		if err := json.Unmarshal([]byte(buf), &decodeddata); err != nil {
			//fmt.Println("Failed to encode message", err) //This error is ok as the format of data is different
			return
		}

		 databuffer, err = json.Marshal(decodeddata)
		 if err != nil {
			 fmt.Println("Failed to encode message", err)
			 return
		 }
    case "tetraedrelogt550":
		var decodeddata decodedTetraedreLoGT550data
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
	
}


	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	
	res1 := strings.Split(sEndpointdest, "*") 
	
	fmt.Println("res1[0]",res1[0])
	fmt.Println("res1[1]",res1[1])
	fmt.Println("res1[2]",res1[2])
	fmt.Println("res1[3]",res1[3])
	fmt.Println("res1[4]",res1[4])
	
	
	
	
	
	
	//fmt.Println("appEui Adil: ", appEui) 
	
	opts := MQTT.NewClientOptions().AddBroker(res1[0])
	opts.SetClientID("LoraDataPublisher")
	opts.SetDefaultPublishHandler(f)
	opts.SetUsername(res1[1])
	opts.SetPassword(res1[2])
	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	defer c.Disconnect(300)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Publishing")
	//Publish messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message

	
	
	
	//postdata := "{ \"d\" :" + jsonStr + "}"
	//fmt.Println("Data: ", postdata)
	
	if  dev1.RawData==false {
	postdata := string(databuffer)
	appEui:= postdata[71:94]
	topic := "/tplora/" +res1[1]+"/"+res1[3]+"/"+res1[4]+"/"+appEui+"/"+ dev1.Deviceeui + "/tpingress"
	token := c.Publish(topic, 0, false, postdata)
	//fmt.Println("Token --> ", token)
	token.Wait()
	//time.Sleep(3 * time.Second)

	c.Disconnect(250)
	} else	{
	appEui:= jsondata[71:94]
	topic := "/tplora/" +res1[1]+"/"+res1[3]+"/"+res1[4]+"/"+appEui+"/"+ dev1.Deviceeui + "/tpingress"
	token := c.Publish(topic, 0, false, jsondata)
	//fmt.Println("Token --> ", token)
	token.Wait()
	//time.Sleep(3 * time.Second)

	c.Disconnect(250)
	
	}
}
