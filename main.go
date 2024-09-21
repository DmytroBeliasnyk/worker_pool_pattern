package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	startTime := time.Now()

	users := make(chan User, 1000)
	go generateUsers(1000, users)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go saveUserInfo(users, wg, runtime.NumCPU())

	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(users <-chan User, wg *sync.WaitGroup, workers int) {
	workerWg := &sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		workerWg.Add(1)
		go func() {
			for user := range users {
				fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

				filename := fmt.Sprintf("users/uid%d.txt", user.id)
				file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}

				file.WriteString(user.getActivityInfo())
				time.Sleep(time.Millisecond * 10)

				file.Close()
			}

			workerWg.Done()
		}()
	}

	workerWg.Wait()
	wg.Done()
}

func generateUsers(count int, users chan<- User) {
	defer close(users)
	for i := 0; i < count; i++ {
		users <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}

		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 10)
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
