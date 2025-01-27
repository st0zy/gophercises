package mapping

type PathMapping struct {
	Path         string `yaml:"path" json:"path"`
	RedirectPath string `yaml:"url" json:"url"`
}
