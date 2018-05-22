package group_test

// import (
// 	"testing"

// 	"github.com/dc0d/club/group"
// 	"github.com/stretchr/testify/assert"
// )

// func Test01(t *testing.T) {
// 	assert := assert.New(t)
// 	g := group.New()
// 	c := g.Child()
// 	sum := 0
// 	got := make(chan int, 20)
// 	for i := 1; i <= 10; i++ {
// 		i := i
// 		sum += i
// 		c.Add(1)
// 		go func() {
// 			defer c.Done()
// 			<-c.Stopped()
// 			got <- i
// 		}()
// 	}
// 	s := g.Supervisor()
// 	s.Stop()
// 	s.Wait()
// 	close(got)
// 	sumRes := 0
// 	for v := range got {
// 		sumRes += v
// 	}
// 	assert.Equal(sum, sumRes)
// }

// func Test02(t *testing.T) {
// 	assert := assert.New(t)
// 	c := group.Child()
// 	sum := 0
// 	got := make(chan int, 20)
// 	for i := 1; i <= 10; i++ {
// 		i := i
// 		sum += i
// 		c.Add(1)
// 		go func() {
// 			defer c.Done()
// 			<-c.Stopped()
// 			got <- i
// 		}()
// 	}
// 	s := group.Supervisor()
// 	s.Stop()
// 	s.Wait()
// 	close(got)
// 	sumRes := 0
// 	for v := range got {
// 		sumRes += v
// 	}
// 	assert.Equal(sum, sumRes)
// }
