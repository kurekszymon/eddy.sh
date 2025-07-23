package shell

type Shell interface {
	GitClone(repoURL string, cloneDir string) error
	Curl(url string, outputDir string) error
	Echo(message string) error
	Unzip(archivePath string, targetDir string) error
	RunScriptFile(filename string) error
	RunCustomScript(script string) error
	RunScriptFileInDir(filename string, dir string, args ...string) error
	GetEddyDir() (string, error)
	GetEddyBinDir() (string, error)
	CheckCommand(command string) error
	Brew(pkg string) error
	Symlink(source string, linkPath string) error
}
