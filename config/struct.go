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
