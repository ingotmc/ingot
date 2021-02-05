package main

import (
	"flag"
	"fmt"
	"github.com/ingotmc/ingot"
	"github.com/ingotmc/ingot/cmd/ingot/chunk"
	"github.com/ingotmc/ingot/cmd/ingot/player"
	"github.com/ingotmc/ingot/world"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/worldgen/noise/octavenoise"
	"log"
	"os"
	"path"
	"path/filepath"
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

	paletteFile, err := os.Open(filepath.Join(*gameFolder, "data", "blocks.json"))
	if err != nil {
		log.Fatal(fmt.Errorf("coudln't load palette file: %w", err))
	}
	mc.GlobalPalette.FromJson(paletteFile)
	_ = paletteFile.Close()

	w := world.World{
		Seed:          "ingot",
		PlayerService: player.NewService(5),
		Overworld: world.Dimension{
			ID:             mc.Overworld,
			ChunkStore:     new(chunk.Store),
			ChunkGenerator: chunk.NewSimplexGenerator(octavenoise.WithFrequency(1600.0), octavenoise.WithScalingFactor(0.3)),
		},
	}

	srv, err := ingot.NewServer(w)
	if err != nil {
		log.Fatal(err)
	}

	srv.Start()
	fmt.Println("Goodbye")
}
