package commands

type Command interface {
	Parse(args []string) error
	Validate() error
	Exec()
}
