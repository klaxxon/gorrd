// This is go-bindings package for librrd
package rrd

// #cgo LDFLAGS: -lrrd
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "rrd.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

const (
	CF_AVERAGE     = 0
	CF_MINIMUM     = 1
	CF_MAXIMUM     = 2
	CF_LAST        = 3
	CF_HWPREDICT   = 4
	CF_SEASONAL    = 5
	CF_DEVPREDICT  = 6
	CF_DEVSEASONAL = 7
	CF_FAILURES    = 8
	CF_MHWPREDICT  = 9

	DS_GAUGE    = 0
	DS_COUNTER  = 1
	DS_DERIVE   = 2
	DS_ABSOLUTE = 3
	DS_COMPUTE  = 4
)

// Convenience type definition for DS types
type DsType int

type RrdValue struct {
	Time  time.Time
	Value int64
}

func (this RrdValue) ToString() string {
	return fmt.Sprintf("%d:%d", this.Time.Unix(), this.Value)
}

func CfToString(cf int) string {
	switch cf {
	case CF_AVERAGE:
		return "AVERAGE"
	case CF_MINIMUM:
		return "MIN"
	case CF_MAXIMUM:
		return "MAX"
	case CF_LAST:
		return "LAST"
	case CF_HWPREDICT:
		return "HWPREDICT"
	case CF_SEASONAL:
		return "SEASONAL"
	case CF_DEVPREDICT:
		return "DEVPREDICT"
	case CF_DEVSEASONAL:
		return "DEVSEASONAL"
	case CF_FAILURES:
		return "FAILURES"
	case CF_MHWPREDICT:
		return "MHWPREDICT"

	default:
		return ""
	}
	return ""
}

// The Create function lets you set up new Round Robin Database (RRD) files.
// The file is created at its final, full size and filled with *UNKNOWN* data.
//
//      filename::
//          The name of the RRD you want to create. RRD files should end with the
//          extension .rrd. However, it accept any filename.
//      step::
//          Specifies the base interval in seconds with which data will be
//          fed into the RRD.
//      startTime::
//          Specifies the time in seconds since 1970-01-01 UTC when the first
//          value should be added to the RRD. It will not accept any data timed
//          before or at the time specified.
//      values::
//          A list of strings identifying datasources (in format "DS:ds-name:DST:dst arguments")
//          and round robin archives - RRA (in format "RRA:CF:cf arguments").
//          There should be at least one DS and RRA.
//
// See http://oss.oetiker.ch/rrdtool/doc/rrdcreate.en.html for detauls.
//
func Create(filename string, step uint64, startTime time.Time, values []string) (err error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	cvalues := makeCStringArray(values)
	defer freeCStringArray(cvalues)

	clearError()
	ret := C.rrd_create_r(cfilename, C.ulong(step), C.time_t(startTime.Unix()),
		C.int(len(values)), getCStringArrayPointer(cvalues))

	if int(ret) != 0 {
		err = errors.New(getError())
	}
	return
}

// The Update function feeds new data values into an RRD. The data is time aligned
// (interpolated) according to the properties of the RRD to which the data is written.
//
//      filename::
//          The name of the RRD you want to create. RRD files should end with the
//          extension .rrd. However, it accept any filename.
//      template::
//          The template switch allows you to specify which data sources you are going
//          to update and in which order. If the data sources specified in the
//          template are not available in the RRD file, the update process will
//          abort with an error. Format: "ds-name[:ds-name]..."
//      values::
//          A list of strings identifying values to be updated with corresponding
//          timestamps.
//
// See http://oss.oetiker.ch/rrdtool/doc/rrdupdate.en.html for detauls.
//
func Update(filename, template string, values []string) (err error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	ctemplate := C.CString(template)
	defer C.free(unsafe.Pointer(ctemplate))

	cvalues := makeCStringArray(values)
	defer freeCStringArray(cvalues)

	clearError()
	ret := C.rrd_update_r(cfilename, ctemplate,
		C.int(len(values)), getCStringArrayPointer(cvalues))

	if int(ret) != 0 {
		err = errors.New(getError())
	}
	return
}

// The UpdateValues function wraps the Update function, but provides a layer
// of abstraction, in that it does not require manual formatting of the
// RRD values being passed.
//
//      filename::
//          The name of the RRD you want to create. RRD files should end with the
//          extension .rrd. However, it accept any filename.
//      template::
//          The template switch allows you to specify which data sources you are going
//          to update and in which order. If the data sources specified in the
//          template are not available in the RRD file, the update process will
//          abort with an error. Format: "ds-name[:ds-name]..."
//      values::
//          A list of RrdValues identifying values to be updated with
//          corresponding timestamps.
//
func UpdateValues(filename, template string, values []RrdValue) (err error) {
	rrds := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		rrds[i] = values[i].ToString()
	}
	err = Update(filename, template, rrds)
	return err
}

