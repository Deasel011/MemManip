import com.sun.jna.Memory;
import com.sun.jna.Pointer;
import com.sun.jna.platform.win32.WinNT;
import com.sun.jna.ptr.IntByReference;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by Philippe on 2017-03-26.
 */
public class Test {
    public static int SIZE = 4;


    public static void main(String[] args) throws Exception {
        int readSize = 0;
        int readValue = 0;
        int readAddress = 0;
        int readOffset = 0;
        MemManip manip = new MemManip();

        manip.PID = manip.FindProcessId("Game.exe");
        System.out.println("is opened? "+manip.OpenProcess());
        manip.loadPageRanges();
        manip.searchFor(384,4);
        //manip.getRegions(4);
        //manip.searchFor(384,4);
//        List<WinNT.MEMORY_BASIC_INFORMATION> pages = MemManip.getPageRanges(manip.processHandle);
//        System.out.println(pages.size());
//        for(WinNT.MEMORY_BASIC_INFORMATION page : pages){
//            System.out.println(page);
//        }


    }
}
