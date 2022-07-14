// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (

	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"fmt"
	"trailers/server/src"

	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
	// 	db := src.ConnectToDb()
	// 	defer db.Close()
	// 	src.SendMqttMessageToDb(db, string(msg.Payload()), string(msg.Topic()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}

func main() {
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
	src.SetupRoutes()

	//client.Disconnect(100)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	fmt.Fprintf(w, "Hello World!")
}
