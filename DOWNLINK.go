package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"net/http"
	"strings"
)

type DownlinkData struct {
	Port      string  
	command	string
	
	
}

type decodeddownlinkdata struct {
	DeviceEui string       `json:"deviceEui"`
	
	Command    string       `json:"command"`
	AccessToken  string      `json:"access_token"`
	
}


func addZero(xpos int, sVal string) string {
var zeros string;
j:=len(sVal)

k:=xpos - j
 for i := 0; i < k ; i++ {
 {
zeros+="0"
 }
	
	
}
return zeros+sVal
}

func publishDownlinkData(dev device, entry downlinkdata ) error {

var command_val string
var iPort string	
fmt.Println("Adil you are here 1: ", dev.Devicetype)

	switch dev.Devicetype {
		//Download Messages for Smart Valve
		
		case "LR210":	
		
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		if entry.Command=="get_relay_state"{
		    command_val="0222"
			iPort="1" 
			
		}
		if entry.Command=="set_relay_state"{
		     Haxval:=strconv.FormatUint(entry.value, 16)
		    commandP1:=addZero(8,Haxval)
			command_val="0122"+commandP1
			iPort="1" 
			
		}
		if entry.Command=="set_both_relay_active"{
		    command_val="012200030003"
			iPort="1" 
			
		}
		if entry.Command=="set_both_relay_deactive"{
		    command_val="012200030000"
			iPort="1" 
			
		}
		if entry.Command=="set_second_relay_active"{
		    command_val="012200020002"
			iPort="1" 
			
		}
		
		
		if entry.Command=="set_interval"{
		    Haxval:=strconv.FormatUint(entry.value, 16)
		    commandP1:=addZero(4,Haxval)
			command_val="0123"+commandP1
			iPort="1" 
				
		}
		
		
		case "SmartValve":	
		if entry.Command=="open"{
		    command_val="1"
			iPort="1" 
			
		}
		if entry.Command=="close"{
		    command_val="1"
			iPort="1" 
			
		}
		if entry.Command=="class_a"{
		    command_val="00"
			iPort="9" 
			
		}
		if entry.Command=="class_c"{
		    command_val="01"
			iPort="9" 
			
		}
		if entry.Command=="sync_clock"{
		    command_val="1"
			iPort="13" 
			
		}
		
		
		case "OY1700","OY1700V1" :	
		if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		if entry.Command=="stabilization_time"{
		    command_val="0340"
			iPort="1" 
			
		}
		
		if entry.Command=="status"{
		    command_val="0220"
			iPort="1" 
			
		}
		if entry.Command=="set_interval"{
		    Haxval:=strconv.FormatUint(entry.value, 16)
		    commandP1:=addZero(4,Haxval)
			command_val="0123"+commandP1
			iPort="1" 
				
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
	    if entry.Command=="removed_alarm" {
		     command_val="012000"
		     iPort="1" 
				
		}
		//Download Messages for OY1700V1
	case "OY1400":
	
	
		if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		
		
		
	    if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0126"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="group_measurement"{
			if entry.value==2 {
				command_val="012702"
				iPort="1" 
			}
			if entry.value==1 {
				command_val="012701"
				iPort="1" 
			}
			
					
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
		
	
	
	
	case "OY1410":
	 if entry.Command=="reset_status" {
		     command_val="012D00"
		     iPort="1" 
				
		}
		if entry.Command=="set_counter" {
		      Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(8,Haxval)
			   command_val="012C"+commandP1
			   iPort="1"  
				
		}
			if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0126"+commandP1
			   iPort="1"  
					
		}
	
	  if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
	

	case "OY1320","OY1320V1","OY1310" :
        if entry.Command=="remove_alarm" {
		     command_val="012000"
		     iPort="1" 
				
		}
		if entry.Command=="set_meter_reading" {
		       Haxval:=strconv.FormatUint(entry.value, 16)
		       commandP1:=addZero(8,Haxval)
			   command_val="0121"+commandP1
			   iPort="1" 
				
		}
       
		if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0123"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="query_backflow_volume" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(8,Haxval)
			   command_val="0127"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
	
		if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		if entry.Command=="stabilization_time"{
		    command_val="0340"
			iPort="1" 
			
		}
		
		if entry.Command=="status"{
		    command_val="0220"
			iPort="1" 
			
		}
		if entry.Command=="set_q3maxflow"{
		    command_val="012B"
			iPort="1" 
			
		}
		
		
	case "OY1210" :

if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}

		
		if entry.Command=="status"{
		    command_val="0220"
			iPort="1" 
			
		}
			if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0130"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="tx_group_size"{
			if entry.value==3 {
				command_val="013103"
				iPort="1" 
			}
			if entry.value==1 {
				command_val="013101"
				iPort="1" 
			}
			if entry.value==2 {
				command_val="013102"
				iPort="1" 
			}
			if entry.value==4 {
				command_val="013104"
				iPort="1" 
			}
				if entry.value==5 {
				command_val="013105"
				iPort="1" 
			}
					
		}
		if entry.Command=="CO2_variation_threshold"{
			  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0132"+commandP1
			   iPort="1" 
			
					
		}
		
		if entry.Command=="CO2_absolute_threshold"{
			  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0132"+commandP1
			   iPort="1" 
			
					
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
		
		
		
		
		
		
	case "OY1200":	
	
	if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		if entry.Command=="stabilization_time"{
		    command_val="0340"
			iPort="1" 
			
		}
		
		if entry.Command=="status"{
		    command_val="0220"
			iPort="1" 
			
		}
			if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0123"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="tx_group_size"{
			if entry.value==3 {
				command_val="012303"
				iPort="1" 
			}
			if entry.value==1 {
				command_val="012301"
				iPort="1" 
			}
			if entry.value==2 {
				command_val="012302"
				iPort="1" 
			}
			if entry.value==4 {
				command_val="012304"
				iPort="1" 
			}
				if entry.value==5 {
				command_val="012305"
				iPort="1" 
			}
					
		}
		if entry.Command=="group_measurement"{
			if entry.value==3 {
				command_val="012403"
				iPort="1" 
			}
			if entry.value==1 {
				command_val="012401"
				iPort="1" 
			}
			
					
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
	
	case "OY1100" ,"OY1110":
	
		if entry.Command=="build_hash"{
		    command_val="0203"
			iPort="1" 
			
		}
		if entry.Command=="cpu_voltage"{
		    command_val="0206"
			iPort="1" 
			
		}
		if entry.Command=="cpu_temperature"{
		    command_val="020A"
			iPort="1" 
			
		}
		if entry.Command=="stabilization_time"{
		    command_val="0340"
			iPort="1" 
			
		}
		
		if entry.Command=="status"{
		    command_val="0220"
			iPort="1" 
			
		}
			if entry.Command=="set_interval" {
        	  Haxval:=strconv.FormatUint(entry.value, 16)
		      commandP1:=addZero(4,Haxval)
			   command_val="0123"+commandP1
			   iPort="1"  
					
		}
		if entry.Command=="group_measurement"{
			if entry.value==3 {
				command_val="012403"
				iPort="1" 
			}
			if entry.value==1 {
				command_val="012401"
				iPort="1" 
			}
			
					
		}
		if entry.Command=="device_reset"{
			
				command_val="0305"
				iPort="1" 
								
		}
		
		
       default :
		
		
	}

	
	
	
	
	fmt.Println("Adil you are here 2: ", dev.Devicetype)
  
	   
	   return loginToNSserver(entry.DeviceEui,iPort,command_val,entry.AccessToken,dev.AccessToken)
  

	// fmt.Println("Data sent: ", string(loradatabytes))
		// transferDatatoEndPoint(loradatabytes, dev)
	

}
func loginToNSserver( deviceEUI string, port string,dataval string,eAccessToken string,devAccessToken string) error   {
	
	jsonData := map[string]string{"username": "datamux", "password": "vg@#BJV$6RyMR_8x"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("https://ns.talkpool.com/login", "application/json", bytes.NewBuffer(jsonValue))
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
	var jsonStr1 = []byte("{\"deviceEui\":\""+deviceEUI+"\",\"type\" : \"data\",\"port\" : \"" + port+ "\",\"confirmedFrame\" : 1,\"data\" : \""+ dataval + "\"}")
	
	
	fmt.Println(bytes.NewBuffer(jsonStr1))
	var url2="https://ns.talkpool.com/internal/device/"+deviceEUI+"/downlink/data"
	fmt.Println(url2)

    req2, err := http.NewRequest("POST", url2 , bytes.NewBuffer(jsonStr1))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-AUTH-TOKEN", res3)
	client := &http.Client{}
	
	 for _, cookie := range response.Cookies() {
		req2.Header.Set("Cookie", cookie.String())
		fmt.Println(cookie.String())
	}
	 if err != nil {
        fmt.Println("Error1:", err)
    }
	 
    resp2, err := client.Do(req2)
	  if err != nil {
		fmt.Println("Error reading response. ", err)
	}
	
    defer resp2.Body.Close()
	fmt.Println("response Status adil:", resp2.Status)
	

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}

	fmt.Printf("this is body %s\n", body)
	
	return json.NewDecoder(resp2.Body).Decode("Message")
	
	
	// bodyr, err := ioutil.ReadAll(resp2.Body)
    //if err != nil {
     //  fmt.Println("Error3:", err)
   // }
//	fmt.Println(string(bodyr))
	
	
	
	
	
	
}
