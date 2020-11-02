package errorhandling

import "log"

//HandleError - handle error
func HandleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
