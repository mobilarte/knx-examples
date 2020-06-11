// Copyright 2020 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

// Prints a list of supported DPTs, with default value and unit.
// Does not work with the original knx-go library, because of
// missing function IsValid() on dpts.
package main

import (
	"fmt"
	"github.com/vapourismo/knx-go/knx/dpt"
	"sort"
)

type byMe []string

func (s byMe) Len() int {
	return len(s)
}

func (s byMe) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byMe) Less(i, j int) bool {
	var il, ir int
	var jl, jr int

	fmt.Sscanf(s[i], "%d.%d", &il, &ir)
	fmt.Sscanf(s[j], "%d.%d", &jl, &jr)

	if il < jl {
		return true
	} else if il == jl {
		return ir < jr
	} else {
		return false
	}
}

func main() {

	fmt.Println("List of known DPTs")
	fmt.Println()

	fmt.Printf("%10s %25s %6s %35s\n", "DPT", "Default Value", "Valid?", "Unit")
	keys := dpt.ListSupportedTypes()
	sort.Sort(byMe(keys))

	for _, key := range keys {
		d, _ := dpt.Produce(key)
		// Trying to pack and unpack
		buf := d.Pack()
		d.Unpack(buf)
		fmt.Printf("%10s %25s %6t %35s\n", key, d, d.IsValid(), d.(dpt.DatapointMeta).Unit())
	}
	fmt.Println("")
}
