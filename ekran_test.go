package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/lucian95511/and/dil/diltest"
	"github.com/lucian95511/and/internal/pluginapi"
)

// TestEkranlar_HerDilde: konu açma formu ve taslaklar ekranı her dilde
// terminali tam kaplar ve çevrilmemiş metin içermez. Kategori ağdaki Türkçe
// değeriyle seçilir, ekranda çevrilir.
func TestEkranlar_HerDilde(t *testing.T) {
	denetci := diltest.YeniDenetci(t, nil, diltest.Oku(t, "ceviri.json", ceviriVerisi))
	diltest.HerDilde(t, func(kod string) {
		for _, sz := range [][2]int{{120, 34}, {80, 24}} {
			m := newModel(&pluginapi.Client{}, t.TempDir(), "Yapay Zeka")
			next, _ := m.Update(tea.WindowSizeMsg{Width: sz[0], Height: sz[1]})
			form := next.(model)
			taslaklar := form
			taslaklar.scr = screenTaslak
			taslaklar.taslaklar = []taslak{{Kategori: "Yapay Zeka", Baslik: "Hello", Icerik: "body"}}
			for ad, m := range map[string]model{"form": form, "taslaklar": taslaklar} {
				ekran := m.View()
				satirlar := strings.Split(ekran, "\n")
				if len(satirlar) != sz[1] {
					t.Errorf("%s %s %dx%d: %d satır", kod, ad, sz[0], sz[1], len(satirlar))
				}
				for i, s := range satirlar {
					if g := lipgloss.Width(s); g != sz[0] {
						t.Errorf("%s %s %dx%d: satır %d genişliği %d: %q", kod, ad, sz[0], sz[1], i, g, ansi.Strip(s))
					}
				}
				for _, k := range denetci.Kalintilar(kod, ansi.Strip(ekran)) {
					t.Errorf("%s %s %dx%d: %s", kod, ad, sz[0], sz[1], k)
				}
			}
		}
	})
}
