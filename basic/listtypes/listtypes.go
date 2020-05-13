// Licensed under the MIT license which can be found in the LICENSE file.

package main

import (
    "fmt"
    "sort"

    "github.com/vapourismo/knx-go/knx/dpt"
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

    fmt.Printf("%10s   %20s %20s\n", "DPT", "Default Value", "Unit")
    keys :=  dpt.ListSupportedTypes()
    sort.Sort(byMe(keys))

    for _, key := range keys {
        d, _ := dpt.Produce(key)
        fmt.Printf("%10s   %20s %20s\n", key, d, d.(dpt.DatapointMeta).Unit())
    }
    fmt.Println("")
}
