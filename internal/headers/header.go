package headers

type Header map[string]string

func NewHeader() Header {
	return make(Header)
}

func (h Header) Parse(data []byte) (n int, done bool, err error) {
	return 0, false, nil
}
