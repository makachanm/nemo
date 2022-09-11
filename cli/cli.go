package cli

type CliInterface struct {
}

func MakeCliInterface() CliInterface {
	return CliInterface{}
}

func (ci *CliInterface) Handle() {

}
