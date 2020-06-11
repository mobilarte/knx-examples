package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/util"
	"log"
	//	"os"
	"time"
)

const UPDATE_INTERVALL = 360 * time.Second

// Monitor is just a GroupRouter with a defined source address
type Monitor struct {
	GroupRouter *knx.GroupRouter
	verbose     bool
	SrcAddr     cemi.IndividualAddr
	Devices     *Devices
	Groups      *DataPoints
	Scheduler   *cron.Cron
}

// NewMonitor returns a Monitor with given source address
func NewMonitor(src string, verbose bool) (*Monitor, error) {
	var groupList = make(DataPoints)
	var deviceList = make(Devices)

	srcAddr, err := cemi.NewIndividualAddrString(src)
	if err != nil {
		return nil, err
	}
	groupRouter, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
	if err != nil {
		return nil, err
	}

	err = deviceList.FromCSVFile("devices.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = groupList.FromCSVFile("groups.csv")
	if err != nil {
		log.Fatal(err)
	}

	scheduler = cron.New()

	mon := &Monitor{
		GroupRouter: &groupRouter,
		verbose:     verbose,
		SrcAddr:     srcAddr,
		Devices:     &deviceList,
		Groups:      &groupList,
		Scheduler:   cron.New(),
	}

	return mon, nil
}

func (mon *Monitor) updateGroupEvent(message knx.GroupEvent) {

	// Still update the internal database
	srcName, _ := mon.Devices.Lookup(message.Source)
	_, dstType, _ := mon.Groups.Lookup(message.Destination)
        if dstType != "test" {
                return
        }
	switch message.Command {
	case knx.GroupWrite:
		dstName, value := mon.Groups.LookupAndUpdate(message.Destination, message.Data)
		if mon.verbose {
			log.Printf("write    %8s %-25s -> %8s %-45s %15s %20x",
				message.Source.String(), srcName, message.Destination.String(),
				dstName, value, message.Data)
		}
	case knx.GroupRead:
		dstName,_,  _ := mon.Groups.Lookup(message.Destination)
		if mon.verbose {
			log.Printf("read     %8s %-25s -> %8s %-45s %15s %20x",
				message.Source.String(), srcName, message.Destination.String(),
				dstName, "", message.Data)
		}
	case knx.GroupResponse:
		dstName, value := mon.Groups.LookupAndUpdate(message.Destination, message.Data)
		if mon.verbose {
			log.Printf("response %8s %-25s -> %8s %-45s %15s %20x",
				message.Source.String(), srcName, message.Destination.String(),
				dstName, value, message.Data)
		}
	default:
		log.Printf("SHOULD NEVER HAPPEN")
	}
}

func (mon *Monitor) GroupEventsListener() {
	for {
		message, open := <-mon.GroupRouter.Inbound()
		if !open {
			util.Log(mon, "Channel is not open")
		} else {
			mon.updateGroupEvent(message)
		}
	}
}

// GroupUpdater sends a periodic knx.GroupRead to all groups.
func (mon *Monitor) GroupUpdater() {
	util.Log(mon, "GroupUpdater running")

	for _, row := range *mon.Groups {
		if row.Type != "sensor" {
			c, ok := dpt.Produce(row.DptType)
			if !ok {
				log.Print(c)
				continue
			}
			err := mon.GroupRouter.Send(knx.GroupEvent{
				Command: knx.GroupRead,
				Source:  mon.SrcAddr,
				Destination: row.Address,
				Data:    c.Pack()})
			if err != nil {
				log.Print(err)
			}
			// Throttle sending of messages
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (mon *Monitor) StartScheduler() {
	mon.Scheduler.AddFunc("@every 0h0m60s", mon.GroupUpdater)
	mon.Scheduler.Start()
}

func (mon *Monitor) PrintDevicesAndGroups() {
	fmt.Println("List of Group Addresses")
	fmt.Print(mon.Groups)
	fmt.Println("List of Devices")
	fmt.Print(mon.Devices)
}

// Close calls the GroupRouter Close
func (mon *Monitor) Close() {
	mon.Scheduler.Stop()
	mon.GroupRouter.Close()
}
