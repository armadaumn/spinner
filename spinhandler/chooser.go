package spinhandler

type ChooseFunction func (ClientMap) (string, error)

type Chooser interface {
	F		ChooseFunction
}