package queue_test

import (
	"sync"

	"github.com/unchartedsoftware/veldt/util/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testRequest struct {
	count int
	mu    *sync.Mutex
}

func newTestRequest() *testRequest {
	return &testRequest{
		mu: &sync.Mutex{},
	}
}

func (r *testRequest) Create() ([]byte, error) {
	r.mu.Lock()
	r.count++
	r.mu.Unlock()
	return nil, nil
}

func (r *testRequest) Decrement() {
	r.mu.Lock()
	r.count--
	r.mu.Unlock()
}

func (r *testRequest) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

var _ = Describe("Queue", func() {

	var q *queue.Queue

	BeforeEach(func() {
		q = queue.NewQueue()
	})

	Describe("Send", func() {

		It("should never execute more than max concurrent at the same time", func() {
			req := newTestRequest()
			n := 8
			m := 256
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
			for i := 0; i < n; i++ {
				q.Send(newTestRequest())
			}
		})

	})
})
