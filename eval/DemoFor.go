package main
import (
	"fmt"
	"time"
)
var flag bool
func main(){
	flag=true
	go runtrue()
	go runflase()
	time.Sleep(1*time.Second)
	flag=false
	time.Sleep(1*time.Second)
}

func runtrue() {

	for {
		if flag == true {
			fmt.Println("I am true flag")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func runflase() {
	for {
		if flag == false {
			fmt.Println("I am false flag")
			time.Sleep(100 * time.Millisecond)
		}
	}
}
