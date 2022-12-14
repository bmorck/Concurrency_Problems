package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var partyAChan chan int
var partyACount = 0
var partyAMutex *sync.Mutex

var partyBChan chan int
var partyBCount = 0
var partyBMutex *sync.Mutex

func seatPartyA(wg *sync.WaitGroup, done chan int) {
	wg.Add(1)
	defer wg.Done()

	partyAMutex.Lock()
	if partyACount >= 3 {
		// release 3 party A threads
		for i := 0; i < 3; i++ {
			<-partyAChan
		}

		fmt.Println("Seating a member of party A")
		partyAMutex.Unlock()
		return
	}

	partyBMutex.Lock()
	if partyACount >= 1 && partyBCount >= 2 {
		<-partyAChan
		<-partyBChan
		<-partyBChan
		fmt.Println("Seating a member of party A")
		partyAMutex.Unlock()
		partyBMutex.Unlock()
		return
	}

	partyBMutex.Unlock()
	partyACount += 1
	partyAMutex.Unlock()

	// if the previous conditions are not met then we should block
	select {
	case partyAChan <- 0:
		fmt.Println("Seating a member of party A")
	case <-done:
		fmt.Println("Run out of members, party A member walking home")
	}

	partyAMutex.Lock()
	partyACount -= 1
	partyAMutex.Unlock()
}

func seatPartyB(wg *sync.WaitGroup, done chan int) {
	wg.Add(1)
	defer wg.Done()
	partyBMutex.Lock()
	if partyBCount >= 3 {
		// release 3 party B threads
		for i := 0; i < 3; i++ {
			<-partyBChan
		}
		fmt.Println("Seating a member of party B")
		partyBMutex.Unlock()
		return
	}

	partyAMutex.Lock()
	if partyBCount >= 1 && partyACount >= 2 {
		<-partyBChan
		<-partyAChan
		<-partyAChan
		fmt.Println("Seating a member of party B")
		partyBMutex.Unlock()
		partyAMutex.Unlock()
		return
	}

	partyAMutex.Unlock()
	partyBCount += 1
	partyBMutex.Unlock()

	// if the previous conditions are not met then we should block
	select {
	case partyBChan <- 0:
		fmt.Println("Seating a member of party B")
	case <-done:
		fmt.Println("Run out of members, party B member walking home")
	}

	partyBMutex.Lock()
	partyBCount -= 1
	partyBMutex.Unlock()

}

func main() {
	partyAChan = make(chan int, 1)
	partyBChan = make(chan int, 1)
	partyAMutex = &sync.Mutex{}
	partyBMutex = &sync.Mutex{}

	doneChan := make(chan int)

	// Ensure that both channels will block initially
	partyAChan <- 0
	partyBChan <- 0

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		n := rand.Intn(2)
		time.Sleep(time.Millisecond * 100)

		if n == 1 {
			fmt.Println("Party A Member is looking to be seated")
			go seatPartyA(wg, doneChan)
		} else {
			fmt.Println("Party B Member is looking to be seated")
			go seatPartyB(wg, doneChan)
		}

	}
	time.Sleep(time.Second)

	close(doneChan)

	wg.Wait()

	close(partyAChan)
	close(partyBChan)

}
