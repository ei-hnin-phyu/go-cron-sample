package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"

	redislock "github.com/go-co-op/gocron-redis-lock"
)

func main() {
	redisOptions := &redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
	}
	redisClient := redis.NewUniversalClient(redisOptions)
	// Create a new scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	locker, err := redislock.NewRedisLocker(redisClient, redislock.WithTries(3))

	// Initialize a distributed locker
	// Set the locker for the scheduler
	scheduler.WithDistributedLocker(locker)

	// Schedule a job to run every 5 seconds
	scheduler.Every(5).Seconds().Do(func() {
		fmt.Println("Running job...")
	})

	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return
	}

	// Start the scheduler
	scheduler.StartAsync()
	<-make(chan chan struct{})

	// Wait for a while to allow the job to run
	time.Sleep(2 * time.Second)

	// Remove the job from the scheduler

	// Stop the scheduler
	scheduler.Stop()
}
