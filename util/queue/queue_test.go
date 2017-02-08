package queue_test

import (
	"sync"
	"time"

	"github.com/unchartedsoftware/veldt/util/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type countRequest struct {
	count int
	mu    *sync.Mutex
}

func newTestRequest() *countRequest {
	return &countRequest{
		mu: &sync.Mutex{},
	}
}

func (r *countRequest) Create() ([]byte, error) {
	r.mu.Lock()
	r.count++
	r.mu.Unlock()
	return nil, nil
}

func (r *countRequest) Decrement() {
	r.mu.Lock()
	r.count--
	r.mu.Unlock()
}

func (r *countRequest) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

type sleepRequest struct {
}

func newSleepRequest() *sleepRequest {
	return &sleepRequest{}
}

func (r *sleepRequest) Create() ([]byte, error) {
	time.Sleep(time.Millisecond * 200)
	return nil, nil
}

var _ = Describe("Queue", func() {

	var q *queue.Queue
	var m int

	BeforeEach(func() {
		q = queue.NewQueue()
		m = 256
	})

	Describe("Send", func() {

		It("should execute the requests when there is availability", func() {
			req := newTestRequest()
			for i := 0; i < m; i++ {
				_, err := q.Send(req)
				Expect(err).To(BeNil())
			}
			count := req.Count()
			Expect(count).To(Equal(m))
		})

	})

	Describe("SetMaxConcurrent", func() {

		It("should set the maximum number of concurrent requests", func() {
			req := newTestRequest()
			n := 8
			q.SetMaxConcurrent(n)
			for i := 0; i < m; i++ {
				_, err := q.Send(req)
				Expect(err).To(BeNil())
				count := req.Count()
				Expect(count).To(BeNumerically("<=", n))
				req.Decrement()
			}
			count := req.Count()
			Expect(count).To(Equal(0))
		})

		It("should handle setting above the current max number", func() {
			req := newTestRequest()
			n := 64
			q.SetMaxConcurrent(n)
			for i := 0; i < m; i++ {
				_, err := q.Send(req)
				Expect(err).To(BeNil())
				count := req.Count()
				Expect(count).To(BeNumerically("<=", n))
				req.Decrement()
			}
			count := req.Count()
			Expect(count).To(Equal(0))
		})

		It("should handle setting below the current max number", func() {
			req := newTestRequest()
			n := 8
			q.SetMaxConcurrent(n)
			for i := 0; i < m; i++ {
				_, err := q.Send(req)
				Expect(err).To(BeNil())
				count := req.Count()
				Expect(count).To(BeNumerically("<=", n))
				req.Decrement()
			}
			count := req.Count()
			Expect(count).To(Equal(0))
		})

	})

	Describe("SetLength", func() {

		It("should set the queue length", func() {
			n := 20
			q.SetLength(n)
			var err error
			c := make(chan error)
			for i := 0; i < m; i++ {
				select {
				case err = <-c:
					break
				default:
					go func() {
						_, err = q.Send(newSleepRequest())
						if err != nil {
							c <- err
						}
					}()
				}

			}
			Expect(err).ToNot(BeNil())
		})

	})
})
