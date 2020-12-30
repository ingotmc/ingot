package main

import (
	"fmt"
	"github.com/ingotmc/ingot"
	"github.com/ingotmc/ingot/mc"
	"log"
	"os"
	"os/signal"
)

func main() {
	s, err := ingot.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		<-sigc
		s.Stop()
	}()
	paletteFile, err := os.Open("./cmd/ingot/generated/reports/blocks.json")
	if err != nil {
		log.Fatal(err)
	}
	mc.GlobalPalette.FromJson(paletteFile)
	paletteFile.Close()
	s.Start()
	fmt.Println("Goodbye.")
}
