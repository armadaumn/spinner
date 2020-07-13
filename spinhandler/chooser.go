package spinhandler

type Chooser interface {
	F(ClientMap) (string, error)
}