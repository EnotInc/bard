package config

// About |Config|
// |RLN| - relative line munber
// |ShowMD| - always show markdown symbols
// |Render| - enable or disable render
// |TabNames| - show file name is tabs or not
type Config struct {
	RLN       bool   `json:"relative_line_numbers"`
	ShowMD    bool   `json:"show_markdown_symbols"`
	Render    bool   `json:"enable_render"`
	TabNames  bool   `json:"show_tab_names"`
	ThemeName string `json:"theme_name"`
}

type Theme struct {
	General  General  `json:"general"`
	Markdown Markdown `json:"markdown"`
	Code     Code     `json:"code"`
}

type General struct {
	CurrentLine string `json:"current_line"`
	LineNumber  string `json:"line_number"`
	BottomBar   string `json:"bottom_bar"`
	EmptyLine   string `json:"empty_line"`
	Selection   string `json:"selection"`
	Command     string `json:"command"`
	Message     string `json:"message"`
	Tab         string `json:"tab"`
}

type Markdown struct {
	NumberList string `json:"number_list"`
	Highlight  string `json:"highlight"`
	CodeLineBg string `json:"code_line_bg"`
	CodeHeader string `json:"code_header"`
	CodeText   string `json:"code_text"`
	Header1    string `json:"header_1"`
	Header2    string `json:"header_2"`
	Header3    string `json:"header_3"`
	Header4    string `json:"header_4"`
	Header5    string `json:"header_5"`
	Header6    string `json:"header_6"`
	Symbol     string `json:"symbol"`
	Quote      string `json:"quote"`
	Image      string `json:"image"`
	Link       string `json:"link"`
	Tag        string `json:"tag"`
}

type Code struct {
	Background string `json:"background"`
	Keyword    string `json:"keyword"`
	Bracket    string `json:"bracket"`
	Comment    string `json:"comment"`
	String     string `json:"string"`
	Number     string `json:"number"`
	Symbol     string `json:"symbol"`
}
