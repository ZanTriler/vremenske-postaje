// Package vreme omogoča delo z vremenskimi postajami: dodajanje postaj in
// meritev temperature ter izpis meritev in statističnih poročil.
//
// Polja strukture Postaja niso izvožena, zato postaje ustvarimo s funkcijo
// DodajPostajo, ki jih shrani v slovar z oznako postaje kot ključem.
package vreme

import (
	"fmt"
	"sort"
)

// Postaja predstavlja vremensko postajo z imenom, krajem in seznamom
// izmerjenih temperatur v °C.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo v slovar postaj doda novo postajo z dano oznako, imenom in
// krajem. Postaja nima nobene meritve. Če postaja z isto oznako že obstaja,
// se izpiše opozorilo, slovar pa ostane nespremenjen.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	if _, obstaja := postaje[oznaka]; obstaja {
		fmt.Printf("Postaja %s že obstaja\n", oznaka)
		return
	}
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
}

// DodajMeritev postaji z dano oznako doda meritev temperature. Meritev je
// veljavna, če leži med minTemp in maxTemp (vključno). Če postaja ne obstaja
// ali temperatura ni veljavna, se izpiše napaka, meritev pa se ne doda.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) {
	postaja, obstaja := postaje[oznaka]
	if !obstaja {
		fmt.Printf("Postaja %s ne obstaja\n", oznaka)
		return
	}

	if temperatura < minTemp || temperatura > maxTemp {
		fmt.Printf("Neveljavna temperatura %f\n", temperatura)
		return
	}

	postaja.meritve = append(postaja.meritve, temperatura)
	postaje[oznaka] = postaja
}

// statistika izračuna najnižjo, najvišjo in povprečno temperaturo postaje.
// Če postaja ne obstaja ali ima manj kot stMeritev meritev, vrne ok = false.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	postaja, obstaja := postaje[oznaka]
	if !obstaja || len(postaja.meritve) < stMeritev || len(postaja.meritve) == 0 {
		return 0, 0, 0, false
	}

	najnizja = postaja.meritve[0]
	najvisja = postaja.meritve[0]
	vsota := 0.0

	for _, m := range postaja.meritve {
		if m < najnizja {
			najnizja = m
		}
		if m > najvisja {
			najvisja = m
		}
		vsota += m
	}

	povprecna = vsota / float64(len(postaja.meritve))
	return najnizja, najvisja, povprecna, true
}

// oznake vrne urejen seznam oznak postaj.
func oznake(postaje map[string]Postaja) []string {
	kljuci := make([]string, 0, len(postaje))
	for k := range postaje {
		kljuci = append(kljuci, k)
	}
	sort.Strings(kljuci)
	return kljuci
}

// IzpisMeritev izpiše vse postaje (po abecednem redu oznak) z njihovimi
// meritvami.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for _, oznaka := range oznake(postaje) {
		p := postaje[oznaka]
		fmt.Printf("%s - %s (%s): %v\n", oznaka, p.ime, p.kraj, p.meritve)
	}
}

// IzpisPorocila za vsako postajo izpiše najnižjo, najvišjo in povprečno
// temperaturo ter opis (Vroče, Toplo, Hladno ali Mrzlo). Postaje z manj kot
// stMeritev meritvami so označene s sporočilom, da je meritev premalo.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	for _, oznaka := range oznake(postaje) {
		p := postaje[oznaka]
		min, max, povp, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev\n", p.ime)
			continue
		}

		var opis string
		switch {
		case povp >= 25.0:
			opis = "Vroče"
		case povp >= 15.0:
			opis = "Toplo"
		case povp >= 5.0:
			opis = "Hladno"
		default:
			opis = "Mrzlo"
		}

		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", p.ime, min, max, povp, opis)
	}
}
