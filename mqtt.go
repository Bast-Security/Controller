package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
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
	}

	log.Printf("%s for PIN %s at %s\n", status, pin, door)
}

func handleCard(client mqtt.Client, message mqtt.Message) {
	door := doorFromMessage(message)
	card := string(message.Payload())
	status := "ACCESS DENIED"


	if cardValidate(card, door) {
		status = "ACCESS GRANTED"
	}


	log.Printf("%s for card %s at %s\n", status, card, door)
}

