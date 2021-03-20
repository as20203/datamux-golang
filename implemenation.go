package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/naoina/toml"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Global structure for configuration data
type tomlConfig struct {
	Listeningport uint64
	Devices       []device
}

//Global structure for configuration data
var config tomlConfig

//File pointers
// var fpOnYieldData *os.File

func readConfig() {

	f, err := os.Open("config.dat")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	f.Close()
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
}

func printConfig() {
	fmt.Println("Listening on Port: ", config.Listeningport)
	fmt.Println("Devices included:")
	for _, elem := range config.Devices {
		fmt.Println("DeviceEUI:", elem.Deviceeui, " Devicetype:", elem.Devicetype)
	}
	fmt.Println(" ================================================================= ")
	fmt.Println("  ")
}

func findDeviceEUI(deviceeui string) (device, error) {
	uri := "mongodb://localhost:27017/datamux"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	var retrievedDevice device
	datamuxDatabase := client.Database("datamux")
	devicesCollection := datamuxDatabase.Collection("devices")

	if err = devicesCollection.FindOne(ctx, bson.M{"deviceeui": deviceeui}).Decode(&retrievedDevice); err != nil {
		return device{}, errors.New("Unknown Device: " + deviceeui)
	} else {
		return retrievedDevice, nil
	}
}

// func createOutputFiles() {
// 	Filename := "OnYieldData" + time.Now().Format("_02_01_15_04_05") + ".csv"
// 	var err error
// 	fpOnYieldData, err = os.Create(Filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	w := csv.NewWriter(fpOnYieldData)
// 	err = w.Write([]string{"Time", "DevEUI", "Temperature", "Humidity"})
// 	if err != nil {
// 		panic(err)
// 	}
// 	w.Flush()
// 	if err = w.Error(); err != nil {
// 		log.Fatal(err)
// 	}
// 	//defer fpOnYieldData.Close()
// }

// func storeOnYieldData(deviceEUI string, parsedvalues []onyielddata) bool {

// 	for _, v := range parsedvalues {
// 		fmt.Println(v.Time.Format(time.StampMilli), deviceEUI, fmt.Sprint(v.Temperature), fmt.Sprint(v.RelativeHumidity), fmt.Sprint(v.AbsoluteHumidity))
// 		w := csv.NewWriter(fpOnYieldData)
// 		err := w.Write([]string{v.Time.Format(time.StampMilli), deviceEUI, fmt.Sprint(v.Temperature), fmt.Sprint(v.RelativeHumidity), fmt.Sprint(v.AbsoluteHumidity)})
// 		if err != nil {
// 			panic(err)
// 		}
// 		w.Flush()
// 		if err := w.Error(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	return true
// }

/*Data Generator for ESAB*/
// func decodeESABData() []byte {
// 	return nil
// }

// func decodeTrackerData() []byte {
// 	return nil
// }

/* Generic function for decoding data
Add a case for each type of device and add a function of its own
so that we can support multiple devices.

If no device type is specified, then 313233343536 is returned as data*/

