package group

// import (
// 	"sync"
// )

// //-----------------------------------------------------------------------------

// // Group groups a set of actions (goroutines) and signals the stop
// type Group struct {
// 	wg       *sync.WaitGroup
// 	stop     chan struct{}
// 	stopOnce *sync.Once
// }

// // New .
// func New() *Group {
// 	return &Group{
// 		wg:       &sync.WaitGroup{},
// 		stop:     make(chan struct{}),
// 		stopOnce: &sync.Once{},
// 	}
// }

// // Add .
// func (g *Group) Add(delta int) {
// 	g.wg.Add(delta)
// }

// // Done .
// func (g *Group) Done() { g.wg.Done() }

// // Stopped .
// func (g *Group) Stopped() <-chan struct{} { return g.stop }

// // Wait .
// func (g *Group) Wait() { g.wg.Wait() }

// // Stop .
// func (g *Group) Stop() { g.stopOnce.Do(func() { close(g.stop) }) }

// // Supervisor returns the supervisor API
// func (g *Group) Supervisor() interface {
// 	Wait()
// 	Stop()
// } {
// 	return g
// }

// // Child returns the child API
// func (g *Group) Child() interface {
// 	Add(delta int)
// 	Done()
// 	Stopped() <-chan struct{}
// } {
// 	return g
// }

// //-----------------------------------------------------------------------------

// var _g *Group

// func init() { _g = New() }

// // Supervisor .
// func Supervisor() interface {
// 	Wait()
// 	Stop()
// } {
// 	return _g
// }

// // Child .
// func Child() interface {
// 	Add(delta int)
// 	Done()
// 	Stopped() <-chan struct{}
// } {
// 	return _g
// }

// //-----------------------------------------------------------------------------
