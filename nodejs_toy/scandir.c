#define _GNU_SOURCE

#include <dlfcn.h>
#include <unistd.h>
#include <dirent.h>
#include <stdio.h>

int scandir64(const char * dirp, struct dirent64 ** * namelist, int( * filter)(const struct dirent64 * ), int( * compar)(const struct dirent64 ** ,
    const struct dirent64 ** )) {
    int( * real_scandir)(const char * dirp, struct
        dirent64 ** * namelist, int( * filter)(const struct dirent64 * ), int( * compar)(const struct dirent64 ** ,
            const struct dirent64 ** ));
    real_scandir = dlsym(RTLD_NEXT, "scandir64");

    printf("Slow system call waiting...\n");
    sleep(1);
    printf("Slow system call wait finished\n");
    return real_scandir(dirp, namelist, filter, compar);
}
