package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

var rdb *redis.Client
var m3 = make([]Message, 1)

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Subscribe() {
	var m1 Message
	pubsub := rdb.PSubscribe("*")
	_, err := pubsub.Receive()
	if err != nil {
		return
	}
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload, "\r\n")
		m1.Channel = msg.Channel
		m1.Msg = msg.Payload
		m3 = append(m3, m1)
	}
	m3 = m3[1:]
}

func Checkchannle(channel string) []Message {
	m2 := make([]Message, 1)
	fmt.Println(len(m3))
	for i := 0; i < len(m3); i++ {
		if channel == m3[i].Channel {
			fmt.Println(m3[i].Channel)
			m2 = append(m2, m3[i])
			fmt.Println(true)
			fmt.Println(m3[i])
		}
	}
	m2 = m2[1:]
	return m2
}

func Publish(channel string, msg string) {
	err := rdb.Publish(channel, msg)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	r := gin.Default()
	initClient()
	go Subscribe()
	r.POST("/publish", Publishweb)
	r.POST("/message", Messageweb)
	r.Run(":6060")
}

func Publishweb(c *gin.Context) {
	user := c.PostForm("user")
	msg := c.PostForm("msg")
	Publish(user, msg)
	c.JSON(http.StatusOK, gin.H{
		"state": "true",
	})
}

func Messageweb(c *gin.Context) {
	m4 := make([]Message, 1)
	channel := c.PostForm("channnel")
	m4 = Checkchannle(channel)
	c.JSON(200, m4)
}

type Message struct {
	Channel string `json:"channel"`
	Msg     string `json:"msg"`
}