// Fetch retrieves the values represented by an RRD file.
func Fetch(filename string, cf int, startTime int64, endTime int64, step uint64) (dsCount uint64, dsNames []string, data []map[int64]float64, err error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	// Newer version requires indices
	cf_idx := CfToString(cf)
	if cf_idx == "" {
		err = errors.New(fmt.Sprintf("Unable to convert cf : %d", cf))
		return
	}
	ccf := C.CString(cf_idx)

	var cdsCount C.ulong
	cstep := C.ulong(step)
	var cdsNames **C.char
	var cdata *C.rrd_value_t

	cst := C.time_t(startTime)
	cet := C.time_t(endTime)

	ret := C.rrd_fetch_r(cfilename, ccf, &cst, &cet, &cstep, &cdsCount, &cdsNames, &cdata)
	if int(ret) != 0 {
		err = errors.New(getError())
		return
	}

	// Figure out count
	dsCount = uint64(cdsCount)

	// Decode names of data sources
	dsNames = make([]string, dsCount)
	data = make([]map[int64]float64, dsCount)

	nptr := *cdsNames
	for iter := 0; iter < int(dsCount); iter++ {
		dsNames[iter] = C.GoString(nptr)
		data[iter] = make(map[int64]float64)
		nptr = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cdsNames)) + uintptr(iter)))
	}

	vptr := cdata
	for ti := startTime + int64(step); ti <= endTime; ti += int64(step) {
		k := int64(ti)
		for ii := 0; ii < int(dsCount); ii++ {
			v := float64(*vptr)

			// Add to appropriate map
			data[ii][k] = v

			vptr = (*C.rrd_value_t)(unsafe.Pointer(uintptr(unsafe.Pointer(vptr)) + 1))
		}
	}

	return
}

// Last retrieves the last timestamp for a value represented by RRD data.
func Last(filename string) time.Time {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	ret := C.rrd_last_r(cfilename)

	return time.Unix(int64(ret), 0)
}

func Dump(filename string, out string) (err error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	cout := C.CString(out)
	defer C.free(unsafe.Pointer(cout))

	ret := C.rrd_dump_r(cfilename, cout)
	if int(ret) != 0 {
		err = errors.New(getError())
	}

	return
}

//----- Helper methods ---------------------------------------------------------

func getError() string {
	return C.GoString(C.rrd_get_error())
}

func clearError() {
	C.rrd_clear_error()
}

func makeCStringArray(values []string) (cvalues []*C.char) {
	cvalues = make([]*C.char, len(values))
	for i := range values {
		cvalues[i] = C.CString(values[i])
	}
	return
}

func freeCStringArray(cvalues []*C.char) {
	for i := range cvalues {
		C.free(unsafe.Pointer(cvalues[i]))
	}
}

func getCStringArrayPointer(cvalues []*C.char) **C.char {
	return (**C.char)(unsafe.Pointer(&cvalues[0]))
}

func stringToCf(cf string) int {
	switch cf {
	case "AVERAGE":
		return CF_AVERAGE
	case "MIN":
		return CF_MINIMUM
	case "MAX":
		return CF_MAXIMUM
	case "LAST":
		return CF_LAST
	case "HWPREDICT":
		return CF_HWPREDICT
	case "MHWPREDICT":
		return CF_MHWPREDICT
	case "DEVPREDICT":
		return CF_DEVPREDICT
	case "SEASONAL":
		return CF_SEASONAL
	case "DEVSEASONAL":
		return CF_DEVSEASONAL
	case "FAILURES":
		return CF_FAILURES
	default:
		return -1
	}
	return -1
}

func cfToString(cf int) string {
	switch cf {
	case CF_AVERAGE:
		return "AVERAGE"
	case CF_MINIMUM:
		return "MIN"
	case CF_MAXIMUM:
		return "MAX"
	case CF_LAST:
		return "LAST"
	case CF_HWPREDICT:
		return "HWPREDICT"
	case CF_SEASONAL:
		return "SEASONAL"
	case CF_DEVPREDICT:
		return "DEVPREDICT"
	case CF_DEVSEASONAL:
		return "DEVSEASONAL"
	case CF_FAILURES:
		return "FAILURES"
	case CF_MHWPREDICT:
		return "MHWPREDICT"
	default:
		return ""
	}
	return ""
}
