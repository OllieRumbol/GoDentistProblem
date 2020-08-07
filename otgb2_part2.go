package main

import (
	"math/rand"
	"time"
)

//Dentist function
func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	timer := time.NewTimer(1000 * time.Millisecond)
	for {
		select {
		//Check if any patients are in the high priority waiting room
		case highChannel := <-hwait:
			patientHandling(highChannel)
			select {
			//Checks if timer is finished
			case <-timer.C:
				//Checks if any patients are in the low priority waiting room
				if len(lwait) > 0 {
					lowValue := <-lwait
					//Add patient to high priority channel
					hwait <- lowValue
					//Reset timer
					timer = time.NewTimer(1000 * time.Millisecond)
				}
			default:
			}
		//Check if any patients are in the low priority waiting room
		case lowChannel := <-lwait:
			patientHandling(lowChannel)
			//Reset timer
			timer = time.NewTimer(1000 * time.Millisecond)
		//Sleep on the dent channel
		default:
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
	dent := make(chan chan int)
	hwait := make(chan chan int, 100)
	lwait := make(chan chan int, 5)
	go dentist(hwait, lwait, dent)
	time.Sleep(3 * time.Second)

	// high := 10
	low := 3
	// for i := low; i < high; i++ {
	// 	go patient(hwait, dent, i)
	// 	time.Sleep(300 * time.Millisecond)
	// }

	for i := 0; i < low; i++ {
		go patient(lwait, dent, i)
		time.Sleep(300 * time.Millisecond)
	}

	time.Sleep(50 * time.Second)
}

/*
2b - Starvation
One example of starvation that can be found is in the patient function when adding to the wait channel.
The wait channel is of a set size so only a limited number of patients can be stored.
Starvation can occur when trying to add to the wait channel because a patient could be blocked by other patients always trying to get access to the channel.
This problem could lead to a patient never being added to the waiting channel and therefore never seeing the dentist.

This can be fixed by using a timer and adding the longest waiting patient to the wait channel, this would then create fairness around patients.
*/
