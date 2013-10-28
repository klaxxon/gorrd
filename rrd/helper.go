package rrd

import (
	"errors"
	"fmt"
	"time"
)

// Define helper routines for rrd package, for creating value strings, etc.

// Convenience method for forming an RRD value.
func CreateRrdValue(Time time.Time, Value int64) string {
	o := RrdValue{
		Time:  Time,
		Value: Value,
	}
	return o.ToString()
}

func CreateDSValue(Name string, Type DsType, Heartbeat int64, Min, Max float64) string {
	t, err := dsToString(Type)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("DS:%s:%s:%d:%f:%f", Name, t, Heartbeat, Min, Max)
}

func CreateComputeDSValue(Name string, RpnExpression string) string {
	return fmt.Sprintf("DS:%s:COMPUTE:%s", Name, RpnExpression)
}

func dsToString(dsType DsType) (string, error) {
	switch dsType {
	case DS_GAUGE:
		return "GAUGE", nil
		break
	case DS_COUNTER:
		return "COUNTER", nil
		break
	case DS_DERIVE:
		return "DERIVE", nil
		break
	case DS_ABSOLUTE:
		return "ABSOLUTE", nil
		break
	case DS_COMPUTE:
		return "COMPUTE", nil
		break
	default:
		return "", errors.New("Invalid DS type")
		break
	}
	return "", errors.New("Invalid DS type")
}
