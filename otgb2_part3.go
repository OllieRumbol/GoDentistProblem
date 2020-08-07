package main

import (
	"math/rand"
	"time"
)

func dentist(wait <-chan chan int, dent <-chan chan int) {
	for {
		select {
		//Check if any patients are waiting in the waiting room
		case waitChannel := <-wait:
			patientHandling(waitChannel)

		default:
			//If there are no patients sleep (block) on the dent channel waiting for a patient to wake the dentist up
			println("Dentist is sleeping.")
			waitChannel := <-dent
			patientHandling(waitChannel)
		}
	}
}

//Patient helper function to reduce code duplication and carry out treatments on patients
func patientHandling(patient chan int) {
	//Get id of patient
	id := <-patient
	println("Patient", id, "is having a treatment.")
	//Need to leave enough time for a patient to join the queue, can have less than 1 second but the dentist will sleep again because no one is in the waiting room
	time.Sleep(time.Duration(rand.Intn(1500-1000)+1000) * time.Millisecond)
	patient <- id
}

func assistant(hwait chan chan int, lwait <-chan chan int, wait chan<- chan int) {
	//Create timer
	timer := time.NewTimer(500 * time.Millisecond)
	for {
		select {
		case highChannel := <-hwait:
			wait <- highChannel
			println("HELLO high")
			select {
			case <-timer.C:
				if len(lwait) > 0 {
					lowValue := <-lwait
					hwait <- lowValue
					//Reset timer
					timer = time.NewTimer(500 * time.Millisecond)
				}
			default:
			}
		case lowChannel := <-lwait:
			wait <- lowChannel
			//Reset timer
			timer = time.NewTimer(500 * time.Millisecond)
		default:
		}
	}
}

//Patient function
func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	println("Patient", id, "wants a treatment.")
	//Create channel for patient
	pat := make(chan int)

	select {
	//Trying communicating on dent channel, if blocked add patient to the waiting room
	case dent <- pat:
	default:
		//Dentist is busy with a patient so patient is added to the wait channel
		println("Patient", id, "is waiting.")
		wait <- pat
	}

	//Send id of patient
	pat <- id
	//Wait for id to return to signify that the dentist is finished
	<-pat
	println("Patient", id, "has shinny teeth.")
}

func main() {
	//Create four channels
	dent := make(chan chan int)
	wait := make(chan chan int, 10)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 5)

	go dentist(wait, dent)
	go assistant(hwait, lwait, wait)
	time.Sleep(3 * time.Second)

	high := 10
	low := 3
	for i := low; i < high; i++ {
		go patient(hwait, dent, i)
		time.Sleep(300 * time.Millisecond)
	}

	for i := 0; i < low; i++ {
		go patient(lwait, dent, i)
		time.Sleep(300 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}

/*
3b - deadlocks
In part 2 a dealock could occur when the timer has run out and I want to add someone from the low priority channel to the high priority channel however if there is no one in the low
priority channel then it would deadlock waiting for someone to be added t the channel. However this never occurs due to having an if statement to check the length of the channel.
I only pull from the low channel if the length of the channel is bigger than 0. This code can also be found is the assistent is part 3.
*/
