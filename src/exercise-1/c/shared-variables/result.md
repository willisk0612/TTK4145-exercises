### Without synronization
The result was 0 before introducing usleep() calls. This is most likely because one of the threads started and finished before the other. Once implementing it, the deviation was in the range +-1000.

### With syncronization
Adding a mutex lock fixed the issue. It prevents both threads from accessing the critical section at the same time.
