// This is go-bindings package for librrd
package rrd

// #cgo LDFLAGS: -lrrd_th 
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include "rrd.h"
// #define HAVE_INTTYPES_H 1
// #include "rrd_client.h"
import "C"

import (
	"errors"
	"unsafe"
)

// int rrdc_connect (const char *addr);
func RrdcConnect(address string) (err error) {
	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))
	ret := C.rrdc_connect(caddr)
	if ret == -1 {
		err = errors.New(getError())
	}
	return
}

// int rrdc_is_connected(const char *daemon_addr);
func RrdcIsConnected(address string) (connected bool, err error) {
	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))
	ret := C.rrdc_is_connected(caddr)
	switch ret {
	case -1:
		{
			err = errors.New(getError())
		}
	case 1:
		{
			connected = true
		}
	default:
		{
			connected = false
		}
	}
	return
}

// int rrdc_disconnect (void);
func RrdcDisconnect() (err error) {
	ret := C.rrdc_disconnect()
	if ret == -1 {
		err = errors.New(getError())
	}
	return
}

// int rrdc_update (const char *filename, int values_num,
//         const char * const *values);

// int rrdc_flush (const char *filename);
func RrdcFlush(filename string) (err error) {
	cfn := C.CString(filename)
	defer C.free(unsafe.Pointer(cfn))
	ret := C.rrdc_flush(cfn)
	if ret == -1 {
		err = errors.New(getError())
	}
	return
}

// int rrdc_flush_if_daemon (const char *opt_daemon, const char *filename);
/*
func RrdcFlushIfDaemon(optDaemon string, filename string) (err error) {
	cod := C.CString(optDaemon)
	defer C.free(unsafe.Pointer(cod))

	cfn := C.CString(filename)
	defer C.free(unsafe.Pointer(cfn))

	ret := C.rrdc_flush_if_daemon(cod, cfn)
	if ret == -1 {
		err = errors.New(getError())
	}
	return
}
*/
