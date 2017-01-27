package promise_test

import (
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/unchartedsoftware/veldt/util/promise"
)

var _ = Describe("map", func() {

	const (
		numConcurrent = 256 * 8
		numKeys       = 36
	)

	Describe("Set", func() {
		It("should set a promise into the map under the provided key", func() {
			m := promise.NewMap()
			p := promise.NewPromise()
			m.Set("test", p)
			err := fmt.Errorf("This is an error")
			p.Resolve(err)
			o, ok := m.Get("test")
			Expect(ok).To(Equal(true))
			Expect(p.Wait()).To(Equal(err))
			Expect(o.Wait()).To(Equal(err))
		})
		It("should be thread safe", func() {
			wg := sync.WaitGroup{}
			m := promise.NewMap()
			for i := 0; i < numConcurrent; i++ {
				wg.Add(1)
				go func(index int) {
					key := fmt.Sprintf("%d", index)
					m.Set(key, promise.NewPromise())
					wg.Done()
				}(i % numKeys)
			}
			wg.Wait()
		})
	})

	Describe("Get", func() {
		It("should return a promise under the provided key", func() {
			m := promise.NewMap()
			p := promise.NewPromise()
			m.Set("test", p)
			err := fmt.Errorf("This is an error")
			p.Resolve(err)
			o, ok := m.Get("test")
			Expect(ok).To(Equal(true))
			Expect(p).To(Equal(o))
			Expect(p.Wait()).To(Equal(err))
			Expect(o.Wait()).To(Equal(err))
		})
		It("should return false if the promise does not exist", func() {
			m := promise.NewMap()
			_, ok := m.Get("test")
			Expect(ok).To(Equal(false))
		})
		It("should return true if the promise exists", func() {
			m := promise.NewMap()
			p := promise.NewPromise()
			m.Set("test", p)
			_, ok := m.Get("test")
			Expect(ok).To(Equal(true))
		})
		It("should be thread safe", func() {
			wg := sync.WaitGroup{}
			m := promise.NewMap()
			p := promise.NewPromise()
			err := fmt.Errorf("This is an error")
			m.Set("test", p)
			for i := 0; i < numConcurrent; i++ {
				wg.Add(1)
				go func() {
					p, ok := m.Get("test")
					Expect(ok).To(Equal(true))
					Expect(p.Wait()).To(Equal(err))
					wg.Done()
				}()
			}
			p.Resolve(err)
			wg.Wait()
		})
	})

	Describe("GetOrCreate", func() {
		It("should always return a promise", func() {
			m := promise.NewMap()
			p, _ := m.GetOrCreate("test")
			Expect(p).To(Not(BeNil()))
		})
		It("should return true if the promise already existed", func() {
			m := promise.NewMap()
			p := promise.NewPromise()
			m.Set("test", p)
			o, ok := m.GetOrCreate("test")
			Expect(ok).To(Equal(true))
			Expect(p).To(Equal(o))
		})
		It("should return false if the promise was created", func() {
			m := promise.NewMap()
			_, ok := m.GetOrCreate("test")
			Expect(ok).To(Equal(false))
		})
		It("should be thread safe", func() {
			wg := sync.WaitGroup{}
			m := promise.NewMap()
			err := fmt.Errorf("This is an error")
			for i := 0; i < numConcurrent; i++ {
				wg.Add(1)
				go func() {
					p, ok := m.GetOrCreate("test")
					if !ok {
						// creator will resolve it
						time.Sleep(time.Millisecond * 100)
						p.Resolve(err)
					} else {
						Expect(p.Wait()).To(Equal(err))
					}
					wg.Done()
				}()
			}
			wg.Wait()
		})
	})

	Describe("Remove", func() {
		It("should remove a promise from under a provided key", func() {
			m := promise.NewMap()
			p := promise.NewPromise()
			m.Set("test", p)
			_, ok := m.Get("test")
			Expect(ok).To(Equal(true))
			m.Remove("test")
			_, ok = m.Get("test")
			Expect(ok).To(Equal(false))
		})
		It("should be thread safe", func() {
			wg := sync.WaitGroup{}
			m := promise.NewMap()
			p := promise.NewPromise()
			err := fmt.Errorf("This is an error")
			m.Set("test", p)
			for i := 0; i < numConcurrent; i++ {
				wg.Add(1)
				go func() {
					m.Remove("test")
					wg.Done()
				}()
			}
			p.Resolve(err)
			wg.Wait()
		})
	})

})
