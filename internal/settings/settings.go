package settings

import "github.com/EnotInc/Bard/config"

type Settings struct {
	bools  map[setting]bool
	change func()
	uptade bool
	w, h   int
	cursor int
}

type setting string

const (
	RelativeNumbers setting = "RelativeNumbers"
	ShowMDSymbols   setting = "ShowMDSymbols"
	ShowTabNames    setting = "ShowTabNames"
	EnableRender    setting = "EnableRender"
	ShowDotFiles    setting = "ShowDotFiles"
	ShowBorders     setting = "ShowBorders"
	ShowIcons       setting = "ShowIcons"
	ShowEmpty       setting = "ShowEmpty"
)

func IntiSettings(onChange func()) *Settings {
	s := &Settings{
		cursor: 0,
		change: onChange,
	}

	s.updateSettings()
	return s
}

func (s *Settings) updateSettings() {
	cfg := config.GetConfig()
	s.bools = map[setting]bool{
		RelativeNumbers: cfg.RLN,
		ShowMDSymbols:   cfg.ShowMD,
		ShowTabNames:    cfg.TabNames,
		EnableRender:    cfg.Render,
		ShowDotFiles:    cfg.ShowDot,
		ShowBorders:     cfg.ShowBorder,
		ShowIcons:       cfg.ShowIcons,
		ShowEmpty:       cfg.ShowEmpty,
	}
}

func (s *Settings) apply() {
	cfg := config.GetConfig()

	cfg.RLN = s.bools[RelativeNumbers]
	cfg.ShowMD = s.bools[ShowMDSymbols]
	cfg.TabNames = s.bools[ShowTabNames]
	cfg.Render = s.bools[EnableRender]
	cfg.ShowDot = s.bools[ShowDotFiles]
	cfg.ShowBorder = s.bools[ShowBorders]
	cfg.ShowIcons = s.bools[ShowIcons]
	cfg.ShowEmpty = s.bools[ShowEmpty]
}
