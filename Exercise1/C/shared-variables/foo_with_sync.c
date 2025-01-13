// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>
#include <unistd.h> // for usleep

int i = 0;

// Note the return type: void*

pthread_mutex_t lock;

void* increment(){
    // TODO: increment i 1_000_000 times
    for (int counter = 0;counter<1000000;counter++) {
        pthread_mutex_lock(&lock);
        i++;
        pthread_mutex_unlock(&lock);
        usleep(1);
    }
    return NULL;
}

void* decrement(){
    // TODO: decrement i 1_000_000 times
    for (int counter = 0;counter<1000000;counter++) {
        pthread_mutex_lock(&lock);
        i--;
        pthread_mutex_unlock(&lock);
        usleep(1);
    }
    return NULL;
}


int main(){
    // TODO:
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_t thread_1;
    pthread_t thread_2;

    pthread_create(&thread_1, NULL, increment, NULL);
    pthread_create(&thread_2, NULL, decrement, NULL);


    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`
    pthread_join(thread_1, NULL);
    pthread_join(thread_2, NULL);

    printf("The magic number is: %d\n", i);
    return 0;
}
