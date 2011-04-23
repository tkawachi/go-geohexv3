package main

import (
	"fmt"
	"geohexv3"
)


func main() {
	fmt.Println("hello")
	z := geohexv3.GetZoneByLocation(30.0, 120.1, 5)
	if z != nil {
		fmt.Println(z.Code)
		fmt.Println(z)
	}
	z2 := geohexv3.GetZoneByCode(z.Code)
	fmt.Println(z2)
}
