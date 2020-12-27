package main

import (
	"fmt"
	"github.com/ingotmc/ingot"
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
	s.Start()
	fmt.Println("Goodbye.")
}
