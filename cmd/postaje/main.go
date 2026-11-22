// Ukaz postaje prikaže delovanje paketa vreme: ustvari nekaj vremenskih
// postaj, jim doda meritve in izpiše meritve ter poročilo.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ZanTriler/vremenske-postaje/vreme"
)

func main() {
	cmd := &cli.Command{
		Name:  "postaje",
		Usage: "izpiše meritve in poročilo za vremenske postaje",
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "minTemp",
				Value: -50.0,
				Usage: "najnižja dovoljena temperatura meritve (°C)",
			},
			&cli.Float64Flag{
				Name:  "maxTemp",
				Value: 50.0,
				Usage: "najvišja dovoljena temperatura meritve (°C)",
			},
			&cli.IntFlag{
				Name:  "stMeritev",
				Value: 3,
				Usage: "najmanjše število meritev, potrebnih za izračun statistike",
			},
		},
		Action: zagon,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func zagon(ctx context.Context, cmd *cli.Command) error {
	minTemp := cmd.Float64("minTemp")
	maxTemp := cmd.Float64("maxTemp")
	stMeritev := int(cmd.Int("stMeritev"))

	if minTemp > maxTemp {
		return fmt.Errorf("minTemp (%v) ne sme biti večja od maxTemp (%v)", minTemp, maxTemp)
	}

	postaje := make(map[string]vreme.Postaja)

	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Center", "Ljubljana")
	vreme.DodajPostajo(postaje, "KP01", "Luka Koper", "Koper")
	vreme.DodajPostajo(postaje, "KR01", "Kredarica", "Triglav")
	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")

	začetne := map[string][]float64{
		"LJ01": {18.2, 21.5, 19.8},
		"KP01": {24.5, 27, 26.1},
		"KR01": {2, -1.5, 3.5},
		"MB01": {12.4},
	}
	for oznaka, meritve := range začetne {
		for _, m := range meritve {
			vreme.DodajMeritev(postaje, oznaka, m, minTemp, maxTemp)
		}
	}

	// uspešno
	vreme.DodajMeritev(postaje, "KP01", 28.4, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "LJ01", 16.5, minTemp, maxTemp)

	// neuspešno
	vreme.DodajMeritev(postaje, "KR01", -65.0, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "KP01", 55.2, minTemp, maxTemp)

	// postaja ne obstaja
	vreme.DodajMeritev(postaje, "CE01", 20.0, minTemp, maxTemp)
	fmt.Println()

	vreme.IzpisMeritev(postaje)
	fmt.Println()

	vreme.IzpisPorocila(postaje, stMeritev)
	return nil
}
