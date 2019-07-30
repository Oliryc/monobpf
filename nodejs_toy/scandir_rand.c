#define _GNU_SOURCE

#include <dlfcn.h>
#include <unistd.h>
#include <dirent.h>
#include <stdio.h>

#include <sys/time.h>
#include <stdlib.h>

// Probability to be slow
#define SLOW_PROBA   0.1
// Duration of the sleep time in the middle of the call (μsec)
#define SLEEP_DURATION   1500*1000

// Randomly slow scandir
int scandir64(const char * dirp, struct dirent64 ** * namelist, int( * filter)(const struct dirent64 * ), int( * compar)(const struct dirent64 ** ,
    const struct dirent64 ** )) {
    int( * real_scandir)(const char * dirp, struct
        dirent64 ** * namelist, int( * filter)(const struct dirent64 * ), int( * compar)(const struct dirent64 ** ,
            const struct dirent64 ** ));
    real_scandir = dlsym(RTLD_NEXT, "scandir64");

    struct timeval tv;
    gettimeofday(&tv,NULL);
    srand(tv.tv_usec); // Should be a bit different even for two consecutive call
    int r = rand();
    if (r < RAND_MAX*SLOW_PROBA) {
      printf("Slow system call waiting...\n");
      usleep(SLEEP_DURATION);
      printf("Slow system call wait finished\n");
    } else {
      printf("Normal system call\n");
    }
    return real_scandir(dirp, namelist, filter, compar);
}
