package jnal

type Config struct {
	BaseDirectory  string `mapstructure:"base_directory"`
	DateFormat     string `mapstructure:"date_format"`
	FileNameFormat string `mapstructure:"file_name_format"`
	FileTemplate   string `mapstructure:"file_template"`
	OpenCommand    string `mapstructure:"open_command"`
	ListCommand    string `mapstructure:"list_command"`
}
