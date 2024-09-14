package types

// Context type tracks top level debugging flag.
type Context struct {
	Debug   bool
	Success bool
}

type DebounceCommand struct {
	Quantity string
	Unit     string
	Command  []string
	Debug    bool
	Status   bool
}