func decodeDataDownlink(dev device, entry downlinkdata) bool {

	fmt.Println("downlink correct place")

	publishDownlinkData(dev, entry)
	return true

}
func decodeData(dev device, entry loradata) bool {
	switch strings.ToLower(dev.Devicetype) {
	case "oy1100":
		fmt.Println("OY1100")
		parsedvalues := parseOnYieldData(extractTimewithTimeZone(entry.Time), entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}
		//storeOnYieldData(deviceEUI, parsedvalues)
		//Same data format for both OY1100 and OY1110. Only parsing changes. That is already done.
		publishOY1110Data(dev, entry, parsedvalues)
		return true
	case "oy1110":
		fmt.Println("OY1110")
		parsedvalues := parseOY1110Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1110Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1700":
		fmt.Println("OY1700")
		parsedvalues := parseOY1700Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1700Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1200":
		fmt.Println("OY1200")
		parsedvalues := parseOY1200Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1200Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1700v1":
		fmt.Println("OY1700V1")
		parsedvalues := parseOY1700V1Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishOY1700V1Data(dev, entry, parsedvalues)

		return true
	case "peoplecounter":
		fmt.Println("PeopleCounter")
		parsedvalues := parsePeopleCounterData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishPeopleCounterData(dev, entry, parsedvalues)

		return true

	case "tetraedretillquistdiz":
		fmt.Println("tetraedretillquistdiz")

		parsedvalues := parseTETRAEDRETillquistDIZDataData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTETRAEDRETillquistDIZData(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "tetraedre_abb_b24":
		fmt.Println("tetraedre_abb_b24")

		parsedvalues := parseTETRAEDRE_ABB_B24DataData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTETRAEDRE_ABB_B24Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "tetraedregwfmtk3":
		fmt.Println("tetraedregwfmtk3")

		parsedvalues := parsedTetraedreGWFMTK3Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTetraedreGWFMTK3Data(dev, entry, parsedvalues)

		return true

	case "adeunius_modbus_ecs_link":
		fmt.Println("adeunius_modbus_ecs_link")

		parsedvalues := parsedAdeunius_ModBus_ECS_LinkData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishAdeunius_ModBus_ECS_LinkData(dev, entry, parsedvalues)

		return true
	case "tetraedre_cyble_2":
		fmt.Println("Tetraedre_Cyble_2")

		parsedvalues := parsedTetraedre_Cyble_2Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTetraedre_Cyble_2Data(dev, entry, parsedvalues)

		return true

	case "tetraedrelogt550":
		fmt.Println("tetraedrelogt550")

		parsedvalues := parsedTetraedreLoGT550Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTetraedreLoGT550Data(dev, entry, parsedvalues)

		return true
	case "tetraedrearmatec":
		fmt.Println("TetraedreArmaTec")

		parsedvalues := parsedTetraedreArmaTecData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTetraedreArmaTecData(dev, entry, parsedvalues)

		return true

	case "cubecellconductivity":
		fmt.Println("CubeCellConductivity")

		parsedvalues := parsedCubeCellConductivityData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishCubeCellConductivityData(dev, entry, parsedvalues)

		return true
	case "tetraedraedcba":
		fmt.Println("TetraedraEDCBA")

		parsedvalues := parseTetraedraEDCBADataData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTetraedraEDCBAData(dev, entry, parsedvalues)

		return true
	case "oy1210":
		fmt.Println("OY1210")
		parsedvalues := parseOY1210Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1210Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1320":
		fmt.Println("OY1320")
		/*
			TODO:
			parsedvalues := parseOY1320Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
			if len(parsedvalues) == 0 {
				return false
			}
		*/
		var parsedvalues []OY1320Data
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1320Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "oy1320v1":
		fmt.Println("OY1320VI")

		parsedvalues := parseOY1320V1Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//var parsedvalues []OY1320VIData
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1320V1Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "wateriwmlr3":
		fmt.Println("wateriwmlr3")

		parsedvalues := parseWaterIWMLR3Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//var parsedvalues []OY1320VIData
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishWaterIWMLR3Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "digimondometer":
		fmt.Println("DigimondoMeter")

		parsedvalues := parsedDigimondoMeterData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//var parsedvalues []OY1320VIData
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishDigimondoMeterData(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "oy1400":
		fmt.Println("OY1400")
		parsedvalues := parseOY1400Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)

		if len(parsedvalues) == 0 {
			fmt.Println("oops")
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1400Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1410":
		fmt.Println("OY1410")
		parsedvalues := parseOY1410Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)

		if len(parsedvalues) == 0 {
			fmt.Println("oops")
			return false
		}

		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1410Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "oy1600v1":
		fmt.Println("OY1600V1")

		parsedvalues := parseOY1600V1Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//var parsedvalues []OY1600V1Data
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1600V1Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "unknown":
		fmt.Println("unknown")

		parsedvalues := parseUNKNOWNData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishUNKNOWNData(dev, entry, parsedvalues)

		return true
	case "smartmc0101door":
		fmt.Println("smartmc0101door")
		parsedvalues := parsedSmartMC0101DoorData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}
		//var parsedvalues []OY1600Data
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishSmartMC0101DoorData(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true

	case "oy1600":
		fmt.Println("OY1600")
		parsedvalues := parseOY1600Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}
		//var parsedvalues []OY1600Data
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishOY1600Data(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "landisgyr":
		fmt.Println("LandisGyr")

		parsedvalues := parseLandisGyrData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		//var parsedvalues []OY1320VIData
		//Irrespective of the publish result, we must respond SUCCESS to the network server
		publishLandisGyrData(dev, entry, parsedvalues)
		//storeOnYieldData(deviceEUI, parsedvalues)
		return true
	case "smartvalve":
		fmt.Println("SmartValve")

		parsedvalues := parseSmartValveData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishSmartValveData(dev, entry, parsedvalues)

		return true

	case "lr210":
		fmt.Println("LR210")

		parsedvalues := parseLR210Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishLR210Data(dev, entry, parsedvalues)

		return true

	case "tds_100f_flow":
		fmt.Println("TDS_100F_Flow")
		parsedvalues := parsedTDS_100F_FlowData(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishTDS_100F_FlowData(dev, entry, parsedvalues)

		return true
	case "honeywell_ew773":
		fmt.Println("honeywell_ew773")
		parsedvalues := parsedhoneywellew773Data(extractTimewithTimeZone(entry.Time), entry.Port, entry.Data)
		if len(parsedvalues) == 0 {
			return false
		}

		publishHoneywell_EW773Data(dev, entry, parsedvalues)

		return true

	case "esab":
		fmt.Println("ESAB")
		//return decodeESABData()
		return false
	case "tracker":
		fmt.Println("Tracker")
		//return decodeTrackerData()
		return false

	case "downlink":

		fmt.Println("downlink Wrong area")
		return true

	default:
		return false
	}
}

func printIncomingData(entry loradata) {
	fmt.Println("Received:", entry.Time, entry.DeviceEui, entry.Data)
}
func printIncomingDownlinkData(entry downlinkdata) {
	fmt.Println("Received:", entry.DeviceEui, entry.value, entry.Command, entry.AccessToken, entry.Devicetype)
}
func extractTimewithTimeZone(timestring string) time.Time {
	//timestringwithtimezone := strings.Replace(timestring, " ", "T", 1) + "+00:00"
	var timeentry time.Time
	//timeentry.UnmarshalText([]byte(timestringwithtimezone))
	timeentry.UnmarshalText([]byte(timestring))
	return timeentry
}

//Add a new key entry
func processLoraData(entry loradata) bool {

	//Print the incoming data
	printIncomingData(entry)

	//Determine the device type
	//if it is undefined, throw an error
	dev, err := findDeviceEUI(entry.DeviceEui)

	if err != nil {
		return false
	}

	//Decode the data
	return decodeData(dev, entry)
}

func processDownlinkData(entry downlinkdata) bool {

	//Print the incoming data
	printIncomingDownlinkData(entry)

	//DetermindecodeDataDownlinke the device type
	//if it is undefined, throw an error
	dev, err := findDeviceEUI(entry.DeviceEui)

	if err != nil {
		return false
	}

	return decodeDataDownlink(dev, entry)
}
