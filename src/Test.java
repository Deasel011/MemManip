import com.sun.jna.Memory;
import com.sun.jna.Native;
import com.sun.jna.Pointer;
import com.sun.jna.platform.win32.Kernel32;
import com.sun.jna.platform.win32.WinNT;
import com.sun.jna.ptr.IntByReference;
import com.sun.jna.win32.W32APIOptions;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;


/**
 * Created by Philippe on 2017-03-26.
 */
public class Test {
    public static int SIZE = 4;
    static Kernel32 kernel32 = (Kernel32) Native.loadLibrary(Kernel32.class, W32APIOptions.UNICODE_OPTIONS);



    public static void main(String[] args) throws Exception {
        MemManip manip = new MemManip();

        manip.PID = manip.FindProcessId("Darkest.exe");
        System.out.println(manip.PID);
        System.out.println("is opened? "+manip.OpenProcess());
        manip.loadPageRanges();
        manip.searchFor(5345,SIZE);
//        System.out.println("Sell item now");
//        TimeUnit.SECONDS.sleep(10);
//        manip.narrow(222,SIZE);
//        System.out.println("Rebuy item now");
//        TimeUnit.SECONDS.sleep(10);
//        manip.narrow(225,SIZE);
//        System.out.println(manip.intAtSingleEntry(4));


//        System.out.println(manip.intAt("0x1F706D60",SIZE));




    }
}
