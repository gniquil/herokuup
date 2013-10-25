package herokuup

type Response struct {
	url    string
	status int
}

func (res *Response) failed() bool {
	return res.status != 200
}
