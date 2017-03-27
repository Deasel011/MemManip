//import com.sun.jna.Memory;
//import com.sun.jna.Native;
//import com.sun.jna.Pointer;
//import com.sun.jna.platform.win32.*;
//import com.sun.jna.platform.win32.Kernel32;
//import com.sun.jna.ptr.IntByReference;
//
///**
// * Created by Philippe on 2017-03-27.
// */
//public class Tutorial {
//    static Kernel32 kernel32 = (Kernel32) Native.loadLibrary("kernel32", Kernel32.class);
//    static User32 user32 = (User32) Native.loadLibrary("user32", User32.class);
//
//    public static void main(String[] args){
//        long pid = getProcessId("Game.exe");
//        System.out.println("long pid:"+pid);
//        pid = 4160;
//        WinNT.HANDLE readprocess = openHandle(0x0010,(int)pid);
//        int size = 20;
//        Memory read = readMemory(readprocess,0x02992D58L,size);
//        System.out.println(read.getCharArray(0,0));
//    }
//
//    public static long getProcessId(String processName)
//    {
//        Tlhelp32.PROCESSENTRY32.ByReference processInfo = new Tlhelp32.PROCESSENTRY32.ByReference();
//        WinNT.HANDLE processesSnapshot = kernel32.CreateToolhelp32Snapshot(Tlhelp32.TH32CS_SNAPPROCESS, new WinDef.DWORD(0L));
//
//        try{
//            kernel32.Process32First(processesSnapshot, processInfo);
//            if(processName.equals(Native.toString(processInfo.szExeFile)))
//            {
//                return processInfo.th32ProcessID.longValue();
//            }
//
//            while(kernel32.Process32Next(processesSnapshot, processInfo))
//            {
//                if(processName.equals(Native.toString(processInfo.szExeFile)))
//                {
//                    return processInfo.th32ProcessID.longValue();
//                }
//            }
//            return 0L;
//        }
//        finally
//        {
//            kernel32.CloseHandle(processesSnapshot);
//        }
////        IntByReference pid = new IntByReference(0);
////        user32.GetWindowThreadProcessId(user32.FindWindow(null,window), pid);
////
////        return pid.getValue();
//    }
//
//    public static WinNT.HANDLE openHandle(int permissions, int pid)
//    {
//        System.out.println("Int pid:"+pid);
//        WinNT.HANDLE process = kernel32.OpenProcess(permissions,true, pid);
//        return process;
//    }
//
//    public static Memory readMemory(WinNT.HANDLE process, long address, int bytesToRead)
//    {
//        IntByReference read = new IntByReference(0);
//        Pointer baseAdd = new Pointer(address);
//        Memory output = new Memory(bytesToRead);
//
//        kernel32.ReadProcessMemory(process, baseAdd, output, bytesToRead, read);
//        return output;
//    }
//}
