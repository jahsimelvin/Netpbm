package test

import (
	"fmt"

	netppm "github.com/jahsimelvin/Netpbm"
)

func testppm() {
	filename := "../duckP3.ppm"
	ppm, err := netppm.ReadPBM(filename)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", ppm.MagicNumber)
	fmt.Printf("Width: %d\n", ppm.Width)
	fmt.Printf("Height: %d\n", ppm.Height)
	fmt.Println("Data:")
	for _, row := range ppm.Data {
		fmt.Println(row)
	}

}