/**
 * Created by Philippe on 2017-03-27.
 */
    import com.sun.jna.Native;
    import com.sun.jna.win32.W32APIOptions;
    import com.sun.jna.platform.win32.Kernel32;
    import com.sun.jna.platform.win32.WinDef.DWORD;
    import com.sun.jna.platform.win32.Tlhelp32;
    import com.sun.jna.platform.win32.WinNT;

//Imports for MyKernel32
    import com.sun.jna.win32.StdCallLibrary;
    import com.sun.jna.ptr.IntByReference;
    import com.sun.jna.Pointer;
    import com.sun.jna.Memory;

    public class MyProcessReaderExample
    {
        private static final boolean DEBUG = false;
        static Kernel32 kernel32 = (Kernel32)Native.loadLibrary(Kernel32.class, W32APIOptions.UNICODE_OPTIONS);
        static MyKernel32 myKernel32 = (MyKernel32)Native.loadLibrary("kernel32", MyKernel32.class);
        static String usage = "Usage: java MyProcessReaderExample [processName] [readSize] [readAddress] [readOffset]" + System.getProperty("line.separator") + "processName example: firefox.exe" + System.getProperty("line.separator") + "readSize (in bytes) example : 4" + System.getProperty("line.separator") + "readAddress (hexadecimal) example : 00010000" + System.getProperty("line.separator") + "readOffset (decimal) example : 0";

        interface MyKernel32 extends StdCallLibrary
        {
            boolean WriteProcessMemory(Pointer p, int address, Pointer buffer, int size, IntByReference written);
            boolean ReadProcessMemory(Pointer hProcess, int inBaseAddress, Pointer outputBuffer, int nSize, IntByReference outNumberOfBytesRead);
            Pointer OpenProcess(int desired, boolean inherit, int pid);
            int GetLastError();
        }

        public static void main(String[] args)
        {
            if(args.length < 4)
            {
                System.err.println(usage);
                System.exit(1);
            }
            String processName = "";
            int readSize = 0;
            int readAddress = 0;
            int readOffset = 0;
            try{
                processName = args[0];//Defines the name of the process to search for
                readSize = Integer.parseInt(args[1], 10);//Defines the size of the value to read
                readAddress = Integer.parseInt(args[2], 16);//Defines the address of the value to read (hexadecimal)
                readOffset = Integer.parseInt(args[3], 10);//Defines the offset of the value to read
            }
            catch(NumberFormatException e)
            {
                System.err.println(usage);
                System.exit(1);
            }

            long processId = FindProcessId(processName);
            if(processId == 0L)
            {
                System.err.println("The searched process was not found : " + processName);
                System.exit(1);
            }
            System.out.println(processName + " id : " + processId);

            Pointer openedProcess = myKernel32.OpenProcess(0x0010 + 0x0020 + 0x0008, true, (int)processId);

            Memory output = new Memory(readSize);
            myKernel32.ReadProcessMemory(openedProcess, readAddress, output, readSize, new IntByReference(0));
            int readValue = output.getInt(readOffset);//This should change if the readSize is different than 4
            System.out.println(readValue + " has been read at 0x" + args[2] + "+" + readOffset + " on " + readSize + "bytes.");

            IntByReference written = new IntByReference(0);
            Memory toWrite = new Memory(readSize);
            toWrite.setInt(readOffset, readValue+1);//This should change if the readSize is different than 4
            myKernel32.WriteProcessMemory(openedProcess, readAddress, toWrite, 1, written);
            System.out.println("Wrote " + written.getValue() + " times.");

            output = new Memory(readSize);
            myKernel32.ReadProcessMemory(openedProcess, readAddress, output, readSize, new IntByReference(0));
            readValue = output.getInt(readOffset);//This should change if the readSize is different than 4
            System.out.println(readValue + " has been read at 0x" + args[2] + "+" + readOffset + " on " + readSize + "bytes.");

            written = new IntByReference(0);
            toWrite = new Memory(readSize);
            toWrite.setInt(readOffset, readValue-1);//This should change if the readSize is different than 4
            myKernel32.WriteProcessMemory(openedProcess, readAddress, toWrite, 1, written);
            System.out.println("Wrote " + written.getValue() + " times.");

            output = new Memory(readSize);
            myKernel32.ReadProcessMemory(openedProcess, readAddress, output, readSize, new IntByReference(0));
            readValue = output.getInt(readOffset);//This should change if the readSize is different than 4
            System.out.println(readValue + " has been read at 0x" + args[2] + "+" + readOffset + " on " + readSize + "bytes.");
        }

        static long FindProcessId(String processName)
        {
            Tlhelp32.PROCESSENTRY32.ByReference processInfo = new Tlhelp32.PROCESSENTRY32.ByReference();
            WinNT.HANDLE processesSnapshot = kernel32.CreateToolhelp32Snapshot(Tlhelp32.TH32CS_SNAPPROCESS, new DWORD(0L));

            try{
                kernel32.Process32First(processesSnapshot, processInfo);
                if(processName.equals(Native.toString(processInfo.szExeFile)))
                {
                    return processInfo.th32ProcessID.longValue();
                }

                while(kernel32.Process32Next(processesSnapshot, processInfo))
                {
                    if(processName.equals(Native.toString(processInfo.szExeFile)))
                    {
                        return processInfo.th32ProcessID.longValue();
                    }
                }
                return 0L;
            }
            finally
            {
                kernel32.CloseHandle(processesSnapshot);
            }
        }
    }
