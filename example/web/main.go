package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adiclepcea/go-omega2gpio"
)

//PWM is used for commanding pwm on gpio pins
type PWM struct {
	Gpio      int `json:"gpio"`
	HZ        int `json:"hz"`
	DutyCycle int `json:"duty"`
}

//Gpio is used for commanding on or off on gpio pins
type Gpio struct {
	Gpio int `json:"gpio"`
}

func pwm(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	pwm := PWM{}

	err := decoder.Decode(&pwm)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	go onion.StopPwm(pwm.Gpio)

	time.Sleep(100 * time.Millisecond)

	go onion.SPwm(pwm.Gpio, pwm.HZ, pwm.DutyCycle)

	log.Println(pwm)
}

func gpio(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	gpio := Gpio{}

	err := decoder.Decode(&gpio)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	onion.SetDirection(gpio.Gpio, 1)

	now := onion.Read(gpio.Gpio)

	if now == 0 {
		onion.Write(gpio.Gpio, uint8(1))
	} else {
		onion.Write(gpio.Gpio, uint8(0))
	}

	now = onion.Read(gpio.Gpio)

	w.Write([]byte(fmt.Sprintf("%d", now)))

	log.Println(gpio)
}

func main() {
	onion.Setup()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/pwm", pwm)
	http.HandleFunc("/gpio", gpio)

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
