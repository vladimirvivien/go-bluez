package main

/*
#cgo LDFLAGS: -lbluetooth
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/socket.h>
#include <bluetooth/bluetooth.h>
#include <bluetooth/hci.h>
#include <bluetooth/hci_lib.h>
*/
import "C"

import (
	"fmt"
	"os"
)

func main() {
	devId := C.hci_get_route(nil)
	sock := C.hci_open_dev(devId)
	if devId < 0 || sock < 0 {
		fmt.Println("Error opening socket")
		os.Exit(1)
	}
	flags := C.IREQ_CACHE_FLUSH
	var ii *C.inquiry_info
	numRsp := C.hci_inquiry(devId, 8, 255, nil, &ii, C.long(flags))
	fmt.Println(numRsp)
}
