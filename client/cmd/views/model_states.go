package views

import (
	"github.com/charmbracelet/bubbles/list"
)

// AppState defines the state of the application
// Each state maps to a sub-model that handles it

type AppState int

const (
	StateSearching AppState = iota
	StateSelecting
	StateProcessing
	StateQuitting
)

type SearchFinishedMsg struct {
	Results []list.Item
	Err     error
}

type SongSelectedMsg struct {
	Item list.Item
}

type CallbackFinishedMsg struct {
}
