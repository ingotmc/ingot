package main

import (
	"fmt"
	"github.com/ingotmc/ingot/anvil"
	"github.com/ingotmc/ingot/mc"
	"log"
	"os"
)

func main() {
	paletteFile, err := os.Open("./cmd/ingot/generated/reports/blocks.json")
	if err != nil {
		log.Fatal(err)
	}
	mc.GlobalPalette.FromJson(paletteFile)
	paletteFile.Close()
	dim := anvil.Dimension("./.gamesave/saves/Ingot Superflat/region")
	for z := -1; z <= 1; z++ {
		for x := -1; x <= 1; x++ {
			block, err := dim.BlockAt(x, 3, z)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("block at (%d, %d, %d): %+v\n", x, 3, z, block)
		}
	}
	return
}
