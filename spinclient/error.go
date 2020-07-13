package spinclient

type MalformedClientRequestError struct {
	err  string 
}

func (e *MalformedClientRequestError) Error() string {
	return e.err
}