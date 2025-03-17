#include "cgomemprof.h"
#include <stdbool.h>
#include <stdlib.h>
#include <backtrace.h>

int EnableMemoryProfiling() {
  bool enable = true;
  return mallctl("prof.active", NULL, NULL, &enable, sizeof(enable));
}

int DisableMemoryProfiling() {
  bool enable = false;
  return mallctl("prof.active", NULL, NULL, &enable, sizeof(enable));
}

int DumpMemoryProfileIntoFile(const char *filename) {
  return mallctl("prof.dump", NULL, NULL, &filename, sizeof(const char *));
}

void syminfoCallback(void *data, uintptr_t pc, const char *symname, uintptr_t symval, uintptr_t symsize) {
    char* buf = (char*)data;
    if (symname != NULL) {
        snprintf(buf, 4096, "0x%lx\t%s",pc, symname);
    } else {
        snprintf(buf, 4096, "0x%lx\t??", pc);
    }
}

void errorCallback(void *data, const char *msg, int errnum) {
    fprintf(stderr, "Error: %s (code: %d)\n", msg, errnum);
}

char* GetSymbol(uintptr_t addr) {
    struct backtrace_state *state = backtrace_create_state(NULL, 1, errorCallback, NULL);
    char* buf = malloc(4096);
    backtrace_syminfo(state, addr, syminfoCallback, errorCallback, buf);
    return buf;
}
