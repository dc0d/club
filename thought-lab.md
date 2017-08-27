# 2017-08-27
[errgroup](https://github.com/golang/sync/blob/master/errgroup/errgroup.go):

* **context**: for (hierarchial) cancelation
* **wait group**: wait (confirm) goroutines stopped

[group](https://github.com/oklog/oklog/blob/master/pkg/group/group.go):

* notify other to stop
* notify others *why* (the error) they have to stop

Is using contexts, always a good choice? What if (OK, what-ifs are bad) child goroutines should ignore the upper context cancelation and wait for the parent to explicitly stop them?
