package bus

import (
	"bus-practice/internal"
	"context"
	"errors"
	"sync"
)

type Bus struct {
	sync.Mutex
	readers  map[string]map[*Reader]struct{}
	messages *internal.Queue
}

type ConnectOpts struct {
	Stream string
	Name   string
	Type   string
}

func NewBus(historySize int) *Bus {
	return &Bus{
		messages: internal.NewQueue(historySize),
		readers:  make(map[string]map[*Reader]struct{}),
	}
}

func (b *Bus) Connect(opts ConnectOpts) (interface{}, error) {
	switch opts.Stream {
	case "in":
		return b.addFeader(opts)
	case "out":
		return b.addReader(opts)
	default:
		return nil, errors.New("opts.Stream must has one of values: \"in\", \"out\"")
	}
}

func (b *Bus) addFeader(opts ConnectOpts) (*Feeder, error) {
	ctx, cancel := context.WithCancel(context.TODO())

	f := &Feeder{
		Name:   opts.Name,
		feedCh: make(chan *FeedMsg),
	}

	f.Close = func() error {
		close(f.feedCh)
		cancel()
		return nil
	}

	go func() {
		for {
			select {
			case msg := <-f.feedCh:
				b.feedIn(&ReadMsg{
					Type:   msg.Type,
					Msg:    msg.Msg,
					Feeder: opts.Name,
				})
			case <-ctx.Done():
				return
			}
		}
	}()

	return f, nil
}

func (b *Bus) addReader(opts ConnectOpts) (*Reader, error) {
	b.Lock()
	defer b.Unlock()
	rdrs, ok := b.readers[opts.Type]
	if !ok {
		b.readers[opts.Type] = map[*Reader]struct{}{}
		rdrs = b.readers[opts.Type]
	}
	r := &Reader{
		Type:   opts.Type,
		readCh: make(chan *ReadMsg),
	}

	rdrs[r] = struct{}{}
	r.Close = func() error {
		delete(rdrs, r)
		close(r.readCh)
		return nil
	}
	go b.sendLastMessages(r)
	return r, nil
}

func (b *Bus) feedIn(msg *ReadMsg) {

	go func() {
		b.messages.In() <- msg
	}()

	rdrs, readersExist := b.readers[msg.Type]
	if readersExist {
		for reader := range rdrs {
			go func(r *Reader, message *ReadMsg) {
				r.readCh <- &ReadMsg{
					Type:   message.Type,
					Msg:    message.Msg,
					Feeder: message.Feeder,
				}
			}(reader, msg)
		}
	}
}

func (b *Bus) sendLastMessages(r *Reader) {
	for _, m := range b.messages.Messages() {
		if m, ok := m.(*ReadMsg); ok && m.Type == r.Type {
			go func() { r.readCh <- m }()
		}
	}
}
