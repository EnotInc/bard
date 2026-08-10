package config

type Config struct {
	ThemeName  string `json:"theme_name,omitempty"`
	ResizeTime int    `json:"resize_time_duration,omitempty"`
	TabStop    int    `json:"tab_stop,omitempty"`
	RLN        bool   `json:"relative_line_numbers"`
	ShowEmpty  bool   `json:"show_empty_line_symbol"`
	ShowMD     bool   `json:"show_markdown_symbols"`
	Render     bool   `json:"enable_render"`
	TabNames   bool   `json:"show_tab_names"`
	KeepTabs   bool   `json:"keep_tabs"`
	ShowIcons  bool   `json:"show_icons"`
	ShowBorder bool   `json:"show_borders"`
	ShowDot    bool   `json:"show_dot_files"`
}

type Theme struct {
	General  General  `json:"general"`
	Markdown Markdown `json:"markdown"`
	Code     Code     `json:"code"`
}

type General struct {
	SelectedTile string `json:"selected_tile" type:"foreground"`
	CurrentLine  string `json:"current_line" type:"foreground"`
	LineNumber   string `json:"line_number" type:"foreground"`
	BottomBar    string `json:"bottom_bar" type:"background"`
	EmptyLine    string `json:"empty_line" type:"foreground"`
	Selection    string `json:"selection" type:"background"`
	Command      string `json:"command" type:"foreground"`
	Message      string `json:"message" type:"foreground"`
	Error        string `json:"error" type:"foreground"`
	Tab          string `json:"tab" type:"foreground"`
}

type Markdown struct {
	NumberList string `json:"number_list" type:"foreground"`
	Highlight  string `json:"highlight" type:"background"`
	CodeLineBg string `json:"code_line_bg" type:"background"`
	CodeHeader string `json:"code_header" type:"background"`
	HTMLSymbol string `json:"html_tag_symbol" type:"foreground"`
	HTMLText   string `json:"html_tag_text" type:"foreground"`
	CodeText   string `json:"code_text" type:"foreground"`
	Header1    string `json:"header_1" type:"foreground"`
	Header2    string `json:"header_2" type:"foreground"`
	Header3    string `json:"header_3" type:"foreground"`
	Header4    string `json:"header_4" type:"foreground"`
	Header5    string `json:"header_5" type:"foreground"`
	Header6    string `json:"header_6" type:"foreground"`
	Symbol     string `json:"symbol" type:"foreground"`
	Quote      string `json:"quote" type:"foreground"`
	Image      string `json:"image" type:"foreground"`
	Link       string `json:"link" type:"foreground"`
	Tag        string `json:"tag" type:"foreground"`
}

type Code struct {
	Background string `json:"background" type:"background"`
	Keyword    string `json:"keyword" type:"foreground"`
	Bracket    string `json:"bracket" type:"foreground"`
	Comment    string `json:"comment" type:"foreground"`
	String     string `json:"string" type:"foreground"`
	Number     string `json:"number" type:"foreground"`
	Symbol     string `json:"symbol" type:"foreground"`
}
