package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	onion "github.com/adiclepcea/go-omega2gpio"
)

func showUsage() {
	fmt.Printf(`Usage
	%s set-input <gpio>
	%s set-output <gpio>
	%s get-direction <gpio>
	%s read <gpio>
	%s set <gpio> <value: 0 or 1>
	%s pwm <gpio> <freq in HZ> <duty cycle percentage>
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
		return
	}

	pin, err := strconv.ParseUint(os.Args[2], 10, 64)

	if err != nil {
		showUsage()
		return
	}

	if pin > 46 {
		fmt.Println("Onion only exposes up to port 46")
		return
	}

	onion.Setup()
	switch os.Args[1] {
	case "set-input":
		onion.SetDirection(int(pin), 0)
		break
	case "set-output":
		onion.SetDirection(int(pin), 1)
		break
	case "get-direction":
		directionInt := onion.GetDirection(int(pin))
		directionStr := "output"
		if directionInt == 0 {
			directionStr = "input"
		}
		fmt.Printf("> Get direction GPIO%d: %s\n", pin, directionStr)
		break
	case "read":
		valueInt := onion.Read(int(pin))
		fmt.Printf("> Read GPIO%d: %d\n", pin, valueInt)
		break
	case "set":
		if len(os.Args) < 4 {
			showUsage()
			break
		}

		argInt, err := strconv.ParseUint(os.Args[3], 10, 64)

		if err != nil {
			showUsage()
			return
		}

		if argInt != 0 {
			argInt = 1
		}

		onion.Write(int(pin), uint8(argInt))
		break
	case "pwm":
		if len(os.Args) < 5 {
			showUsage()
			break
		}
		freqInt, err := strconv.ParseUint(os.Args[3], 10, 64)

		if err != nil {
			showUsage()
			return
		}

		dutyInt, err := strconv.ParseUint(os.Args[4], 10, 64)

		if err != nil || dutyInt > 100 {
			showUsage()
			return
		}

		go onion.SPwm(int(pin), int(freqInt), int(dutyInt))

		time.Sleep(3000 * time.Millisecond)

		go onion.StopPwm(int(pin))

		time.Sleep(3 * time.Second)

	}

}
