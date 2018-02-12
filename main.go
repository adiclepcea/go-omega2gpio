package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

var (
	memlock   sync.Mutex
	base      int64
	memLength = 2048
	mmap      []uint32

	regBlockAddr int64 = 0x10000000
	regBlockSize       = 0x6AC

	//GPIO_CTRL_0 10000600(Directions for GPIO0-GPIO31)
	registerCtrlOffset = []int{384, 385, 386}
	//GPIO_CTRL_1 10000604(Directions for GPIO32-GPIO63)
	registerCtrl1Offset = 385
	//GPIO_CTRL_2 10000608(Directions for GPIO64-GPIO95)
	registerCtrl2Offset = 386

	//DATA REGISTERS: STATES OF GPIOS

	//GPIO_DATA_0 10000620(GPIO0-31)
	registerDataOffset = []int{392, 393, 394}
	//GPIO_DATA_1 10000624(GPIO32-63)
	registerData1Offset = 393
	//GPIO_DATA_2 10000628(GPIO64-95)
	registerData2Offset = 394

	//DATA SET REGISTERS: SET STATES OF GPIO_DATA_x registers

	//GPIO_DSET_0 10000630(GPIO0-31)
	registerDsetOffset = []int{396, 397, 398}
	//GPIO_DSET_1 10000634(GPIO31-63)
	registerDset1Offset = 397
	//GPIO_DSET_2 10000638(GPIO64-95)
	registerDset2Offset = 398

	//DATA CLEAR REGISTERS: CLEAR BITS OF GPIO_DATA_x registers

	//GPIO_DCLR_0 10000640(GPIO0-31)
	registerDclrOffset = []int{400, 401, 402}
	//GPIO_DCLR_1 10000644(GPIO31-63)
	registerDclr1Offset = 401
	//GPIO_DCLR_2 10000648(GPIO64-95)
	registerDclr2Offset = 402
)

func readRegistry(index int) uint32 {

	offset := registerCtrlOffset[index]

	regVal := uint32(mmap[offset+index])

	return regVal
}

func getDirection(pinNo int) uint32 {
	index := (pinNo) / 32

	regVal := readRegistry(index)

	gpio := uint32(pinNo % 32)

	val := ((regVal >> gpio) & 0x1)

	log.Printf("%s => %d, gpio=%d,pinNo=%d, byteVal=%d, index=%d\n", strconv.FormatInt(int64(regVal), 2), val, gpio, pinNo, regVal, index)

	return val

}

func setDirection(pinNo int, val uint8) {

	index := (pinNo) / 32

	regVal := readRegistry(index)
	gpio := uint32(pinNo % 32)

	if val == 1 {
		regVal |= (1 << gpio)
	} else {
		regVal &= ^(1 << gpio)
	}

	offset := registerCtrlOffset[index]

	memlock.Lock()

	defer memlock.Unlock()

	mmap[offset] = regVal
}

func write(pinNo int, val uint8) {

	var offset int
	gpio := uint32(pinNo % 32)
	index := (pinNo) / 32

	if val == 0 {
		offset = registerDclrOffset[index]
	} else {
		offset = registerDsetOffset[index]
	}

	regVal := (uint32(1) << gpio)

	log.Printf("reg=%d\n", regVal)

	memlock.Lock()
	defer memlock.Unlock()

	mmap[offset] = regVal
}

func read(pinNo int) uint32 {
	var offset int
	gpio := uint32(pinNo % 32)
	index := (pinNo) / 32

	offset = registerDataOffset[index]

	memlock.Lock()

	defer memlock.Unlock()

	regVal := uint32(mmap[offset+index])

	return ((regVal >> gpio) & 0x1)

}

func setup() {
	mfd, err := os.OpenFile("/dev/mem", os.O_RDWR, 0)

	if err != nil {
		log.Panic(err)
	}

	defer mfd.Close()

	memlock.Lock()

	defer memlock.Unlock()

	mmap8, err := syscall.Mmap(int(mfd.Fd()), regBlockAddr, memLength, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_FILE|syscall.MAP_SHARED)

	if err != nil {
		log.Panicf("Error mapping: %s\n", err.Error())
	}

	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mmap8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)

	mmap = *(*[]uint32)(unsafe.Pointer(&header))
}

func main() {
	var err error
	status := int64(1)
	if len(os.Args) > 1 {
		status, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Printf("Error parsing arg %s\n", err.Error())
		}
		if status != 0 {
			status = 1
		}

	}

	setup()

	setDirection(18, 0)

	fmt.Println("OK")
	fmt.Printf("offsets: pin 18=%d, pin 31=%d, pin 32=%d", getDirection(18), getDirection(31), getDirection(32))
	if getDirection(18) != 1 {
		setDirection(18, 1)
	}
	write(18, uint8(status))
	log.Printf("Pin 18 has now value: %d\n", read(18))

}
