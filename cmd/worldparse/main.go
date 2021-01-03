package main

import (
	"fmt"
	"github.com/ingotmc/ingot/anvil"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/simulation"
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
	world := simulation.NewWorldStore(
		anvil.Dimension("./.gamesave/saves/Ingot/region"),
		anvil.Dimension("./.gamesave/saves/Ingot/DIM-1"),
		anvil.Dimension("./.gamesave/saves/Ingot/DIM1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Print("enter space separated coordinates: ")
		var x, y, z int
		n, err := fmt.Scanln(&x, &y, &z)
		if n != 3 || err != nil {
			fmt.Println("error:", err)
			continue
		}
		coords := mc.BlockCoords{
			X: x,
			Y: y,
			Z: z,
		}
		chunk, err := world.Dimension(mc.Overworld).ChunkAt(coords.ChunkCoords())
		if err != nil {
			fmt.Printf("\terror: %v", err)
			continue
		}
		block, err := chunk.BlockAt(coords)
		if err != nil {
			fmt.Println("\terror:", err)
			continue
		}
		fmt.Println("\t", block)
	}
	return
}
