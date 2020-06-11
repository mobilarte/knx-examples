// Copyright 2020 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

// Log windspeed, temperature, luminosity from an MDT weather station into
// a csv file for further processing. Flushes log to file after a defined
// number of measures have been received or when receiving sigUSR1 signal.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
)

const (
	TIME_FORMAT     = "2006/01/02 15:04:05"
	FLUSH_FREQUENCY = 10 // regularly flush data to disk after 10 measures received.
	// increase to meaningful value!
	FILENAME = "log.csv" // no path component
)

type logEntry struct {
	time  string
	src   string
	value string
}

var logEntries []logEntry
var logMutex sync.Mutex
var verbose *bool

// Individual address of the weather station
var weatherStation = "1.1.7"

func logData(message knx.GroupEvent) {

	var h dpt.DPT_9004 // Illuminance in [Lux]
	var t dpt.DPT_9001 // Temperature in [°C]
	var w dpt.DPT_9005 // Windspeed in [m/s]

	tm := time.Now()

	dstAddr := message.Destination.String()
	switch dstAddr {
	case "1/2/0":
		h.Unpack(message.Data)
		l := logEntry{time: tm.Format(TIME_FORMAT), src: "HS", value: fmt.Sprintf("%f", float32(h))}
		logEntries = append(logEntries, l)
		if *verbose {
			fmt.Printf("%s,Illuminance South,%s\n", tm.Format(TIME_FORMAT), h)
		}
	case "1/2/1":
		h.Unpack(message.Data)
		l := logEntry{time: tm.Format(TIME_FORMAT), src: "HW", value: fmt.Sprintf("%f", float32(h))}
		logEntries = append(logEntries, l)
		if *verbose {
			fmt.Printf("%s,Illuminance West,%s\n", tm.Format(TIME_FORMAT), h)
		}
	case "1/2/2":
		h.Unpack(message.Data)
		l := logEntry{time: tm.Format(TIME_FORMAT), src: "HO", value: fmt.Sprintf("%f", float32(h))}
		logEntries = append(logEntries, l)
		if *verbose {
			fmt.Printf("%s,Illuminance East,%s\n", tm.Format(TIME_FORMAT), h)
		}
	case "1/2/6":
		t.Unpack(message.Data)
		l := logEntry{time: tm.Format(TIME_FORMAT), src: "T", value: fmt.Sprintf("%f", float32(t))}
		logEntries = append(logEntries, l)
		if *verbose {
			fmt.Printf("%s,Temperature,%s\n", tm.Format(TIME_FORMAT), t)
		}
	case "1/2/7":
		w.Unpack(message.Data)
		l := logEntry{time: tm.Format(TIME_FORMAT), src: "W", value: fmt.Sprintf("%f", float32(w))}
		logEntries = append(logEntries, l)
		if *verbose {
			fmt.Printf("%s,Windspeed,%s\n", tm.Format(TIME_FORMAT), w)
		}
	default:
		fmt.Printf("%s,ignoring message from source", tm.Format(TIME_FORMAT))
	}
}

func writeToFile() {
	f, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot create logfile")
	}
	defer f.Close()

	logMutex.Lock()
	for _, row := range logEntries {
		s := fmt.Sprintf("%s,%s,%s\n", row.time, row.src, row.value)
		_, err := f.WriteString(s)
		if err != nil {
			log.Println("Cannot write to file: #%v\n", row)
		}
	}
	logMutex.Unlock()
	*verbose = !*verbose // toggle verbose, in case it becomes annoying
	logEntries = nil
	f.Sync()
}

func main() {
	verbose = flag.Bool("verbose", false, "Be verbose, print all messages to stdout")

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		for {
			sig := <-sigs
			fmt.Printf("Signal %s received, flushing logs to file %s, (toggling verbose)\n", sig, FILENAME)
			writeToFile()
			if sig != syscall.SIGUSR1 {
				os.Exit(0)
			}
		}
	}()

	srcAddr, err := cemi.NewIndividualAddrString(weatherStation)

	groupRouter, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
	if err != nil {
		log.Fatal("Could not create GroupRouter")
	}
	for {
		message, open := <-groupRouter.Inbound()
		if !open {
			log.Fatal("Channel is closed")
		} else if message.Source == srcAddr && message.Command == knx.GroupWrite {
			logData(message)
			if len(logEntries) == FLUSH_FREQUENCY {
				writeToFile()
			}
		}
	}
}
