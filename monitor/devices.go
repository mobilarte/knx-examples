package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/vapourismo/knx-go/knx/cemi"
	"sync"
	"io"
	"os"
)

type Device struct {
	Name    string
	Type    string
	Address cemi.IndividualAddr
}

type Devices map[string]*Device

var deviceList Devices
var deviceListMutex sync.Mutex

func (dv *Devices) String() string {
	str := ""
	for _, entryPtr := range *dv {
		str = str + fmt.Sprintf("%8s %-45s %15s\n", entryPtr.Address.String(), entryPtr.Name, entryPtr.Type) 
	}
	return str
}

// Lookup returns the name and address as strings, given the individual address.
// Performs a linear search through the device list.
func (dv *Devices) Lookup(src cemi.IndividualAddr) (string, string) {
	for _, entryPtr := range *dv {
		if (*entryPtr).Address == src {
			return (*entryPtr).Name, (*entryPtr).Address.String()
		}
	}
	return "unknown", "na"
}

func (dv *Devices) FromCSVFile(filename string) error {
	deviceList = make(Devices)
	deviceListMutex = sync.Mutex{}

	csvfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	r := csv.NewReader(bufio.NewReader(csvfile))
	r.Comma = ';'
	r.Comment = '#'
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		addr, err := cemi.NewIndividualAddrString(rec[2])
		if err != nil {
			return err
		}
		row := Device{Name: rec[1], Type: rec[0], Address: addr}
		if _, ok := (*dv)[rec[1]]; !ok {
			deviceListMutex.Lock()
			(*dv)[rec[1]] = &row
			deviceListMutex.Unlock()
		} else {
			s := fmt.Sprintf("Duplicate key: %s", rec[1])
			return errors.New(s)
		}
	}
	return nil
}
