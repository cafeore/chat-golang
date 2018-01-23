package trace

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちはtraceパッケージ")
		if buf.String() != "こんにちはtraceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	//fmt.Println(a...)
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
