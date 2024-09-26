package types

// Context type tracks top level debugging flag.
type Context struct {
	Debug   bool
	Success bool
}

type DebounceCommand struct {
	Quantity string
	Unit     string
	CacheDir string
	Command  []string
	Debug    bool
	Local    bool
	Status   bool
}
