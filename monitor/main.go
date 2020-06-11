package main

import (
	//"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/vapourismo/knx-go/knx/util"
	"time"
)

var monitor *Monitor

func cleanup(dstFilter string) {
	//deviceList.DumpAsList()
	//dpList.DumpAsList()
	//dpList.DumpMinimal(dstFilter)
	fmt.Println("cleanup")
}

/*

func updateGroupEvent(message knx.GroupEvent, srcFilter *string, dstFilter *string) {

	// Still update the internal database
	srcName, _ := deviceList.Lookup(message.Source)

	// Log selectively
	//if (*srcFilter == "" && *dstFilter == "") || *srcFilter == srcAddr || *dstFilter == message.Destination.String() {
	dstName, dstType, _ := groupList.Lookup(message.Destination)
	if dstType != "test" {
		return
	}
	switch message.Command {
	case knx.GroupWrite:
		dstName, value := groupList.LookupAndUpdate(message.Destination, message.Data)
		log.Printf("write    %8s %-25s -> %8s %-45s %15s %20x", message.Source.String(), srcName, message.Destination.String(), dstName, value, message.Data)
	case knx.GroupRead:
		dstName, _ := groupList.Lookup(message.Destination)
		log.Printf("read     %8s %-25s -> %8s %-45s %15s %20x", message.Source.String(), srcName, message.Destination.String(), dstName, "", message.Data)
	case knx.GroupResponse:
		dstName, value := groupList.LookupAndUpdate(message.Destination, message.Data)
		log.Printf("response %8s %-25s -> %8s %-45s %15s %20x", message.Source.String(), srcName, message.Destination.String(), dstName, value, message.Data)
	default:
		log.Printf("SHOULD NEVER HAPPEN")
	}
	//}
}

// runGroupReader listens to multicasted group events.
func runGroupReader(srcFilter *string, dstFilter *string) {
	//log.Print("groupRouter is", grpRouter)
	for {
		msg, open := <-grpRouter.Inbound()
		if !open {
			log.Print("Channel is not open")
		} else {
			updateGroupEvent(msg, srcFilter, dstFilter)
		}
	}
	defer grpRouter.Close()
}

// GroupUpdater sends a periodic knx.GroupRead to all groups.
func GroupUpdater() {
	for {
		//log.Print("groupRouter is", grpRouter)
		for _, row := range groupList {
			if row.Type != "sensor" {
				c, ok := dpt.Produce(row.DptType)
				if !ok {
					log.Print(c)
					continue
				}
				err := grpRouter.Send(knx.GroupEvent{
					Command:     knx.GroupRead,
					Destination: row.Address,
					Data:        c.Pack()})
				if err != nil {
					log.Print(err)
				}
				time.Sleep(20 * time.Millisecond)
			}
		}
		time.Sleep(UPDATE_INTERVALL)
	}
}

func checkFilters(srcFilter *string, dstFilter *string) bool {
	if *srcFilter != "" {
		_, err := cemi.NewIndividualAddrString(*srcFilter)
		if err != nil {
			return false
		} else {
			return true
		}
	}
	if *dstFilter != "" {
		_, err := cemi.NewGroupAddrString(*dstFilter)
		if err != nil {
			return false
		} else {
			return true
		}
	}
	return true
}

*/
func main() {
	var err error

	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	/*dstFilter := flag.String("dst", "", "Destination to filter for: eg 1/2/7, a group address")
	srcFilter := flag.String("src", "", "Source to filter for: eg 1.0.1, an individual address")

	dumpGroups := flag.Bool("dump", false, "Only dump group addresses")
	typeFilter := flag.String("type", "", "Type to filter for dumping: eg \"light\"")

	flag.Parse() */

	sigs := make(chan os.Signal)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- sigs
		cleanup("")

		fmt.Println("Ctrl+C pressed in Terminal")
		fmt.Println("Exiting gracefully!")
		os.Exit(0)
	}()

	verbose := true
	monitor, err = NewMonitor("1.1.6", verbose)
	if err != nil {
		log.Fatal(err)
	}

	go monitor.GroupEventsListener()
	go monitor.StartScheduler()

	//monitor.PrintDevicesAndGroups()
	for {
		time.Sleep(1*time.Hour)
	}	

	/*if *dumpGroups {
		fmt.Print(groupList.StringFiltered(typeFilter))
	} else if !checkFilters(srcFilter, dstFilter) {
		log.Fatal("Bad filters!")
	} else {
		log.Printf("Filter active: src = [%s], dst = [%s]", *srcFilter, *dstFilter)

		grpRouter, err = knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
		if err != nil {
			log.Fatal(err)
		}
		go runGroupReader(srcFilter, dstFilter)
		go runScheduler()
		runRest()
	}
	*/

}
