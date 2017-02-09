package queue_test

import (
	"sync"

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

type pauseRequest struct {
	c chan bool
}

func newPauseRequest() *pauseRequest {
	return &pauseRequest{
		c: make(chan bool),
	}
}

func (r *pauseRequest) Create() ([]byte, error) {
	<-r.c
	return nil, nil
}

func (r *pauseRequest) Unpause() {
	r.c <- true
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

		It("should set the queue length, returning an error when surpassed", func() {
			n := 20
			p := 20
			q.SetMaxConcurrent(p)
			q.SetLength(n)
			reqs := make([]*pauseRequest, m)
			mu := &sync.Mutex{}
			wg := &sync.WaitGroup{}
			wg.Add(m - n - p - 1)
			errCount := 0
			// asynchronously send m requests
			for i := 0; i < m; i++ {
				go func(index int) {
					req := newPauseRequest()
					reqs[index] = req
					_, err := q.Send(req)
					if err != nil {
						mu.Lock()
						errCount++
						reqs[index] = nil
						mu.Unlock()
						wg.Done()
					}
				}(i)
			}
			// ensure all err'd requests have returned at this point
			wg.Wait()
			// unpause all paused requests
			for _, req := range reqs {
				if req != nil {
					go func(req *pauseRequest) {
						req.Unpause()
					}(req)
				}
			}
			Expect(errCount).To(BeNumerically("==", m-n-p-1))
		})

	})
})
