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
	monokai     theme_name = "monokai.json"
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
			CodeLineBg: "#252538",
			CodeHeader: "#303049",
			CodeText:   "#f5c542",
			Image:      "#c997d9",
			Link:       "#6fc1d4",
		},
		Code: Code{
			Background: "#252538",
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
			SelectedTile: "#45475a",
			LineNumber:   "#6c7086",
			CurrentLine:  "#313244",
			BottomBar:    "#181825",
			Selection:    "#45475a",
			Command:      "#cba6f7",
			EmptyLine:    "#89dceb",
			Message:      "#a6e3a1",
			Error:        "#f38ba8",
			Tab:          "#89b4fa",
		},
		Markdown: Markdown{
			Header1:    "#89b4fa",
			Header2:    "#cba6f7",
			Header3:    "#89dceb",
			Header4:    "#a6e3a1",
			Header5:    "#f9e2af",
			Header6:    "#fab387",
			Highlight:  "#f9e2af",
			HTMLSymbol: "#6c7086",
			HTMLText:   "#f38ba8",
			Symbol:     "#6c7086",
			Quote:      "#a6e3a1",
			NumberList: "#cba6f7",
			Tag:        "#f38ba8",
			CodeLineBg: "#1e1e2e",
			CodeHeader: "#1e1e2e",
			CodeText:   "#cdd6f4",
			Image:      "#89dceb",
			Link:       "#89b4fa",
		},
		Code: Code{
			Background: "#1e1e2e",
			Keyword:    "#cba6f7",
			String:     "#a6e3a1",
			Number:     "#fab387",
			Bracket:    "#cdd6f4",
			Symbol:     "#89dceb",
			Comment:    "#6c7086",
		},
	}
}

func getDraculaTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#6272a4",
			LineNumber:   "#6272a4",
			CurrentLine:  "#44475a",
			BottomBar:    "#282a36",
			Selection:    "#44475a",
			Command:      "#ff79c6",
			EmptyLine:    "#8be9fd",
			Message:      "#50fa7b",
			Error:        "#ff5555",
			Tab:          "#bd93f9",
		},
		Markdown: Markdown{
			Header1:    "#ff79c6",
			Header2:    "#bd93f9",
			Header3:    "#8be9fd",
			Header4:    "#50fa7b",
			Header5:    "#f1fa8c",
			Header6:    "#ffb86c",
			Highlight:  "#f1fa8c",
			HTMLSymbol: "#6272a4",
			HTMLText:   "#ff5555",
			Symbol:     "#6272a4",
			Quote:      "#50fa7b",
			NumberList: "#bd93f9",
			Tag:        "#ff79c6",
			CodeLineBg: "#21222c",
			CodeHeader: "#21222c",
			CodeText:   "#f8f8f2",
			Image:      "#8be9fd",
			Link:       "#8be9fd",
		},
		Code: Code{
			Background: "#21222c",
			Keyword:    "#ff79c6",
			String:     "#f1fa8c",
			Number:     "#bd93f9",
			Bracket:    "#8be9fd",
			Symbol:     "#ff79c6",
			Comment:    "#6272a4",
		},
	}
}

func getNordTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#4c566a",
			LineNumber:   "#616e88",
			CurrentLine:  "#2e3440",
			BottomBar:    "#2e3440",
			Selection:    "#434c5e",
			Command:      "#88c0d0",
			EmptyLine:    "#81a1c1",
			Message:      "#a3be8c",
			Error:        "#bf616a",
			Tab:          "#8fbcbb",
		},
		Markdown: Markdown{
			Header1:    "#88c0d0",
			Header2:    "#81a1c1",
			Header3:    "#8fbcbb",
			Header4:    "#a3be8c",
			Header5:    "#ebcb8b",
			Header6:    "#d08770",
			Highlight:  "#ebcb8b",
			HTMLSymbol: "#616e88",
			HTMLText:   "#bf616a",
			Symbol:     "#616e88",
			Quote:      "#a3be8c",
			NumberList: "#b48ead",
			Tag:        "#81a1c1",
			CodeLineBg: "#242933",
			CodeHeader: "#242933",
			CodeText:   "#d8dee9",
			Image:      "#88c0d0",
			Link:       "#88c0d0",
		},
		Code: Code{
			Background: "#242933",
			Keyword:    "#81a1c1",
			String:     "#a3be8c",
			Number:     "#b48ead",
			Bracket:    "#88c0d0",
			Symbol:     "#8fbcbb",
			Comment:    "#616e88",
		},
	}
}

func getGruvboxTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#504945",
			LineNumber:   "#665c54",
			CurrentLine:  "#282828",
			BottomBar:    "#1d2021",
			Selection:    "#504945",
			Command:      "#fb4934",
			EmptyLine:    "#83a598",
			Message:      "#b8bb26",
			Error:        "#fb4934",
			Tab:          "#fabd2f",
		},
		Markdown: Markdown{
			Header1:    "#fabd2f",
			Header2:    "#83a598",
			Header3:    "#b8bb26",
			Header4:    "#d3869b",
			Header5:    "#fe8019",
			Header6:    "#8ec07c",
			Highlight:  "#fabd2f",
			HTMLSymbol: "#928374",
			HTMLText:   "#fb4934",
			Symbol:     "#928374",
			Quote:      "#b8bb26",
			NumberList: "#d3869b",
			Tag:        "#fe8019",
			CodeLineBg: "#282828",
			CodeHeader: "#282828",
			CodeText:   "#ebdbb2",
			Image:      "#83a598",
			Link:       "#83a598",
		},
		Code: Code{
			Background: "#282828",
			Keyword:    "#fb4934",
			String:     "#b8bb26",
			Number:     "#d3869b",
			Bracket:    "#ebdbb2",
			Symbol:     "#8ec07c",
			Comment:    "#928374",
		},
	}
}

func getSolarizedTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#073642",
			LineNumber:   "#586e75",
			CurrentLine:  "#002b36",
			BottomBar:    "#002b36",
			Selection:    "#073642",
			Command:      "#cb4b16",
			EmptyLine:    "#2aa198",
			Message:      "#859900",
			Error:        "#dc322f",
			Tab:          "#268bd2",
		},
		Markdown: Markdown{
			Header1:    "#268bd2",
			Header2:    "#2aa198",
			Header3:    "#6c71c4",
			Header4:    "#859900",
			Header5:    "#b58900",
			Header6:    "#cb4b16",
			Highlight:  "#b58900",
			HTMLSymbol: "#586e75",
			HTMLText:   "#dc322f",
			Symbol:     "#586e75",
			Quote:      "#859900",
			NumberList: "#6c71c4",
			Tag:        "#d33682",
			CodeLineBg: "#073642",
			CodeHeader: "#073642",
			CodeText:   "#839496",
			Image:      "#2aa198",
			Link:       "#268bd2",
		},
		Code: Code{
			Background: "#073642",
			Keyword:    "#859900",
			String:     "#2aa198",
			Number:     "#d33682",
			Bracket:    "#839496",
			Symbol:     "#268bd2",
			Comment:    "#586e75",
		},
	}
}

func getTokyoNightTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#3b4261",
			LineNumber:   "#565f89",
			CurrentLine:  "#1f2335",
			BottomBar:    "#16161e",
			Selection:    "#33467c",
			Command:      "#bb9af7",
			EmptyLine:    "#7dcfff",
			Message:      "#9ece6a",
			Error:        "#f7768e",
			Tab:          "#7aa2f7",
		},
		Markdown: Markdown{
			Header1:    "#7aa2f7",
			Header2:    "#bb9af7",
			Header3:    "#7dcfff",
			Header4:    "#9ece6a",
			Header5:    "#e0af68",
			Header6:    "#ff9e64",
			Highlight:  "#e0af68",
			HTMLSymbol: "#565f89",
			HTMLText:   "#f7768e",
			Symbol:     "#565f89",
			Quote:      "#9ece6a",
			NumberList: "#bb9af7",
			Tag:        "#f7768e",
			CodeLineBg: "#1f2335",
			CodeHeader: "#1f2335",
			CodeText:   "#a9b1d6",
			Image:      "#7dcfff",
			Link:       "#7dcfff",
		},
		Code: Code{
			Background: "#1f2335",
			Keyword:    "#bb9af7",
			String:     "#9ece6a",
			Number:     "#ff9e64",
			Bracket:    "#a9b1d6",
			Symbol:     "#7dcfff",
			Comment:    "#565f89",
		},
	}
}

func getMonokaiTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#49483e",
			LineNumber:   "#75715e",
			CurrentLine:  "#3e3d32",
			BottomBar:    "#272822",
			Selection:    "#49483e",
			Command:      "#f92672",
			EmptyLine:    "#66d9ef",
			Message:      "#a6e22e",
			Error:        "#f92672",
			Tab:          "#ae81ff",
		},
		Markdown: Markdown{
			Header1:    "#a6e22e",
			Header2:    "#66d9ef",
			Header3:    "#ae81ff",
			Header4:    "#e6db74",
			Header5:    "#fd971f",
			Header6:    "#f92672",
			Highlight:  "#e6db74",
			HTMLSymbol: "#75715e",
			HTMLText:   "#f92672",
			Symbol:     "#75715e",
			Quote:      "#a6e22e",
			NumberList: "#ae81ff",
			Tag:        "#f92672",
			CodeLineBg: "#2d2e27",
			CodeHeader: "#2d2e27",
			CodeText:   "#f8f8f2",
			Image:      "#66d9ef",
			Link:       "#66d9ef",
		},
		Code: Code{
			Background: "#2d2e27",
			Keyword:    "#f92672",
			String:     "#e6db74",
			Number:     "#ae81ff",
			Bracket:    "#f8f8f2",
			Symbol:     "#66d9ef",
			Comment:    "#75715e",
		},
	}
}

func getOneDarkTheme() *Theme {
	return &Theme{
		General: General{
			SelectedTile: "#3e4451",
			LineNumber:   "#636d83",
			CurrentLine:  "#282c34",
			BottomBar:    "#21252b",
			Selection:    "#3e4451",
			Command:      "#c678dd",
			EmptyLine:    "#56b6c2",
			Message:      "#98c379",
			Error:        "#e06c75",
			Tab:          "#61afef",
		},
		Markdown: Markdown{
			Header1:    "#61afef",
			Header2:    "#56b6c2",
			Header3:    "#c678dd",
			Header4:    "#98c379",
			Header5:    "#e5c07b",
			Header6:    "#d19a66",
			Highlight:  "#e5c07b",
			HTMLSymbol: "#636d83",
			HTMLText:   "#e06c75",
			Symbol:     "#636d83",
			Quote:      "#98c379",
			NumberList: "#c678dd",
			Tag:        "#e06c75",
			CodeLineBg: "#21252b",
			CodeHeader: "#21252b",
			CodeText:   "#abb2bf",
			Image:      "#56b6c2",
			Link:       "#61afef",
		},
		Code: Code{
			Background: "#21252b",
			Keyword:    "#c678dd",
			String:     "#98c379",
			Number:     "#d19a66",
			Bracket:    "#abb2bf",
			Symbol:     "#56b6c2",
			Comment:    "#5c6370",
		},
	}
}
