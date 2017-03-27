/**
 * Created by Philippe on 2017-03-26.
 */
import com.sun.jna.Memory;
import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.Structure;
import com.sun.jna.platform.win32.*;
import com.sun.jna.ptr.IntByReference;
import com.sun.jna.win32.W32APIOptions;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

import static sun.security.krb5.internal.Krb5.DEBUG;

public class MemManip {
    static Kernel32 kernel32 = (Kernel32) Native.loadLibrary(Kernel32.class, W32APIOptions.UNICODE_OPTIONS);
    public static int PROCESS_VM_READ= 0x0010;
    public static int PROCESS_VM_WRITE = 0x0020;
    public static int PROCESS_VM_OPERATION = 0x0008;
    public static int PROCESS_VM_QUERY_INFO = 0x0400;
    public List<WinNT.MEMORY_BASIC_INFORMATION> readablePages;
    public LinkedHashMap<String,Integer> valueContainer;
    WinNT.HANDLE processHandle;
    Pointer process;
    int size_t = 0;
    int PID=0;
    Memory memBuffer;

    public int [] FindAddressesWith(String query) throws Exception {
        if(PID == 0){
            throw new Exception("No PID has been defined yet.");
        }
        if(this.process == null){
                this.process = OpenProcess(PID);
        }
        return null;
    }

    public int FindProcessId(String processName)
    {
        // This Reference will contain the processInfo that will be parsed to recover the ProcessId
        Tlhelp32.PROCESSENTRY32.ByReference processInfo = new Tlhelp32.PROCESSENTRY32.ByReference();

        // This Handle allows us to parse the process map
        WinNT.HANDLE processesSnapshot = kernel32.CreateToolhelp32Snapshot(Tlhelp32.TH32CS_SNAPPROCESS, new WinDef.DWORD(0L));
        if(processesSnapshot == kernel32.INVALID_HANDLE_VALUE)
        {
            if(DEBUG) System.err.println("INVALID_HANDLE_VALUE");
            return 0;
        }

        try{// This will parse all the processes to find the process id corresponding to the process name
            kernel32.Process32First(processesSnapshot, processInfo);
            if(processName.equals(Native.toString(processInfo.szExeFile)))
            {
                if(DEBUG) System.out.println("Process " + processName + " found : " + processInfo.th32ProcessID.longValue());
                return processInfo.th32ProcessID.intValue();
            }

            while(kernel32.Process32Next(processesSnapshot, processInfo))
            {
                if(processName.equals(Native.toString(processInfo.szExeFile)))
                {
                    if(DEBUG) System.out.println("Process " + processName + " found : " + processInfo.th32ProcessID.longValue());
                    return processInfo.th32ProcessID.intValue();
                }
            }

            if(DEBUG) System.out.println("Did not found the requested Process: " + processName);
            return 0;
        }
        finally
        {
            kernel32.CloseHandle(processesSnapshot);
        }
    }

    private static Pointer OpenProcess(int pid) {
        WinNT.HANDLE process = kernel32.OpenProcess(PROCESS_VM_READ|0x0400|PROCESS_VM_WRITE|PROCESS_VM_OPERATION, true, pid);
        return process.getPointer();
    }

    public boolean OpenProcess() {
        this.processHandle = kernel32.OpenProcess(PROCESS_VM_READ|0x0400|PROCESS_VM_WRITE|PROCESS_VM_OPERATION, true, this.PID);
        this.process = this.processHandle.getPointer();
        return this.process != null;
    }

    public boolean setSize(int size_t){
        this.size_t = size_t;
        return size_t!=0;
    }

    public boolean loadMemBuffer(int bufferSize){
        memBuffer = new Memory(bufferSize);
        Pointer baseAdd = new Pointer(0x0);
        IntByReference readBytes = new IntByReference(0);
        kernel32.ReadProcessMemory(this.processHandle,baseAdd,memBuffer,8,readBytes);
        System.out.println("Read bytes: "+readBytes.toString());
        System.out.println(memBuffer.getCharArray(0,1));
        return false;
    }

    public int searchFor(int value,int size){
        this.valueContainer = new LinkedHashMap<>();
        memBuffer = new Memory(size);
        Pointer add = new Pointer(0);
        IntByReference readBytes = new IntByReference(0);
        int count=0;
        for(WinNT.MEMORY_BASIC_INFORMATION page : this.readablePages){
            int offset=0;
            Pointer current = page.baseAddress;
            Pointer last = new Pointer(pointerToAddress(page.baseAddress)+page.regionSize.longValue());
            long memSize = pointerToAddress(last)-pointerToAddress(current);
            memBuffer = new Memory((int)memSize);
            kernel32.ReadProcessMemory(this.processHandle,current,memBuffer,(int)memSize,readBytes);
            while(offset<(int)memSize-2){
                memBuffer.getInt(offset);
                if(value == memBuffer.getInt(offset)){
                    valueContainer.put(String.format("0x%08X", offset+pointerToAddress(current)),memBuffer.getInt(offset));
                    System.out.println(String.format("0x%08X", offset+pointerToAddress(current))+":"+memBuffer.getInt(offset));
                }
                offset+=size;
            }

        }
        return this.valueContainer.size();
    }

