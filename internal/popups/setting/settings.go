package setting

type Settings struct {
	onChange func()
	w, h     int
	cursor   int
}

const header_offset = 1
const settings_amount int = 7

type setting int

const (
	_ setting = iota
	RelativeNumbers
	ShowMDSymbols
	ShowTabNames
	EnableRender
	ShowDotFiles
	ShowBorders
	ShowIcons
	ShowEmpty
)

func IntiSettings(onChange func()) *Settings {
	s := &Settings{
		onChange: onChange,
	}
	return s
}
