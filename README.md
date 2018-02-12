# Go-omega2gpio

This is a go library that exposes the GPIO ports of the Onion Omega2.

The inspiration for this is from Onion fast-gpio https://github.com/OnionIoT/fast-gpio.

## Usage

First you have to obtain the library:

```
go get gihub.com/adiclepcea/go-omega2gpio
```

That will put the library on your GOPATH.

A simple example usage is:

```
package main

import (
    "fmt"

    onion "github.com/adiclepcea/go-omega2gpio"
)

func main(){
    onion.Setup()
    fmt.Println(onion.GetDirection(18))
}
```

This would show you the direction of the pin 18 (__0__ for __input__ and __1__ for __output__)

You have the following methods:

- __Setup()__ - this needs to be called once in every program that uses this library. This will map the memory for GPIO.
- __Read(pinNo int) uint32__ - will read the value of the pin __pinNo__ (0 or 1).
- __Write(pinNo int, val uint8)__ - will set the value of the pin __pinNo__ to the value __val__ (0 or 1).
- __SetDirection(pinNo int, val uint8)__ - will set the direction of the pin __pinNo__ to the __val__ (0 for __input__ or 1 for __output__) direction
- __func GetDirection(pinNo int) uint32__ - will return the direction (0 for __input__ or 1 for __output__) of the pin __pinNo__


## Compilation

Please remember to compile the program using this library for the __mips__ architecture.

```
GOOS=linux GOARCH=mipsle go build -o myprogram main.go
```

This will compile the program that you can just copy to your Onion Omega2 and run.

## Example

There is one example in the _example_ directory. This is a partial reimplementation of the __fast-gpio__ executable.

The present library does not implement pwm, but this should be trivial.

To compile you can just cd to the example directory and run ```make```

The directory should now contain an executable called ```go-fast-gpio```. You can copy it to your Onion Omega2:

```
scp go-fast-gpio root@<ipaddress_of_omega2>:/root/
```

You will asked for the password. After you write it and hit enter, the executable will be copied on your Omega2.
