package theme

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
