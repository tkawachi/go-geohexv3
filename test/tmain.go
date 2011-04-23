package main

import (
	"fmt"
	"geohex3"
)


func main() {
	fmt.Println("hello")
	z := geohex3.GetZoneByLocation(30.0, 120.1, 5)
	if z != nil {
		fmt.Println(z.Code)
		fmt.Println(z)
	}
	z2 := geohex3.GetZoneByCode(z.Code)
	fmt.Println(z2)
}
