package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	endpoint := os.Getenv("MEMORYDB_ENDPOINT")
	if endpoint == "" {
		log.Fatal("MEMORYDB_ENDPOINT env variable missing")
	}

	username := os.Getenv("MEMORYDB_USERNAME")
	if username == "" {
		log.Fatal("MEMORYDB_USERNAME env variable missing")
	}

	password := os.Getenv("MEMORYDB_PASSWORD")
	if password == "" {
		log.Fatal("MEMORYDB_PASSWORD env variable missing")
	}

	listName := os.Getenv("LIST_NAME")
	if listName == "" {
		log.Fatal("LIST_NAME env variable missing")
	}

	fmt.Println("connecting to Redis host ", endpoint)

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:     []string{endpoint},
		Username:  username,
		Password:  password,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()

	if err != nil {
		log.Fatal("failed to connect to redis", err)
		return
	}
	fmt.Println("successfully connected to redis", endpoint)

	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("failed to close client conn ", err)
			return
		}
		fmt.Println("closed client connection")

	}()

	go func() {
		pipe := client.Pipeline()
		num := 50

		for {
			log.Println("sending", num, "items into redis list", listName)

			for i := 1; i <= num; i++ {
				err := pipe.LPush(ctx, listName, "message-"+strconv.Itoa(rand.Intn(50))).Err()
				if err != nil {
					fmt.Println("unable to pipe data to redis list", err)
					continue
				}
			}

			_, err = pipe.Exec(ctx)
			if err != nil {
				fmt.Println("unable to exec pipeline to send data to redis list", err)
			}

			time.Sleep(2 * time.Second)

			//time.Sleep(time.Duration(rand.Intn(3)+2) * time.Second)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	fmt.Println("press ctrl+c to exit...")
	<-exit
	fmt.Println("program exited")

}
