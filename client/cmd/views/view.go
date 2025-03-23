package views

import "net/rpc"

// UIContext holds shared context for all views
type UIContext struct {
	Client *rpc.Client
	Theme  Themes
}

// View is the interface for all simple render views
type View interface {
	Render(ctx UIContext) error
}

// InteractiveView is for Bubble Tea-based views
type InteractiveView interface {
	Run(ctx UIContext) error
}