    public int narrow(int value, int size){
        System.out.println("Started narrowing");
        LinkedHashMap<String, Integer> temp=new LinkedHashMap<>();
        for(Map.Entry<String,Integer>entry : this.valueContainer.entrySet()){
            String entryKey = entry.getKey();
            memBuffer=new Memory(size);
            kernel32.ReadProcessMemory(this.processHandle,new Pointer(addressToLong(entryKey)),memBuffer,size,new IntByReference(0));
            if(memBuffer.getInt(0)==value){
                temp.put(entryKey,memBuffer.getInt(0));
                System.out.println(entryKey+":"+memBuffer.getInt(0));
            }
        }
        this.valueContainer=temp;
        return this.valueContainer.size();
    }

    public static List<WinNT.MEMORY_BASIC_INFORMATION> getPageRanges(WinNT.HANDLE hOtherProcess) {
        List<WinNT.MEMORY_BASIC_INFORMATION> ret = new ArrayList<>();
        WinNT.MEMORY_BASIC_INFORMATION mbi;
        WinBase.SYSTEM_INFO si = new WinBase.SYSTEM_INFO();
        Kernel32.INSTANCE.GetSystemInfo(si);
        Pointer lpMem = si.lpMinimumApplicationAddress;
        System.out.println("lpMem: "+pointerToAddress(lpMem));
        System.out.println("lpMaxMem: "+pointerToAddress(si.lpMaximumApplicationAddress));
        while (pointerToAddress(lpMem) < pointerToAddress(si.lpMaximumApplicationAddress)) {
            mbi = new WinNT.MEMORY_BASIC_INFORMATION();
            BaseTSD.SIZE_T t = Kernel32.INSTANCE.VirtualQueryEx(hOtherProcess, lpMem, mbi, new BaseTSD.SIZE_T(mbi.size()));
            if (t.longValue() == 0) {
                System.out.println("Cannot get page ranges. Last error:" + Kernel32.INSTANCE.GetLastError());
                break;
            }
            ret.add(mbi);
            lpMem = new Pointer(pointerToAddress(mbi.baseAddress) + mbi.regionSize.longValue());
        }
        return ret;
    }

    public boolean loadPageRanges() {
        List<WinNT.MEMORY_BASIC_INFORMATION> ret = new ArrayList<>();
        WinNT.MEMORY_BASIC_INFORMATION mbi;
        WinBase.SYSTEM_INFO si = new WinBase.SYSTEM_INFO();
        Kernel32.INSTANCE.GetSystemInfo(si);
        Pointer lpMem = si.lpMinimumApplicationAddress;
        while (pointerToAddress(lpMem) < pointerToAddress(si.lpMaximumApplicationAddress)) {
            mbi = new WinNT.MEMORY_BASIC_INFORMATION();
            BaseTSD.SIZE_T t = Kernel32.INSTANCE.VirtualQueryEx(this.processHandle, lpMem, mbi, new BaseTSD.SIZE_T(mbi.size()));
            if (t.longValue() == 0) {
                System.out.println("Cannot get page ranges. Last error:" + Kernel32.INSTANCE.GetLastError());
                return false;
            }
            ret.add(mbi);
            lpMem = new Pointer(pointerToAddress(mbi.baseAddress) + mbi.regionSize.longValue());
        }
        this.readablePages = ret;
        if(!fractureMemChunks())
            return false;
        return true;
    }

    public static long pointerToAddress(Pointer ptr){
        String val = ptr.toString();
        String cut = val.substring(9);
        return new BigInteger(cut,16).longValue();
    }

    public static long addressToLong(String addr){
        String cut = addr.substring(2);
        return new BigInteger(cut,16).longValue();
    }

    public boolean fractureMemChunks(){
        boolean fractured = false;
        WinNT.MEMORY_BASIC_INFORMATION newPage = null, oldPage = null, lastPage=null;
        for(WinNT.MEMORY_BASIC_INFORMATION page:this.readablePages){
            Pointer current = page.baseAddress;
            Pointer last = new Pointer(pointerToAddress(page.baseAddress)+page.regionSize.longValue());
            long memSize = pointerToAddress(last)-pointerToAddress(current);
            if (memSize > 1955555){
                oldPage = page;
                newPage = page;
                lastPage = page;
                oldPage.regionSize = new BaseTSD.SIZE_T(1955554);
                newPage.baseAddress = new Pointer(pointerToAddress(page.baseAddress)+1955554);

                fractured = true;
                break;
            }

        }
        /**
         *  TODO:This is VERY RISKY CODE... MUST REFACTOR TO SAFE HAVEN (aka no possibility for infinite loop)
         */
        if(fractured) {
            this.readablePages.remove(this.readablePages.indexOf(lastPage));
            this.readablePages.add(oldPage);
            this.readablePages.add(newPage);
            return fractureMemChunks();
        }
        return true;
    }
}
