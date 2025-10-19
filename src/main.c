#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <Windows.h>
#include <stdbool.h>
#include "utils.h"

int main(int argc, char **argv) {
    // First I wanna check for cobalt
    char *cobaltPath = "./assets/Cobalt.dll";
    if (!file_exists(cobaltPath)) {
        printf("Error: Cobalt DLL not found!");
        return 1;
    }
    
    char *playPath = "";
    char **allTweaks;
    size_t allTweaksSize = 0;
    bool bIsServer = false;
    for (int i = 0; i < argc; i++) {
        if (strncmp("--play", argv[i], 7) == 0 || strncmp("-p", argv[i], 3) == 0) {
            playPath = argv[i+1];
        } else if (strncmp("--tweak", argv[i], 8) == 0 || strncmp("-t", argv[i], 3) == 0) {
            allTweaksSize++;
            allTweaks = realloc(allTweaks, allTweaksSize * sizeof(char*));
            allTweaks[allTweaksSize-1] = argv[i+1];
        } else if (strncmp("--server", argv[i], 9) == 0 || strncmp("-s", argv[i], 3) == 0) {
            bIsServer = true;
        }
    }
    
    char *serverPath = "./assets/Reboot.dll";
    if (bIsServer && !file_exists(serverPath)) {
        printf("Error: Reboot DLL not found!");
        return 1;
    }

    char* fortniteBinariesPath = combine_path(playPath, "FortniteGame\\Binaries\\Win64");
    char* fortniteLauncherPath = combine_path(fortniteBinariesPath, "FortniteLauncher.exe");
    char* fortniteClientPath = combine_path(fortniteBinariesPath, "FortniteClient-Win64-Shipping.exe");
    char* fortniteEACPath = combine_path(fortniteBinariesPath, "FortniteClient-Win64-Shipping_EAC.exe");

    char* launchArgs = "-epicapp=Fortnite -epicenv=Prod -epiclocale=en-us -epicportal -skippatchcheck -NOSSLPINNING -nobe -fromfl=eac -fltoken=3db3ba5dcbd2e16703f3978d -caldera=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiYmU5ZGE1YzJmYmVhNDQwN2IyZjQwZWJhYWQ4NTlhZDQiLCJnZW5lcmF0ZWQiOjE2Mzg3MTcyNzgsImNhbGRlcmFHdWlkIjoiMzgxMGI4NjMtMmE2NS00NDU3LTliNTgtNGRhYjNiNDgyYTg2IiwiYWNQcm92aWRlciI6IkVhc3lBbnRpQ2hlYXQiLCJub3RlcyI6IiIsImZhbGxiYWNrIjpmYWxzZX0.VAWQB67RTxhiWOxx7DBjnzDnXyyEnX7OljJm-j2d88G_WgwQ9wrE6lwMEHZHjBd1ISJdUO1UVUqkfLdU5nofBQ";

    // Launching
    HANDLE launcherHandle = start_process(fortniteLauncherPath, launchArgs, fortniteBinariesPath);
    NtSuspendProcess(launcherHandle);
    HANDLE acHandle = start_process(fortniteEACPath, launchArgs, fortniteBinariesPath);
    NtSuspendProcess(acHandle);
    HANDLE gameHandle = start_process(fortniteClientPath, launchArgs, fortniteBinariesPath);

    printf("Launched Fortnite...\n");

    //TODO: Inject dlls
    WaitForSingleObject(gameHandle, INFINITE);

    // cleanup
    TerminateProcess(launcherHandle, 0);
    TerminateProcess(acHandle, 0);

    CloseHandle(gameHandle);
    CloseHandle(launcherHandle);
    CloseHandle(acHandle);

    free(fortniteLauncherPath);
    free(fortniteClientPath);
    free(fortniteEACPath);
    free(fortniteBinariesPath);
    free(playPath);
    free(allTweaks);

    return 0;
}