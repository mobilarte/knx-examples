// Copyright 2020 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

package main

import (
	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/util"
	"log"
	"os"
	"time"
)

// DeviceSimulator is just a GroupRouter with a defined source address
type DeviceSimulator struct {
	GroupRouter knx.GroupRouter
	SrcAddr     cemi.IndividualAddr
}

// NewDeviceSimulator returns a DeviceSimulator with given source address
func NewDeviceSimulator(src string) (*DeviceSimulator, error) {
	srcAddr, err := cemi.NewIndividualAddrString(src)
	if err != nil {
		return nil, err
	}
	groupRouter, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
	if err != nil {
		return nil, err
	}
	ds := &DeviceSimulator{
		GroupRouter: groupRouter,
		SrcAddr:     srcAddr,
	}
	return ds, nil
}

// Close calls the GroupRouter Close
func (ds *DeviceSimulator) Close() {
	ds.GroupRouter.Close()
}

// Send sends a datapointvalue (DPT) to the dst group address
func (ds *DeviceSimulator) Send(dst string, d dpt.DatapointValue) error {
	dstAddr, err := cemi.NewGroupAddrString(dst)
	if err != nil {
		return err
	}
	if err = ds.GroupRouter.Send(knx.GroupEvent{
		Command:     knx.GroupWrite,
		Source:      ds.SrcAddr,
		Destination: dstAddr,
		Data:        d.Pack()}); err != nil {
		return err
	}
	return nil
}

// SendDate sends a date to the dst group address
func (ds *DeviceSimulator) SendDate(dst string) {
	var d dpt.DPT_11001

	// Take current date
	year, month, day := time.Now().Date()
	d.Year = uint16(year)
	d.Day = uint8(day)
	d.Month = uint8(month)

	util.Log(ds, "%#v", d)
	util.Log(ds, "%#v", d.Pack())

	if err := ds.Send(dst, &d); err != nil {
		log.Fatal(err)
	}
}

// SendTime sends a time to the dst group address
func (ds *DeviceSimulator) SendTime(dst string) {
	var t dpt.DPT_10001

	// Take current clock
	hour, min, sec := time.Now().Clock()
	_, _, day := time.Now().Date()
	t.Day = uint8(day)
	t.Hour = uint8(hour)
	t.Minutes = uint8(min)
	t.Seconds = uint8(sec)

	util.Log(ds, "%#v", t)
	util.Log(ds, "%#v", t.Pack())

	if err := ds.Send(dst, &t); err != nil {
		log.Fatal(err)
	}
}

// SendDateTime sends a date_time to the dst group address
func (ds *DeviceSimulator) SendDateTime(dst string) {
	var dt dpt.DPT_19001

	t := time.Now()
	hour, min, sec := t.Clock()
	year, month, day := t.Date()
	dt.Year = uint16(year)
	dt.Month = uint8(month)
	dt.DayOfMonth = uint8(day)
	dt.DayOfWeek = uint8(time.Now().Weekday())
	dt.HourOfDay = uint8(hour)
	dt.Minutes = uint8(min)
	dt.Seconds = uint8(sec)
	// Set summertime
	dt.SUTI = true
	// Set external clock synchronization, aka clock quality
	dt.CLQ = true

	util.Log(ds, "%#v", dt)
	util.Log(ds, "%#v", dt.Pack())

	if err := ds.Send(dst, &dt); err != nil {
		log.Fatal(err)
	}
}

// SendCounter sends a 2 byte unsigned counter to the dst group address
func (ds *DeviceSimulator) SendCounter(dst string, value uint16) {
	var d dpt.DPT_7001

	d = dpt.DPT_7001(value)
	util.Log(ds, "%#v", d)

	if err := ds.Send(dst, &d); err != nil {
		log.Fatal(err)
	}
}

// SendASCIIText sends a short (max. 14 chars) ASCII text to the dst group address
func (ds *DeviceSimulator) SendASCIIText(dst string, value string) {
	var d dpt.DPT_16000

	d = dpt.DPT_16000(value)
	util.Log(ds, "%#v", d)
	util.Log(ds, "%#v", d.Pack())

	if err := ds.Send(dst, &d); err != nil {
		log.Fatal(err)
	}
}

// SendSignedRelative sends a percentage [-128..127] to the dst group address
func (ds *DeviceSimulator) SendSignedRelative(dst string, value int8) {
	var d dpt.DPT_6001

	d = dpt.DPT_6001(value)
	util.Log(ds, "%#v", d)
	util.Log(ds, "%#v", d.Pack())

	if err := ds.Send(dst, &d); err != nil {
		log.Fatal(err)
	}
}

// Send a few message from source 1.1.6, the ABB Router does not care on which
// side the messages come from, it relays them onto the KNX network.
func main() {
	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	ds, err := NewDeviceSimulator("1.1.6")
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	ds.SendSignedRelative("1/5/0", -128)
	ds.SendCounter("1/5/1", 64000)
	ds.SendASCIIText("1/5/2", "KNX is OK, but")
	ds.SendDate("1/5/3")
	ds.SendTime("1/5/4")
	ds.SendDateTime("1/5/5")
}
