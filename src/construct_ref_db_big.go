package main 

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"strconv"
	"encoding/binary"
)

func DumpArrayToFile(arr []uint16, filePath string) bool {
	length := len(arr);
	fp, err := os.Create(filePath)

	if err != nil {
		fmt.Printf("%v\n", err)
		return false;
	}
	wOut := bufio.NewWriter(fp);
	bitMask1 := uint16(0xFF00);
	bitMask2 := uint16(0x00FF);
	for i:=0;i<length;i++ {
		wOut.WriteByte(byte((arr[i]&bitMask1)>>8));
		wOut.WriteByte(byte(arr[i]&bitMask2));
		//fmt.Fprint(wOut, "%d", arr[i]);
	}
	wOut.Flush(); // Don't forget to flush!
	fp.Close();
	
	return true;
}

func ReadKmerFreqFromFile(kmerId uint64, filePath string) uint16 {
	fp, err := os.Open(filePath);

	if err != nil {
		fmt.Printf("%v\n", err);
		os.Exit(1);
	}
	
	//func (f *File) ReadAt(b []byte, off int64) (n int, err error)
	bytes := make([]byte, 2);
	fp.ReadAt(bytes, int64(kmerId<<1));
	//freq := uint16(bytes[1])<<8 | uint16(bytes[1]);
	freq := binary.BigEndian.Uint16(bytes);
	
	fp.Close();
	
	return freq;
}

func main() {
	
	kmer := 14;
	kmerBits := uint64(kmer)<<1;
	kmerCount := 1<<kmerBits;
	
	fmt.Printf("kmerCount = %d\n", kmerCount);
	
	gKmerCounts := make([]uint16, kmerCount);
	
	gKmerCounts[0]=65;
	gKmerCounts[1]=67;
	gKmerCounts[2]=68;
	gKmerCounts[3]=96;
	gKmerCounts[4]=1100;
	gKmerCounts[5]=255;
	gKmerCounts[6]=65535;
	//DumpArrayToFile(gKmerCounts, "abc.bin");
	
	var k_value = flag.String("k", "1234", "value of k in k-mer");
	flag.Parse();
	catchk, _ := strconv.ParseUint(*k_value, 10, 32);
	// assign k value to k
	kmerId := uint64(catchk);
	freq := ReadKmerFreqFromFile(kmerId, "abc.bin");
	fmt.Printf("k-mer = %d, Freq = %d\n", kmerId, freq);
}

