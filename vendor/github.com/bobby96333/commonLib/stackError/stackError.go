package stackError

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
)

type StackError struct {
	stack []byte
	msg   string
}

func New(msg string) *StackError {
	return &StackError{
		stack: debug.Stack(),
		msg:   msg,
	}

}

func (ths *StackError) PrintErr() {
	fmt.Fprintln(os.Stderr, "error:", ths.msg)
	fmt.Fprintln(os.Stderr, string(ths.stack))
}

func NewFromError(err error, stackErrorId int) *StackError {
	if err == nil {
		return nil
	}
	return &StackError{
		stack: debug.Stack(),
		msg:   err.Error() + ",StackError id is" + strconv.Itoa(stackErrorId),
	}
}
