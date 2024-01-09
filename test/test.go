package main

import (
	"fmt"

	netpbm "github.com/jahsimelvin/Netpbm"
)

func main() {
	pbm, err := netpbm.ReadPBM("test.pbm")
	if err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", pbm.MagicNumber)
	fmt.Printf("Width: %d\n", pbm.Width)
	fmt.Printf("Height: %d\n", pbm.Height)
	fmt.Printf("Data:")
	for _, row := range pbm.Data {
		fmt.Println(row)
	}
}
