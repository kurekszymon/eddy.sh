package types

type Git struct {
	CloneDir string   `yaml:"clone_dir"`
	Repos    []string `yaml:"repos"`
}
