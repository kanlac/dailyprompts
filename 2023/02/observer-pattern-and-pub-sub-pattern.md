### 观察者模式与发布订阅模式


观察者模式和发布订阅模式都是一对多的生产消费模式，当发布者发布一条消息时，订阅者们会根据订阅主题来选择取用它们各自需要的消息。它们实现了消息发布方和订阅方的解藕，订阅者和发布者都可以运行时动态添加，实际使用场景有天气预报，或者通知系统中另一个组件更新某些数据缓存等等。

两者的基本概念类似，但观察者模式一般是指通过函数调用来传递数据，而发布订阅模式一般指通过消息队列（golang 中的 channel）传递。

发布订阅模式代码：
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Topic func(v interface{}) bool

type Subscribers map[chan interface{}]Topic // tip: channel is comparable, lambda is not

type Publisher struct {
	subscribers Subscribers
	m           sync.Mutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: make(Subscribers),
	}
}

func (p *Publisher) SubscribeAll() <-chan interface{} {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) SubscribeTopic(topic Topic) <-chan interface{} {
	p.m.Lock()
	defer p.m.Unlock()

	ch := make(chan interface{})
	p.subscribers[ch] = topic
	return ch
}

func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	var wg sync.WaitGroup
	for ch, topic := range p.subscribers {
		if topic != nil && !topic(v) {
			continue
		}
		go publishToCertainSubscriber(v, ch, &wg)
	}
	wg.Wait()
}

func publishToCertainSubscriber(v interface{}, subscriber chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	subscriber <- v
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for ch := range p.subscribers {
		close(ch)
	}
}

func printMessage(in <-chan interface{}) {
	for v := range in {
		fmt.Println("printMessage: ", v)
	}
}

func squareAndPrint(in <-chan interface{}) {
	for v := range in {
		n := v.(int)
		fmt.Println("squeareAndPrint: ", n*n)
	}
}

func main() {
	/* set up */

	p := NewPublisher()

	s1 := p.SubscribeAll()

	numberOnly := func(v interface{}) bool {
		_, ok := v.(int)
		return ok
	}
	s2 := p.SubscribeTopic(numberOnly)

	go printMessage(s1)
	go squareAndPrint(s2)

	/* start messaging */

	p.Publish("hello world")
	p.Publish(25)
	defer p.Close()

	time.Sleep(time.Second * 10)
}
```