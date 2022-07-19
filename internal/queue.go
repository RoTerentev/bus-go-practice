package internal

type Queue struct {
	msgs  []interface{}
	input chan interface{}
}

func NewQueue(size int) *Queue {
	q := Queue{
		msgs:  make([]interface{}, size),
		input: make(chan interface{}),
	}
	go func(queue *Queue) {
		for msg := range q.input {
			q.addMsg(msg)
		}
	}(&q)

	return &q
}

func (q *Queue) In() chan<- interface{} {
	return q.input
}

func (q *Queue) Messages() []interface{} {
	return q.msgs
}

func (q *Queue) addMsg(m interface{}) {
	if len(q.msgs) == cap(q.msgs) {
		q.msgs = q.msgs[1:]
	}
	q.msgs = append(q.msgs, m)
}
