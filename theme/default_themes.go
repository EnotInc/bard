package theme

type theme_name string

const (
	bard        theme_name = "bard.json"
	catppuccin  theme_name = "catppuccin.json"
	dracula     theme_name = "dracula.json"
	nord        theme_name = "nord.json"
	gruvbox     theme_name = "gruvbox.json"
	solarized   theme_name = "solarized.json"
	tokyo_night theme_name = "tokyo_night.json"
	rose_pine   theme_name = "rose_pine.json"
	one_dark    theme_name = "one_dark.json"
)

func getDefaultTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#8c9dff",
			LineNumber:   "#9aa0a6",
			CurrentLine:  "#f5c542",
			BottomBar:    "#2d2d3f",
			Selection:    "#d4d4d4",
			Command:      "#f5c542",
			EmptyLine:    "#6fc1d4",
			Message:      "#7fc97f",
			Error:        "#e8747c",
			Tab:          "#7c9dff",
		},
		Markdown: Markdown{
			Header1:    "#7c9dff",
			Header2:    "#7c9dee",
			Header3:    "#7c9ddd",
			Header4:    "#7c9dcc",
			Header5:    "#7c9dbb",
			Header6:    "#7c9daa",
			Highlight:  "#f5e68b",
			HTMLSymbol: "#9aa0a6",
			HTMLText:   "#e8747c",
			Symbol:     "#9aa0a6",
			Quote:      "#7fc97f",
			NumberList: "#c997d9",
			Tag:        "#c997d9",
			CodeLineBg: "#2a2a3a",
			CodeHeader: "#2a2a3a",
			CodeText:   "#f5c542",
			Image:      "#6fc1d4",
			Link:       "#6fc1d4",
		},
		Code: Code{
			Background: "#2a2a3a",
			Keyword:    "#f5c542",
			String:     "#7fc97f",
			Number:     "#c997d9",
			Bracket:    "#c997d9",
			Symbol:     "#f5c542",
			Comment:    "#9aa0a6",
		},
	}
}

func getCatppuccinTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#cba6f7",
			LineNumber:   "#6c7086",
			CurrentLine:  "#f9e2af",
			BottomBar:    "#1e1e2e",
			Selection:    "#585b70",
			Command:      "#f9e2af",
			EmptyLine:    "#94e2d5",
			Message:      "#a6e3a1",
			Error:        "#f38ba8",
			Tab:          "#89b4fa",
		},
		Markdown: Markdown{
			Header1:    "#89b4fa",
			Header2:    "#74c7ec",
			Header3:    "#94e2d5",
			Header4:    "#a6e3a1",
			Header5:    "#cba6f7",
			Header6:    "#f38ba8",
			Highlight:  "#f9e2af",
			HTMLSymbol: "#6c7086",
			HTMLText:   "#f38ba8",
			Symbol:     "#6c7086",
			Quote:      "#a6e3a1",
			NumberList: "#cba6f7",
			Tag:        "#cba6f7",
			CodeLineBg: "#313244",
			CodeHeader: "#313244",
			CodeText:   "#f9e2af",
			Image:      "#94e2d5",
			Link:       "#94e2d5",
		},
		Code: Code{
			Background: "#1e1e2e",
			Keyword:    "#cba6f7",
			String:     "#a6e3a1",
			Number:     "#fab387",
			Bracket:    "#cba6f7",
			Symbol:     "#f9e2af",
			Comment:    "#6c7086",
		},
	}
}

func getDraculaTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#bd93f9",
			LineNumber:   "#44475a",
			CurrentLine:  "#f1fa8c",
			BottomBar:    "#282a36",
			Selection:    "#44475a",
			Command:      "#f1fa8c",
			EmptyLine:    "#8be9fd",
			Message:      "#50fa7b",
			Error:        "#ff5555",
			Tab:          "#bd93f9",
		},
		Markdown: Markdown{
			Header1:    "#bd93f9",
			Header2:    "#8be9fd",
			Header3:    "#50fa7b",
			Header4:    "#f1fa8c",
			Header5:    "#ff79c6",
			Header6:    "#ff5555",
			Highlight:  "#f1fa8c",
			HTMLSymbol: "#44475a",
			HTMLText:   "#ff5555",
			Symbol:     "#44475a",
			Quote:      "#50fa7b",
			NumberList: "#bd93f9",
			Tag:        "#bd93f9",
			CodeLineBg: "#44475a",
			CodeHeader: "#44475a",
			CodeText:   "#f1fa8c",
			Image:      "#8be9fd",
			Link:       "#8be9fd",
		},
		Code: Code{
			Background: "#282a36",
			Keyword:    "#ff79c6",
			String:     "#f1fa8c",
			Number:     "#bd93f9",
			Bracket:    "#ff79c6",
			Symbol:     "#f1fa8c",
			Comment:    "#6272a4",
		},
	}
}

func getNordTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#88c0d0",
			LineNumber:   "#4c566a",
			CurrentLine:  "#ebcb8b",
			BottomBar:    "#2e3440",
			Selection:    "#434c5e",
			Command:      "#ebcb8b",
			EmptyLine:    "#8fbcbb",
			Message:      "#a3be8c",
			Error:        "#bf616a",
			Tab:          "#88c0d0",
		},
		Markdown: Markdown{
			Header1:    "#88c0d0",
			Header2:    "#81a1c1",
			Header3:    "#8fbcbb",
			Header4:    "#a3be8c",
			Header5:    "#b48ead",
			Header6:    "#bf616a",
			Highlight:  "#ebcb8b",
			HTMLSymbol: "#4c566a",
			HTMLText:   "#bf616a",
			Symbol:     "#4c566a",
			Quote:      "#a3be8c",
			NumberList: "#b48ead",
			Tag:        "#b48ead",
			CodeLineBg: "#434c5e",
			CodeHeader: "#434c5e",
			CodeText:   "#ebcb8b",
			Image:      "#8fbcbb",
			Link:       "#8fbcbb",
		},
		Code: Code{
			Background: "#2e3440",
			Keyword:    "#81a1c1",
			String:     "#a3be8c",
			Number:     "#b48ead",
			Bracket:    "#81a1c1",
			Symbol:     "#ebcb8b",
			Comment:    "#4c566a",
		},
	}
}

func getGruvboxTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#d79921",
			LineNumber:   "#504945",
			CurrentLine:  "#fabd2f",
			BottomBar:    "#282828",
			Selection:    "#504945",
			Command:      "#fabd2f",
			EmptyLine:    "#83a598",
			Message:      "#b8bb26",
			Error:        "#fb4934",
			Tab:          "#83a598",
		},
		Markdown: Markdown{
			Header1:    "#83a598",
			Header2:    "#8ec07c",
			Header3:    "#b8bb26",
			Header4:    "#fabd2f",
			Header5:    "#d3869b",
			Header6:    "#fb4934",
			Highlight:  "#fabd2f",
			HTMLSymbol: "#504945",
			HTMLText:   "#fb4934",
			Symbol:     "#504945",
			Quote:      "#b8bb26",
			NumberList: "#d3869b",
			Tag:        "#d3869b",
			CodeLineBg: "#3c3836",
			CodeHeader: "#3c3836",
			CodeText:   "#fabd2f",
			Image:      "#83a598",
			Link:       "#83a598",
		},
		Code: Code{
			Background: "#282828",
			Keyword:    "#fb4934",
			String:     "#b8bb26",
			Number:     "#d3869b",
			Bracket:    "#fb4934",
			Symbol:     "#fabd2f",
			Comment:    "#928374",
		},
	}
}

func getSolarizedTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#268bd2",
			LineNumber:   "#586e75",
			CurrentLine:  "#b58900",
			BottomBar:    "#002b36",
			Selection:    "#073642",
			Command:      "#b58900",
			EmptyLine:    "#2aa198",
			Message:      "#859900",
			Error:        "#dc322f",
			Tab:          "#268bd2",
		},
		Markdown: Markdown{
			Header1:    "#268bd2",
			Header2:    "#6c71c4",
			Header3:    "#2aa198",
			Header4:    "#859900",
			Header5:    "#d33682",
			Header6:    "#dc322f",
			Highlight:  "#b58900",
			HTMLSymbol: "#586e75",
			HTMLText:   "#dc322f",
			Symbol:     "#586e75",
			Quote:      "#859900",
			NumberList: "#d33682",
			Tag:        "#d33682",
			CodeLineBg: "#073642",
			CodeHeader: "#073642",
			CodeText:   "#b58900",
			Image:      "#2aa198",
			Link:       "#2aa198",
		},
		Code: Code{
			Background: "#002b36",
			Keyword:    "#268bd2",
			String:     "#859900",
			Number:     "#d33682",
			Bracket:    "#268bd2",
			Symbol:     "#b58900",
			Comment:    "#586e75",
		},
	}
}

func getTokyoNightTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#7aa2f7",
			LineNumber:   "#3b4261",
			CurrentLine:  "#e0af68",
			BottomBar:    "#1a1b26",
			Selection:    "#33467c",
			Command:      "#e0af68",
			EmptyLine:    "#73daca",
			Message:      "#9ece6a",
			Error:        "#f7768e",
			Tab:          "#7aa2f7",
		},
		Markdown: Markdown{
			Header1:    "#7aa2f7",
			Header2:    "#0db9d7",
			Header3:    "#73daca",
			Header4:    "#9ece6a",
			Header5:    "#bb9af7",
			Header6:    "#f7768e",
			Highlight:  "#e0af68",
			HTMLSymbol: "#3b4261",
			HTMLText:   "#f7768e",
			Symbol:     "#3b4261",
			Quote:      "#9ece6a",
			NumberList: "#bb9af7",
			Tag:        "#bb9af7",
			CodeLineBg: "#283457",
			CodeHeader: "#283457",
			CodeText:   "#e0af68",
			Image:      "#73daca",
			Link:       "#73daca",
		},
		Code: Code{
			Background: "#1a1b26",
			Keyword:    "#7dcfff",
			String:     "#9ece6a",
			Number:     "#bb9af7",
			Bracket:    "#7dcfff",
			Symbol:     "#e0af68",
			Comment:    "#565f89",
		},
	}
}

func getRosePineTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#c4a7e7",
			LineNumber:   "#6e6a86",
			CurrentLine:  "#f6c177",
			BottomBar:    "#191724",
			Selection:    "#2a283e",
			Command:      "#f6c177",
			EmptyLine:    "#9ccfd8",
			Message:      "#31748f",
			Error:        "#eb6f92",
			Tab:          "#c4a7e7",
		},
		Markdown: Markdown{
			Header1:    "#c4a7e7",
			Header2:    "#9ccfd8",
			Header3:    "#31748f",
			Header4:    "#f6c177",
			Header5:    "#eb6f92",
			Header6:    "#eb6f92",
			Highlight:  "#f6c177",
			HTMLSymbol: "#6e6a86",
			HTMLText:   "#eb6f92",
			Symbol:     "#6e6a86",
			Quote:      "#31748f",
			NumberList: "#c4a7e7",
			Tag:        "#c4a7e7",
			CodeLineBg: "#2a283e",
			CodeHeader: "#2a283e",
			CodeText:   "#f6c177",
			Image:      "#9ccfd8",
			Link:       "#9ccfd8",
		},
		Code: Code{
			Background: "#191724",
			Keyword:    "#c4a7e7",
			String:     "#31748f",
			Number:     "#eb6f92",
			Bracket:    "#c4a7e7",
			Symbol:     "#f6c177",
			Comment:    "#6e6a86",
		},
	}
}

func getOneDarkTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#61afef",
			LineNumber:   "#5c6370",
			CurrentLine:  "#e5c07b",
			BottomBar:    "#282c34",
			Selection:    "#3e4451",
			Command:      "#e5c07b",
			EmptyLine:    "#56b6c2",
			Message:      "#98c379",
			Error:        "#e06c75",
			Tab:          "#61afef",
		},
		Markdown: Markdown{
			Header1:    "#61afef",
			Header2:    "#56b6c2",
			Header3:    "#98c379",
			Header4:    "#e5c07b",
			Header5:    "#c678dd",
			Header6:    "#e06c75",
			Highlight:  "#e5c07b",
			HTMLSymbol: "#5c6370",
			HTMLText:   "#e06c75",
			Symbol:     "#5c6370",
			Quote:      "#98c379",
			NumberList: "#c678dd",
			Tag:        "#c678dd",
			CodeLineBg: "#3e4451",
			CodeHeader: "#3e4451",
			CodeText:   "#e5c07b",
			Image:      "#56b6c2",
			Link:       "#56b6c2",
		},
		Code: Code{
			Background: "#282c34",
			Keyword:    "#c678dd",
			String:     "#98c379",
			Number:     "#d19a66",
			Bracket:    "#c678dd",
			Symbol:     "#e5c07b",
			Comment:    "#5c6370",
		},
	}
}
