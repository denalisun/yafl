#ifndef UTILS_H
#define UTILS_H

#include <windows.h>
#include <stdbool.h>

int NtSuspendProcess(HANDLE ProcessHandle);

char *combine_path(const char* p1, const char* p2);
HANDLE start_process(char* processName, char* args, char* directory);
BOOL file_exists(LPCTSTR szPath);

#endif