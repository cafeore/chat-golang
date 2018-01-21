package trace

//Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェースです。
type Tracer interface {
	Trace(...interface{})
}
