use std::{ffi::CString, os::raw::c_void, ptr::null_mut};

use windows::{Win32::{Foundation::{CloseHandle, HANDLE}, System::{Diagnostics::Debug::WriteProcessMemory, LibraryLoader::{GetModuleHandleA, GetProcAddress}, Memory::{MEM_COMMIT, MEM_RELEASE, MEM_RESERVE, PAGE_READWRITE, VirtualAllocEx, VirtualFreeEx}, Threading::{LPTHREAD_START_ROUTINE, OpenProcess, PROCESS_ALL_ACCESS, PROCESS_QUERY_INFORMATION, PROCESS_SUSPEND_RESUME}}}, core::PCSTR};

#[link(name = "kernel32")]
unsafe extern "system" {
    fn CreateRemoteThread(hProcess: HANDLE, lpThreadAttributes: *mut c_void, dwStackSize: usize, lpStartAddress: LPTHREAD_START_ROUTINE, lpParameter: *mut c_void, dwCreationFlags: u32, lpThreadId: *mut u32) -> HANDLE;
}

#[link(name = "ntdll")]
unsafe extern "system" {
    fn NtSuspendProcess(proc: HANDLE) -> i32;
}

pub fn nt_suspend_process(pid: u32) -> bool {
    unsafe {
        let handle = OpenProcess(
            PROCESS_SUSPEND_RESUME | PROCESS_QUERY_INFORMATION,
            false,
            pid
        ).unwrap();

        if handle.is_invalid() {
            return false;
        }

        let status = NtSuspendProcess(handle);
        let _ = CloseHandle(handle);

        status == 0
    }
}
    
// This is directly skidded from the original project
pub fn inject_dll(pid: u32, dll_path: &str) -> bool {
    let path_size = dll_path.len();

    unsafe {
        let handle = OpenProcess(
            PROCESS_ALL_ACCESS, 
            false, 
            pid
        ).unwrap();

        if handle.is_invalid() {
            return false;
        }

        let mem_loc = VirtualAllocEx(
            handle, 
            None, 
            path_size, 
            MEM_COMMIT | MEM_RESERVE, 
            PAGE_READWRITE
        );

        let write_process_memory = WriteProcessMemory(
            handle, 
            mem_loc, 
            CString::new(dll_path).unwrap().as_ptr() as *mut c_void, 
            path_size, 
            None
        );
        match write_process_memory {
            Ok(_) => {},
            Err(_) => {
                return false;
            }
        }
        
        let kernel32_dll_name = CString::new("kernel32.dll").unwrap();
        let kernel32_handle = GetModuleHandleA(PCSTR::from_raw(kernel32_dll_name.as_ptr() as *const u8));

        let loadlibrary_name = CString::new("LoadLibraryA").unwrap();

        let load_library_addr = GetProcAddress(kernel32_handle.unwrap(), PCSTR::from_raw(loadlibrary_name.as_ptr() as *const u8));
        if load_library_addr.is_none() {
            return false;
        }
        let load_library = load_library_addr.unwrap();

        let load_library_fn: LPTHREAD_START_ROUTINE = std::mem::transmute(load_library);

        let thread_handle = CreateRemoteThread(
            handle, 
            null_mut(), 
            0, 
            load_library_fn, 
            mem_loc, 
            0, 
            null_mut()
        );
        if thread_handle.is_invalid() {
            let _ = VirtualFreeEx(
                handle, 
                mem_loc, 
                0, 
                MEM_RELEASE
            );
            return false;
        }
    }

    true
}