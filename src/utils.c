#include "utils.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

typedef NTSTATUS(NTAPI *pNtSuspendProcess)(HANDLE ProcessHandle);
int NtSuspendProcess(HANDLE ProcessHandle) {
    HMODULE hNtDll = GetModuleHandleA("ntdll.dll");
    if (hNtDll == NULL) {
        printf("Failed to find ntdll.dll!\n");
        return 1;
    }

    pNtSuspendProcess NtSuspendProcessFunc = (pNtSuspendProcess)GetProcAddress(hNtDll, "NtSuspendProcess");
    if (NtSuspendProcessFunc == NULL) {
        printf("Failed to find NtSuspendProcess!\n");
        return 1;
    }

    NTSTATUS status = NtSuspendProcessFunc(ProcessHandle);
    if (status != 0) {
        printf("Failed to suspend process: 0x%x\n", (unsigned int)status);
        return 1;
    }

    return 0;
}

char *combine_path(const char* p1, const char* p2) {
    char *result = malloc(strlen(p1) + strlen(p2) + 2);
    strcpy(result, p1);
    strcat(result, "\\");
    strcat(result, p2);
    return result;
}

HANDLE start_process(char* processName, char* args, char* directory) {
    STARTUPINFO si;
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);

    PROCESS_INFORMATION pi;
    ZeroMemory(&pi, sizeof(pi));

    char* cmdLine = NULL;
    if (args) {
        size_t len = strlen(processName) + 1 + strlen(args) + 1;
        cmdLine = malloc(len);
        snprintf(cmdLine, len, "%s %s", processName, args);
    }

    if (!CreateProcess(
        processName,
        cmdLine,
        NULL,
        NULL,
        FALSE,
        0,
        NULL,
        directory,
        &si,
        &pi)) {
            printf("Failed to create process!\n");
            return NULL;
    }

    free(cmdLine);
    CloseHandle(pi.hThread);
    return pi.hProcess;
}

BOOL file_exists(LPCTSTR szPath)
{
  DWORD dwAttrib = GetFileAttributes(szPath);

  return (dwAttrib != INVALID_FILE_ATTRIBUTES && 
         !(dwAttrib & FILE_ATTRIBUTE_DIRECTORY));
}