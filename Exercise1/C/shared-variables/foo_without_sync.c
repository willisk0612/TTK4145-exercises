// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>
#include <unistd.h> // for usleep

int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    for (int counter = 0;counter<1000000;counter++) {
        i++;
        usleep(1); // to demonstrate a race condition
    }
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    for (int counter = 0;counter<1000000;counter++) {
        i--;
        usleep(1); // to demonstrate a race condition
    }
    return NULL;
}


int main(){
    // TODO:
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_t thread_1;
    pthread_t thread_2;

    // Create thread 1 (incrementing function)
    pthread_create(&thread_1, NULL, incrementingThreadFunction, NULL);
    // Create thread 2 (decrementing function)
    pthread_create(&thread_2, NULL, decrementingThreadFunction, NULL);


    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`
    pthread_join(thread_1, NULL); // joins the first thread with NULL as return value
    pthread_join(thread_2, NULL); // joins the second thread with NULL as return value

    printf("The magic number is: %d\n", i);
    return 0;
}
