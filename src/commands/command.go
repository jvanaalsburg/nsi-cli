package commands

type Command interface {
	Parse(args []string)
	Validate() error
	Exec()
}
