package test

import (
	"fmt"

	netpgm "github.com/jahsimelvin/Netpbm"
)

func testpgm() {
	filename := "../duckP2.pgm"
	pgm, err := netpgm.ReadPGM(filename)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", pgm.MagicNumber)
	fmt.Printf("Width: %d\n", pgm.Width)
	fmt.Printf("Height: %d\n", pgm.Height)
	fmt.Println("Data:")
	fmt.Println("Max: %s\n", pgm.Max)
	for _, row := range pgm.Data {
		fmt.Println(row)
	}

}
