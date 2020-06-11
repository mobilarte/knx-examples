package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
	"sync"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
)

type DataPoint struct {
	Name       string
	Type       string
	Address    cemi.GroupAddr
	DptType    string
	dpt        dpt.DatapointValue
	LastUpdate time.Time
}

const (
	TIME_FORMAT = "2006/01/02 15:04:05"
)

type DataPoints map[string]*DataPoint

//var groupList DataPoints
var groupListMutex sync.Mutex

func (d *DataPoint) Json() []byte {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", d.Name)
	s += fmt.Sprintf("\"address\": \"%s\",", d.Address)
	s += fmt.Sprintf("\"value\": \"%s\",", d.dpt)
	s += fmt.Sprintf("\"lastupdate\": \"%s\"", d.LastUpdate.Format(TIME_FORMAT))
	s += "}"
	return []byte(s)
}

func (gl DataPoints) String() string {
	s := ""
	return gl.StringFiltered(&s)
}

func (gl *DataPoints) StringFiltered(typeFilter *string) string {
	str := ""
	for _, value := range *gl {
		if *typeFilter == "" || *typeFilter == value.Type {
			v := reflect.ValueOf(value.dpt)
			str = str + fmt.Sprintf("%8s %10s %-45s %15s %9s %20s\n", value.Address.String(),
				value.Type, value.Name, v, value.DptType, value.LastUpdate.Format(TIME_FORMAT))
		}
	}
	return str
}

// Lookup returns the name, address and value as strings.
// Performs a linear search through the group list.
func (gl *DataPoints) LookupAndUpdate(destination cemi.GroupAddr, data []byte) (string, string) {
	for key, entryPtr := range *gl {
		if (*entryPtr).Address == destination {
			groupListMutex.Lock()
			(*entryPtr).dpt.Unpack(data)
			(*entryPtr).LastUpdate = time.Now()
			value := (*entryPtr).dpt.(fmt.Stringer).String()
			groupListMutex.Unlock()
			return key, value
		}
	}
	return "unknown", "na"
}

func (gl *DataPoints) Lookup(destination cemi.GroupAddr) (string, string, string) {
	for _, entry := range *gl {
		if entry.Address == destination {
			value := entry.dpt.(fmt.Stringer).String()
			return entry.Name, entry.Type, value
		}
	} 
	return "unknown", "na", "na"
}

func (gl *DataPoints) LookupByName(name string) ([]byte, error) {
	if entryPtr, ok := (*gl)[name]; !ok {
		s := fmt.Sprintf("{\"error\": \"notfound\"}")
		return []byte(s), errors.New("Group not found")
	} else { 
		return entryPtr.Json(), nil
	}
}

func (gl *DataPoints) DptByName(name string) (dpt.DatapointValue, error) {
	if entryPtr, ok := (*gl)[name]; !ok {
		return nil, errors.New("not found")
	} else { 
		return entryPtr.dpt, nil
	}
}

func (gl *DataPoints) FromCSVFile(filename string) error {
	groupListMutex = sync.Mutex{}

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
		grp, err := cemi.NewGroupAddrString(rec[2])
		if err != nil {
			return err
		}
		dpt, ok := dpt.Produce(rec[3])
		if !ok {
			return errors.New("DPT cannot be produced")
		}
		row := DataPoint{Name: rec[1], Type: rec[0], Address: grp, DptType: rec[3], dpt: dpt, LastUpdate: time.Now()}
		if _, ok := (*gl)[string(rec[1])]; !ok {
			groupListMutex.Lock()
			(*gl)[string(rec[1])] = &row
			groupListMutex.Unlock()
		} else {
			s := fmt.Sprintf("Duplicate key: %s", rec[1])
			return errors.New(s)
		}
	}
	return nil
}
