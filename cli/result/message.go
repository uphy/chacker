package result

import (
	"bytes"
	"fmt"
	"io"
)

type MessageResultBody struct {
	Out *bytes.Buffer
}

func NewMessageResultBody() *MessageResultBody {
	return &MessageResultBody{new(bytes.Buffer)}
}

func (s *MessageResultBody) JSON() interface{} {
	return nil
}

func (s *MessageResultBody) Pretty(writer io.Writer) error {
	_, err := fmt.Fprint(writer, s.Out.String())
	return err
}

func (s *MessageResultBody) Plain(writer io.Writer) error {
	return s.Pretty(writer)
}
