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
	"unsafe"
)

func main() {
	devId := C.hci_get_route(nil)
	sock := C.hci_open_dev(devId)
	defer C.close(sock)

	if devId < 0 || sock < 0 {
		fmt.Println("Error opening socket for Bluetooth connection.")
		fmt.Println("Ensure Bluetooth is enabled.")
		os.Exit(1)
	}
	fmt.Println("Searching for Bluetooth devices...")
	flags := C.IREQ_CACHE_FLUSH
	maxRsp := 255
	var ii *C.inquiry_info
	iisize := unsafe.Sizeof(ii)
	iiptr := (*C.inquiry_info)(C.malloc(C.size_t(maxRsp) * C.size_t(iisize)))
	defer C.free(unsafe.Pointer(iiptr))

	numRsp := int(C.hci_inquiry(devId, 8, 255, nil, &iiptr, C.long(flags)))
	if numRsp <= 0 {
		fmt.Println("Unable to find Bluetooth devices.")
		os.Exit(0)
	}

	fmt.Printf("Found %d Bluetooth device(s)\n", numRsp)
	var addr [19]C.char
	var devAddr string
	var name [248]C.char
	var devName string
	for i := 0; i < numRsp; i++ {
		ptrOffset := uintptr(unsafe.Pointer(iiptr)) + (uintptr(iisize) * uintptr(i))
		ii = (*C.inquiry_info)(unsafe.Pointer(ptrOffset))
		C.ba2str(&ii.bdaddr, &addr[0])
		devAddr = C.GoString(&addr[0])
		result := C.hci_read_remote_name(sock, &ii.bdaddr, C.int(len(name)), &name[0], 0)
		if result < 0 {
			devName = "UNKNOWN"
		} else {
			devName = C.GoString(&name[0])
		}
		fmt.Printf("Device %s - %s\n", devName, devAddr)
	}
}
