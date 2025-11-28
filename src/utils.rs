use std::ffi::CString;
use std::ptr::null_mut;
use std::os::raw::c_void;
use std::mem::transmute;

type HANDLE = *mut c_void;
type NTSTATUS = i32;

#[link(name = "kernel32")]
unsafe extern "system" {
    fn LoadLibraryA(lpLibFileName: *const i8) -> HANDLE;
    fn GetProcAddress(hModule: HANDLE, lpProcName: *const i8) -> *mut c_void;
    fn OpenProcess(dwDesiredAccess: u32, bInheritHandle: i32, dwProcessId: u32) -> HANDLE;
    fn CloseHandle(hObject: HANDLE) -> i32;
}

const PROCESS_SUSPEND_RESUME: u32 = 0x0800;
const PROCESS_QUERY_INFORMATION: u32 = 0x0400;

fn nt_suspend_process(processHandle: HANDLE) -> i32 {
    unsafe {
        let ntdll_name = CString::new("ntdll.dll").unwrap();
        let ntdll = LoadLibraryA(ntdll_name.as_ptr());
        if ntdll.is_null() {
            panic!("Failed to load ntdll.dll");
        }

        let func_name = CString::new("NtSuspendProcess").unwrap();
        let addr = GetProcAddress(ntdll, func_name.as_ptr());
        if addr.is_null() {
            panic!("NtSuspendProcess not found!");
        }

        let nt_suspend_process_func: extern "system" fn(HANDLE) -> NTSTATUS = transmute(addr);

        nt_suspend_process_func(processHandle)
    }
}
