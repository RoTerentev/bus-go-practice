package bus

// reader:
//		connect {stream => "out", type => $TYPE} # =>
//		msg {feeder => $NAME, type => $TYPE, msg => $STRUCT} # <=

type ReadMsg struct {
	Type   string
	Feeder string
	Msg    interface{}
}

type Reader struct {
	Type   string
	Close  func() error
	readCh chan *ReadMsg
}

func (r *Reader) ReadCh() <-chan *ReadMsg {
	return r.readCh
}
