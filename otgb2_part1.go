package main

import (
	"math/rand"
	"time"
)

//Dentist function with parameters of read only channels
func dentist(wait <-chan chan int, dent <-chan chan int) {
	for {
		select {
		//Check if any patients are waiting in the waiting room/ channel
		case waitChannel := <-wait:
			patientHandling(waitChannel)
		default:
			//If there are no patients then sleep (block) on the dent channel waiting for a patient to wake the dentist up
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

//Patient function with parameters of write only channels
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

//Main function
func main() {
	//Set the size of the waiting room
	n := 5
	//Number of patients
	m := 5

	dent := make(chan chan int)    // creates a synchronous channel
	wait := make(chan chan int, n) // creates an asynchronous channel

	go dentist(wait, dent)
	time.Sleep(3 * time.Second)

	for i := 0; i < m; i++ {
		go patient(wait, dent, i)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1000 * time.Millisecond)
}
