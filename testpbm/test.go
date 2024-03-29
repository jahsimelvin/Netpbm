package test

import (
	"fmt"

	netpbm "github.com/jahsimelvin/Netpbm"
)

func testpbm() {
	filename := "../imageP1.pbm"
	pbm, err := netpbm.ReadPBM(filename)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", pbm.MagicNumber)
	fmt.Printf("Width: %d\n", pbm.Width)
	fmt.Printf("Height: %d\n", pbm.Height)
	fmt.Println("Data:")
	for _, row := range pbm.Data {
		fmt.Println(row)
	}

}
