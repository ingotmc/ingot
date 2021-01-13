package main

import (
	"flag"
	"fmt"
	"github.com/ingotmc/ingot"
	"github.com/ingotmc/ingot/anvil"
	"github.com/ingotmc/ingot/cmd/ingot/player"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/ingot/world"
	"log"
	"os"
	"path"
)

const defaultGameFolder = ".local/share/ingot/" // join with homedir

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	gameFolder := flag.String("game_folder", path.Join(home, defaultGameFolder), "ingot server saves folder")

	flag.Parse()

	if _, err := os.Stat(*gameFolder); os.IsNotExist(err) {
		log.Fatal(fmt.Errorf("couldn't stat game_folder: %w", err))
	}

	paletteFile, err := os.Open(path.Join(*gameFolder, "/data/blocks.json"))
	if err != nil {
		log.Fatal(fmt.Errorf("coudln't load palette file: %w", err))
	}
	mc.GlobalPalette.FromJson(paletteFile)
	_ = paletteFile.Close()

	overworld := anvil.Dimension(path.Join(*gameFolder, "/region"))
	ws := world.NewWorldStore(overworld, nil, nil)
	world := world.World{
		Seed:       "ingot",
		WorldStore: ws,
	}

	sim := world.Simulation{
		PlayerService: player.NewService(5),
		World:         world,
	}

	srv, err := ingot.NewServer(sim)
	if err != nil {
		log.Fatal(err)
	}

	srv.Start()
	fmt.Println("Goodbye")
}
