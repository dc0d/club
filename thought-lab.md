# 2017-08-28
`errgroup` has this advantage that goroutines can get added dynamically. It uses a context that should be used inside those registered goroutines (and for backward signaling). The original error is available as the return value of `Wait()`, but it can not be used from inside the workers (because of `WaitGroup`). We adapt this by adding the recovery from panic and study it for some time.

And it should be possible to use the controlling/signaling tools inside a struct which normally is the (shared) **closure for those goroutines**.

# 2017-08-27
[errgroup](https://github.com/golang/sync/blob/master/errgroup/errgroup.go):

* **context**: for (hierarchial) cancelation
* **wait group**: wait (confirm) goroutines stopped

[group](https://github.com/oklog/oklog/blob/master/pkg/group/group.go):

* notify other to stop
* notify others *why* (the error) they have to stop

Is using contexts, always a good choice? What if (OK, what-ifs are bad) child goroutines should ignore the upper context cancelation and wait for the parent to explicitly stop them?

There is this pile of goroutines that work as a unit. We need:

* cancellation: which means either 
    * an error happened (from inside the pile) or 
    * just stop (from outside the pile)
* wait for all goroutines to finish
