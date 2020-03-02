package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"fmt"
	"strings"
)

// doorFromMessage returns the name of the door that a message refers to.
func doorFromMessage(msg mqtt.Message) (lock string) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) >= 3 {
		lock = parts[2]
	}
	return
}

func handlePin(client mqtt.Client, message mqtt.Message) {
	door := doorFromMessage(message)
	pin := string(message.Payload())
	status := "ACCESS DENIED"

	if pinValidate(pin, door) {
		status = "ACCESS GRANTED"
		if token := client.Publish(fmt.Sprintf("bast/%s/%s/granted", name, door), 0, false, "granted"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
	} else {
		if token := client.Publish(fmt.Sprintf("bast/%s/%s/denied", name, door), 0, false, "denied"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
	}

	log.Printf("%s for PIN %s at %s\n", status, pin, door)
}

func handleCard(client mqtt.Client, message mqtt.Message) {
	door := doorFromMessage(message)
	card := string(message.Payload())
	status := "ACCESS DENIED"

	if cardValidate(card, door) {
		status = "ACCESS GRANTED"
		if token := client.Publish(fmt.Sprintf("bast/%s/%s/granted", name, door), 0, false, "denied"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
	} else {
		if token := client.Publish(fmt.Sprintf("bast/%s/%s/denied", name, door), 0, false, "granted"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
	}


	log.Printf("%s for card %s at %s\n", status, card, door)
}

