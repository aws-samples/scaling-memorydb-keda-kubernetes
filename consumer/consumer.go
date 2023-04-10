package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
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

	// pod name
	workerName := os.Getenv("INSTANCE_NAME")

	go func() {
		fmt.Println("waiting for items...")
		for {
			items, err := client.BRPop(ctx, 0*time.Second, listName).Result()
			if err != nil {
				fmt.Println("unable to fetch item from list", err)
				continue
			}

			_, err = client.LPush(ctx, workerName, items[1]).Result()
			if err != nil {
				fmt.Println("unable to add to list", err)
				continue
			}

			fmt.Println("processed item", items[1])
			time.Sleep(time.Duration(rand.Intn(5)+2) * time.Second)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	fmt.Println("press ctrl+c to exit...")
	<-exit
	fmt.Println("program exited")
}
