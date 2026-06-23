package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lucian95511/and/internal/forum"
	"github.com/lucian95511/and/internal/pluginapi"
)

var manifest = pluginapi.Manifest{
	Name:        "konu_ac",
	Label:       "",
	Version:     "2.0.0",
	Description: "Forum'da yeni konu oluşturma ve taslak yönetimi (menüde gizli, forumdan n ile açılır)",
	Author:      "AND",
}

const (
	maxBaslik = 100
	maxIcerik = 2000
)

var kategoriler = forum.Categories

type taslak struct {
	Kategori    string `json:"kategori"`
	Baslik      string `json:"baslik"`
	Icerik      string `json:"icerik"`
	KaliciTalep bool   `json:"kalici_talep,omitempty"`
}

type taslakDosya struct {
	Taslaklar []taslak `json:"taslaklar"`
}

func taslakYolu(dataDir, kategori string) string {
	return filepath.Join(dataDir, "taslaklar_"+kategori+".json")
}

func taslakOku(dataDir, kategori string) []taslak {
	data, err := os.ReadFile(taslakYolu(dataDir, kategori))
	if err != nil {
		return nil
	}
	var df taslakDosya
	_ = json.Unmarshal(data, &df)
	return df.Taslaklar
}

func taslakYaz(dataDir, kategori string, ts []taslak) {
	data, _ := json.Marshal(taslakDosya{Taslaklar: ts})
	_ = os.WriteFile(taslakYolu(dataDir, kategori), data, 0o600)
}

var (
	stTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("34"))
	stSel    = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("22")).Foreground(lipgloss.Color("230"))
	stNormal = lipgloss.NewStyle()
	stOk     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	stErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	stDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	stField  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--manifest" {
		data, _ := json.Marshal(manifest)
		fmt.Println(string(data))
		return
	}

	client, err := pluginapi.NewClientFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, "konu_ac:", err)
		os.Exit(1)
	}

	preCategory := pluginapi.Category()
	dataDir := pluginapi.DataDir()

	m := newModel(client, preCategory, dataDir)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "konu_ac:", err)
		os.Exit(1)
	}
}

type ekran int

const (
	ekKategori ekran = iota
	ekForm
	ekTaslak
)

type gonderiMsg struct{ err error }

type model struct {
	client      *pluginapi.Client
	dataDir     string
	ekran       ekran
	kategoriler []string
	katIdx      int
	kategori    string

	baslikInput textinput.Model
	icerikInput textarea.Model
	kaliciTalep bool

	taslaklar []taslak
	taslakIdx int

	notice string
	isErr  bool
	w, h   int
}

