package vreme

import "testing"

func testPostaje() map[string]Postaja {
	p := make(map[string]Postaja)
	DodajPostajo(p, "LJ01", "Ljubljana Center", "Ljubljana")
	for _, m := range []float64{18.2, 21.5, 19.8} {
		DodajMeritev(p, "LJ01", m, -50, 50)
	}
	return p
}

func TestDodajMeritevNeveljavna(t *testing.T) {
	p := testPostaje()
	DodajMeritev(p, "LJ01", 55.2, -50, 50)
	DodajMeritev(p, "LJ01", -65, -50, 50)
	DodajMeritev(p, "CE01", 20, -50, 50)
	if n := len(p["LJ01"].meritve); n != 3 {
		t.Errorf("pričakoval 3 meritve, dobil %d", n)
	}
	if _, ok := p["CE01"]; ok {
		t.Error("postaja CE01 ne bi smela nastati")
	}
}

func TestDodajMeritevMeje(t *testing.T) {
	p := testPostaje()
	DodajMeritev(p, "LJ01", 10, 0, 10)
	DodajMeritev(p, "LJ01", 11, 0, 10)
	if n := len(p["LJ01"].meritve); n != 4 {
		t.Errorf("pričakoval 4 meritve, dobil %d", n)
	}
}

func TestStatistika(t *testing.T) {
	p := testPostaje()
	DodajMeritev(p, "LJ01", 16.5, -50, 50)
	min, max, povp, ok := statistika(p, "LJ01", 3)
	if !ok || min != 16.5 || max != 21.5 || povp < 18.99 || povp > 19.01 {
		t.Errorf("napačna statistika: %v %v %v %v", min, max, povp, ok)
	}
}

func TestStatistikaPremaloMeritev(t *testing.T) {
	p := testPostaje()
	if _, _, _, ok := statistika(p, "LJ01", 4); ok {
		t.Error("pričakoval ok = false pri stMeritev = 4")
	}
	if _, _, _, ok := statistika(p, "XX01", 3); ok {
		t.Error("pričakoval ok = false za neobstoječo postajo")
	}
	if _, _, _, ok := statistika(p, "LJ01", 0); !ok {
		t.Error("pričakoval ok = true pri stMeritev = 0")
	}
}

func TestDodajPostajoDuplikat(t *testing.T) {
	p := testPostaje()
	DodajPostajo(p, "LJ01", "Druga", "Drugje")
	if p["LJ01"].ime != "Ljubljana Center" {
		t.Error("obstoječa postaja ne sme biti prepisana")
	}
}
