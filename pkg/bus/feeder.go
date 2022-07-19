package bus

// feeder:
//     connect {stream => "in", name => $STR} # =>
//     msg {type => $TYPE, msg => $STRUCT} # =>
type FeedMsg struct {
	Type string
	Msg  interface{}
}

type Feeder struct {
	Name   string
	Close  func() error
	feedCh chan *FeedMsg
}

func (f *Feeder) Feed() chan<- *FeedMsg {
	return f.feedCh
}
