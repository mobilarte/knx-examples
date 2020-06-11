package main

import (
//	"fmt"
 //       "github.com/vapourismo/knx-go/knx/dpt"
	"github.com/robfig/cron/v3"
//	"reflect"
)

var scheduler *cron.Cron

func runScheduler () {
	scheduler = cron.New()
	//scheduler.AddFunc("@every 0h0m60s", GroupUpdater)
	//scheduler.AddFunc("@every 0h0m15s", LuminosityChecker)
	scheduler.Start()
}

/*
func LuminosityChecker() {

	// get current value
	t, _ := groupList.DptByName("temperatur")

	//tv := reflect.ValueOf(t).Elem()
	//v := tv.Interface().(dpt.DPT_9001)
	// works
	v := float32(reflect.ValueOf(t).Elem().Interface().(dpt.DPT_9001))
	fmt.Printf("value as float32: %f \n", v)
	if v < 20 {
		fmt.Println("cold!")
	} else {
		fmt.Println("warm!")
	}
	south, _ := groupList.DptByName("helligkeit.süd")
	vs := float32(reflect.ValueOf(south).Elem().Interface().(dpt.DPT_9004))

	fmt.Printf("lux south: %7.2f\n", vs)

	// convert to float32
	fmt.Printf("value as string: %s \n", t)
	t_type := reflect.TypeOf(t)
	fmt.Printf("type as string: %s \n", t_type)
	t_value := reflect.ValueOf(t).Elem()
	fmt.Printf("value as string: %s \n", t_value)
	x := t_value.Interface().(dpt.DPT_9001) 
	f := float32(x)
	fmt.Printf("value as float32: %12.8f\n", f)

	south, _ := groupList.DptByName("helligkeit.süd")
	east, _ := groupList.DptByName("helligkeit.ost")
	west, _ := groupList.DptByName("helligkeit.west")
	wind, _ := groupList.DptByName("wind")

	//if east < 100 {
	//	fmt.Print("most likely nighttime!")
	//	}
}
*/

