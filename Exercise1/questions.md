Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is mostly used for multithreading while parallellism is mostly used for multiprocessing.

What is the difference between a *race condition* and a *data race*?
> Race condition is a more general term used when multiple workds use a shared variable, and the result is unpredictable because of unsynchronized read/write operations. Data race is a subset of a race condition where at least one of the operations is a write.

*Very* roughly - what does a *scheduler* do, and how does it do it?
> A scheduler is a mechanism that controls processes and threads. It does it by using specific algorithms designed for the specific task (FCFS, round robin etc).


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> Using multiple threads is an effective way of utilizing all the resources of a CPU. It can achieve things much faster by using multiple workers (threads) for a specific task.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Green threads are implemented at the application level rather than OS level. They typically use less resources because of less overhead. We would use them for concurrent tasks that do not require heavy CPU processing.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> No, because real time systems have strict requirements. For non concurrent programming we are mostly concerned with creating a working code, without thinking much about time deadlines and synchronization. A real time programmer has to think of every aspect of the code.

What do you think is best - *shared variables* or *message passing*?
> It depends on the situation. But alot of the time message passing can be seen as better, because shared variables can cause issues like race conditions and deadlock. It is simpler to implement a shared variable, so what is best depends on the specific situation.

