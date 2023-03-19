// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"trailers/server/src"

	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var config src.Configuration

func readConfiguration(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())

	var trailerDataMqtt src.TrailerDataMqtt
	err := json.Unmarshal(msg.Payload(), &trailerDataMqtt)
	if err != nil {
		print("Invalid trailers data")
		fmt.Println(err)
		return
	}
	//fmt.Printf("%f %f\n", trailerData.Latt, trailerData.Longt)

	db := src.ConnectToDb(config.DbConfig)
	defer db.Close()

	var trailerData src.TrailerData
	//latt, err := strconv.ParseFloat(trailerDataMqtt.Latt, 64)
	//if err != nil {
	//	print("Cannot parse latt")
	//	return
	//}
	//longt, err := strconv.ParseFloat(trailerDataMqtt.Longt, 64)
	//if err != nil {
	//	print("Cannot parse longt")
	//	return
	//}
	//trailerData.Latt = latt
	//trailerData.Longt = longt

	trailerData.Latt = trailerDataMqtt.Latt
	trailerData.Longt = trailerDataMqtt.Longt

	if trailerDataMqtt.GpsTime == " : : " {
		trailerData.GpsTime = time.Unix(0, 0)
	} else {
		trailerData.GpsTime, err = time.Parse("15:04:05", trailerDataMqtt.GpsTime)
	}
	if err != nil {
		fmt.Println(err)
	}
	//layout := "2022-10-22T17:48:22.592467"
	trailerData.OsTime, err = time.Parse(time.RFC3339, strings.Split(trailerDataMqtt.OsTime, "Z")[0]+"Z")
	//trailerData.OsTime, err = time.Parse(time.RFC3339, trailerDataMqtt.OsTime)
	if err != nil {
		fmt.Println(err)
	}
	trailerData.Weight = trailerDataMqtt.Weight
	trailerData.WeightStatus = trailerDataMqtt.WeightStatus
	trailerData.ShuntVoltage = trailerDataMqtt.ShuntVoltage
	trailerData.PowerSupplyVoltage = trailerDataMqtt.PowerSupplyVoltage
	trailerData.SerialNumber = trailerDataMqtt.SerialNumber
	trailerData.CpuTemp = trailerDataMqtt.CpuTemp

	//src.SendMqttMessageToDb(db, string(msg.Payload()), string(msg.Topic()))
	src.InsertTrailerDataIntoDb(db, trailerData)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}

func main() {
	readConfiguration("configuration/config.yaml")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
	config.MailApiConfig.SendGridApiKey = "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"
	config.MailApiConfig.MailVerifCodeExpiration = 3
	config.MailApiConfig.PassResetCodeExpiration = 30
	config.MailApiConfig.MailVerifTemplateID = "d-765c9b3176b940e0bafee768b5d44124"
	config.MailApiConfig.PassResetTemplateID = "d-8520acc570d64a5686e6fa8ef40ff2cd"
	//mail = src.Mail()

	var broker = "tcp://192.168.1.50:1883"
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("go_mqtt_example")
	options.SetUsername("admin")
	options.SetPassword("panchoididi")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "GARY"
	token = client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)

	// num := 3
	// for i := 0; i < num; i++ {
	// 	text := fmt.Sprintf("%d", i)
	// 	token = client.Publish(topic, 0, false, text)
	// 	token.Wait()
	// 	time.Sleep(time.Second)
	// }
	//time.Sleep(100000 * time.Second)
	src.SetupRoutes(config)

	//client.Disconnect(100)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	fmt.Fprintf(w, "Hello World!")
}
