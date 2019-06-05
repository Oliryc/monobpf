#include <stdio.h>
#include <unistd.h>

int auie() {
  printf("Sleeping…\n");
  sleep(20);
  printf("Done sleeping.\n");
}

int main()
{
  for(;;)
    auie();
  return 0;
}
