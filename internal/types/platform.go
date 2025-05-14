package types

type PkgManager string

const (
	Brew   PkgManager = "brew"
	Manual PkgManager = "manual"
)
