package main

import (
	"fmt"
	"os"
	"strconv"

	onion "github.com/adiclepcea/go-omega2gpio"
)

func showUsage() {
	fmt.Printf(`Usage
	%s set-input <gpio>
	%s set-output <gpio>
	%s get-direction <gpio>
	%s read <gpio>
	%s set <gpio> <value: 0 or 1>
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
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

	}

	/*status := int64(1)
	if len(os.Args) > 1 {
		status, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Printf("Error parsing arg %s\n", err.Error())
		}
		if status != 0 {
			status = 1
		}

	}

	onion.Setup()

	onion.SetDirection(18, 0)

	fmt.Println("OK")
	fmt.Printf("offsets: pin 18=%d, pin 31=%d, pin 32=%d", onion.GetDirection(18), onion.GetDirection(31), onion.GetDirection(32))
	if onion.GetDirection(18) != 1 {
		onion.SetDirection(18, 1)
	}
	onion.Write(18, uint8(status))
	log.Printf("Pin 18 has now value: %d\n", onion.Read(18))
	*/
}