func newModel(c *pluginapi.Client, preCategory, dataDir string) model {
	baslik := textinput.New()
	baslik.Placeholder = "konu başlığı…"
	baslik.CharLimit = maxBaslik
	baslik.Focus()

	icerik := textarea.New()
	icerik.Placeholder = "konu içeriği…"
	icerik.SetHeight(10)
	icerik.CharLimit = maxIcerik
	icerik.ShowLineNumbers = false

	kats := make([]string, len(kategoriler))
	copy(kats, kategoriler)

	m := model{
		client:      c,
		dataDir:     dataDir,
		kategoriler: kats,
		baslikInput: baslik,
		icerikInput: icerik,
	}

	if preCategory != "" {
		for i, k := range kats {
			if strings.EqualFold(k, preCategory) {
				m.katIdx = i
				break
			}
		}
		m.kategori = kats[m.katIdx]
		m.ekran = ekForm
		m.taslaklar = taslakOku(dataDir, m.kategori)
	} else {
		m.ekran = ekKategori
	}
	return m
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func gonderiCmd(c *pluginapi.Client, category, title, body string, kalici bool) tea.Cmd {
	return func() tea.Msg {
		return gonderiMsg{err: c.CreatePost(category, title, body, kalici)}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.icerikInput.SetWidth(msg.Width - 8)
		m.icerikInput.SetHeight(msg.Height - 16)
		return m, nil

	case gonderiMsg:
		if msg.err != nil {
			m.notice = "Gönderme hatası: " + msg.err.Error()
			m.isErr = true
			return m, nil
		}
		taslakYaz(m.dataDir, m.kategori, nil)
		return m, tea.Quit

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	if m.ekran == ekForm {
		var baslikCmd, icerikCmd tea.Cmd
		m.baslikInput, baslikCmd = m.baslikInput.Update(msg)
		m.icerikInput, icerikCmd = m.icerikInput.Update(msg)
		return m, tea.Batch(baslikCmd, icerikCmd)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.ekran {
	case ekKategori:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.katIdx > 0 {
				m.katIdx--
			}
		case "down", "j":
			if m.katIdx < len(m.kategoriler)-1 {
				m.katIdx++
			}
		case "enter":
			m.kategori = m.kategoriler[m.katIdx]
			m.taslaklar = taslakOku(m.dataDir, m.kategori)
			m.ekran = ekForm
			m.baslikInput.Focus()
		}

	case ekForm:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			baslik := strings.TrimSpace(m.baslikInput.Value())
			icerik := strings.TrimSpace(m.icerikInput.Value())
			if baslik != "" || icerik != "" {
				ts := append(m.taslaklar, taslak{
					Kategori:    m.kategori,
					Baslik:      baslik,
					Icerik:      icerik,
					KaliciTalep: m.kaliciTalep,
				})
				taslakYaz(m.dataDir, m.kategori, ts)
			}
			return m, tea.Quit
		case "ctrl+s":
			baslik := strings.TrimSpace(m.baslikInput.Value())
			icerik := strings.TrimSpace(m.icerikInput.Value())
			if baslik == "" {
				m.notice = "Başlık boş olamaz."
				m.isErr = true
				return m, nil
			}
			if icerik == "" {
				m.notice = "İçerik boş olamaz."
				m.isErr = true
				return m, nil
			}
			return m, gonderiCmd(m.client, m.kategori, baslik, icerik, m.kaliciTalep)
		case "ctrl+p":
			m.kaliciTalep = !m.kaliciTalep
		case "ctrl+t":
			if len(m.taslaklar) > 0 {
				m.ekran = ekTaslak
				m.taslakIdx = 0
			}
		default:
			var cmd tea.Cmd
			m.baslikInput, cmd = m.baslikInput.Update(msg)
			return m, cmd
		}

	case ekTaslak:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.ekran = ekForm
		case "up", "k":
			if m.taslakIdx > 0 {
				m.taslakIdx--
			}
		case "down", "j":
			if m.taslakIdx < len(m.taslaklar)-1 {
				m.taslakIdx++
			}
		case "enter":
			if m.taslakIdx < len(m.taslaklar) {
				ts := m.taslaklar[m.taslakIdx]
				m.baslikInput.SetValue(ts.Baslik)
				m.icerikInput.SetValue(ts.Icerik)
				m.kaliciTalep = ts.KaliciTalep
				m.taslaklar = append(m.taslaklar[:m.taslakIdx], m.taslaklar[m.taslakIdx+1:]...)
				taslakYaz(m.dataDir, m.kategori, m.taslaklar)
				m.ekran = ekForm
			}
		case "d", "D":
			if m.taslakIdx < len(m.taslaklar) {
				m.taslaklar = append(m.taslaklar[:m.taslakIdx], m.taslaklar[m.taslakIdx+1:]...)
				taslakYaz(m.dataDir, m.kategori, m.taslaklar)
				if m.taslakIdx >= len(m.taslaklar) && m.taslakIdx > 0 {
					m.taslakIdx--
				}
				if len(m.taslaklar) == 0 {
					m.ekran = ekForm
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.ekran {
	case ekKategori:
		return m.viewKategori()
	case ekForm:
		return m.viewForm()
	case ekTaslak:
		return m.viewTaslak()
	}
	return ""
}

func (m model) viewKategori() string {
	var sb strings.Builder
	sb.WriteString(stTitle.Render("AND — Yeni Konu") + "\n\n")
	sb.WriteString("Kategori seç:\n\n")
	for i, k := range m.kategoriler {
		if i == m.katIdx {
			sb.WriteString(stSel.Render("  "+k) + "\n")
		} else {
			sb.WriteString(stNormal.Render("  "+k) + "\n")
		}
	}
	sb.WriteString("\n" + stDim.Render("↑/↓ seç   enter devam   q çıkış"))
	return sb.String()
}

func (m model) viewForm() string {
	var sb strings.Builder
	kaliciStr := ""
	if m.kaliciTalep {
		kaliciStr = "  " + stOk.Render("[★ Kalıcılık talep ediliyor]")
	}
	sb.WriteString(stTitle.Render(fmt.Sprintf("AND — Yeni Konu  [%s]%s", m.kategori, kaliciStr)) + "\n\n")

	taslakHint := ""
	if len(m.taslaklar) > 0 {
		taslakHint = fmt.Sprintf("  %s", stDim.Render(fmt.Sprintf("(%d taslak — ctrl+t)", len(m.taslaklar))))
	}
	sb.WriteString(stDim.Render("Başlık:") + taslakHint + "\n")
	sb.WriteString(stField.Render(m.baslikInput.View()) + "\n\n")
	sb.WriteString(stDim.Render("İçerik:\n"))
	sb.WriteString(stField.Render(m.icerikInput.View()) + "\n\n")

	if m.isErr {
		sb.WriteString(stErr.Render(m.notice) + "\n")
	} else if m.notice != "" {
		sb.WriteString(stOk.Render(m.notice) + "\n")
	}
	sb.WriteString(stDim.Render("ctrl+s gönder   ctrl+p kalıcılık talebi   esc taslak kaydet ve çık"))
	return sb.String()
}

func (m model) viewTaslak() string {
	var sb strings.Builder
	sb.WriteString(stTitle.Render("Taslaklar") + "\n\n")
	for i, ts := range m.taslaklar {
		line := fmt.Sprintf("%s  %s", ts.Baslik, stDim.Render(fmt.Sprintf("(%d karakter)", len(ts.Icerik))))
		if i == m.taslakIdx {
			sb.WriteString(stSel.Render(line) + "\n")
		} else {
			sb.WriteString(stNormal.Render(line) + "\n")
		}
	}
	sb.WriteString("\n" + stDim.Render("enter yükle   d sil   esc geri"))
	return sb.String()
}
