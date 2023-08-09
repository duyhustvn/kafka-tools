package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(ctx context.Context, brokerUrl []string, topic string, data string, idx int, successChan chan int, failChan chan int) {
	log.Println("Start sending message ", idx)
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokerUrl...),
		Topic: topic,
	}

	if err := w.WriteMessages(ctx, kafka.Message{
		Value: []byte(data),
	}); err != nil {
		log.Printf("There is error produce message %+v", err)
		failChan <- idx
	}

	successChan <- idx
}

func StringPrompt(label string) string {
	fmt.Println(label)
	reader := bufio.NewReader(os.Stdin)

	lines := ""
	for {
		// read line from stdin using newline as separator
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// if line is empty, break the loop
		if len(strings.TrimSpace(line)) == 0 {
			break
		}

		//append the line to a slice
		lines += line

		// Break until \n
		break
	}
	lines = strings.TrimSpace(lines)
	return lines
}

func worker(ctx context.Context, workers chan int, brokers []string, topic string, data string, successChan chan int, failChan chan int) {
	for idx := range workers {
		ProduceMessage(ctx, brokers, topic, data, idx, successChan, failChan)
	}
}

// produce n messages
func ProduceMessages(ctx context.Context, brokers []string, topic string, message string, n int) {
	workers := make(chan int, 1000)

	successChan := make(chan int)
	failChan := make(chan int)

	for i := 0; i < cap(workers); i++ {
		go worker(ctx, workers, brokers, topic, message, successChan, failChan)
	}

	go func() {
		for i := 0; i <= n; i++ {
			workers <- i
		}
	}()

	successCounter := 0
	failedCounter := 0
	for i := 0; i < n; i++ {
		select {
		case <-successChan:
			successCounter++
		case <-failChan:
			failedCounter++
		}
	}

	fmt.Println("Success: ", successCounter)
	fmt.Println("Fail: ", failedCounter)

}

func main() {
	brokerUrl := StringPrompt("Enter broker url (Left empty for default value: 172.17.0.1:9092) >>")
	if brokerUrl == "" {
		brokerUrl = "172.17.0.1:9092"
	}
	fmt.Println("brokerUrl: ", brokerUrl)

	topic := StringPrompt("Enter the topic name >>")
	for topic == "" {
		topic = StringPrompt("Enter the topic name >>")
	}
	fmt.Println("topic: ", topic)

	produceMessage := ""
	for produceMessage == "" {
		produceMessage = StringPrompt("Enter produce Message >>")
		if produceMessage != "" {
			invalidNumber := true
			for invalidNumber {
				message_nums_str := StringPrompt("Enter the number of messages >>")
				message_num, err := strconv.Atoi(message_nums_str)
				if err != nil {
					fmt.Println("Invalid number")
				}
				ProduceMessages(context.Background(), []string{brokerUrl}, topic, produceMessage, message_num)
				invalidNumber = false
			}
		}
		produceMessage = ""
	}
}
