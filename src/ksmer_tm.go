package main 

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/gob"
    "bufio"
    "bytes"
	"path/filepath"
	"flag"
	"strconv"
	"strings"
)

var gkcMap map[uint32]map[uint16]int;


// Catch and process errors using panic
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Translate an entry in the byte array read from the xxxxxx.fasta file to a 2-bit integer code
// N (a mis-read base) must be stripped from input and processed before doing kmer frequency
func getBaseId(b byte) uint32 {
	switch(b) {
		case 'A':
			return 0;
		case 'C':
			return 1;
		case 'G':
			return 2;
		case 'T':
			return 3;
		case 'N':
			return 4;
	}
	return 4;
}

// Increment the frequency of a particular kmer in a particular genome
// Create that kmer entry if not already present
func store(genomeId uint16, kmerId uint32) {
	gMap, ok := gkcMap[kmerId];
	if !ok {
		gMap = make(map[uint16]int);
		gkcMap[kmerId] = gMap;
	}
	gMap[genomeId]++;
}

// Given the genome information in contigStr, a k-value, and a genome ID (1, 2, 3, etc.)
// determine the frequency of each kmer in contigStr
func storeContigMappings(genomeId uint16, contigStr []byte, k uint32) {
	var i uint32;
	var kmerVal uint32 = 0;
	// Spool off k-1 bits of contigStr to find first kmer
func storeContigMappings(genomeId uint16, contigStr []byte, k uint32) {
	var i uint32;
	var kmerVal uint32 = 0;
	for i=0; i<(k-1); i++ {
		kmerVal = (kmerVal<<2) + getBaseId(contigStr[i]);
		//fmt.Printf("%d\n", kmerVal);
	}
	// mask off high bits in kmer
	bitMask := uint32((1<<(k<<1))-1);
	length := uint32(len(contigStr));
	for i=k-1; i<length; i++ {
		kmerVal = (kmerVal<<2 & bitMask) + getBaseId(contigStr[i]);
		store(genomeId, kmerVal);
		//fmt.Printf("%d %c %d %d\n", i, contigStr[i], getBaseId(contigStr[i]), kmerVal);
	}
}

// Given the genome info in contigStr, a genome ID (1, 2, 3, etc.) and a spaced kmer (k1 + s + k2, where k1 the length
// 	of the first half of the spaced kmer, s is the space between the coding portions of the kmer, and k2 is the second
//	half of the spaced kmer) determine the frequency of the spaced kmer in contigStr
func storeContigMappingsSpaced(genomeId uint16, contigStr []byte, k1 uint32, s uint32, k2 uint32) {
	var i,j uint32;
	var kmerValFirst uint32 = 0;
	var kmerValSecond uint32 = 0;
	
	// feed in k1-1 bases into the first half of the kmer
	for i=0; i<(k1-1); i++ {
		kmerValFirst = (kmerValFirst<<2) + getBaseId(contigStr[i]);
	}
	// feed in k2-1 bases into the second half of the kmer, ignoring s bases in between
	for i=0; i<(k1-1); i++ {
		kmerValFirst = (kmerValFirst<<2) + getBaseId(contigStr[i]);
	}
	for i=k1+s; i<(k1+s+k2-1); i++ {
		kmerValSecond = (kmerValSecond<<2) + getBaseId(contigStr[i]);
	}
	
	// mask off high bits in the first and second halves of the spaced kmer
	bitMask1 := uint32((1<<(k1<<1))-1);
	bitMask2 := uint32((1<<(k2<<1))-1);
	length := uint32(len(contigStr));
	var kmerVal uint32;
	i=k1-1;
	// read spaced kmer and record counts in map
	for j=k1+s+k2-1; j<length; j++ {
		kmerValFirst = (kmerValFirst<<2 & bitMask1) + getBaseId(contigStr[i]);
		kmerValSecond = (kmerValSecond<<2 & bitMask2) + getBaseId(contigStr[j]);
		kmerVal = kmerValFirst<<(k2<<1) + kmerValSecond;
		store(genomeId, kmerVal);
		//fmt.Printf("%d %c %d %d\n", j, contigStr[j], getBaseId(contigStr[j]), kmerVal);
		fmt.Printf("%d %c %d %d\n", j, contigStr[j], getBaseId(contigStr[j]), kmerVal);
		i++;
	}
}

// Store map m with frequency counts for all present kmers at filePath using gob Encoder
func saveIndexMapToFile(filePath string, m map[uint32]map[uint16]int) {
	fp, err := os.Create(filePath);
	if err != nil {
		panic("cant open file");
	}
	defer fp.Close();
	
	enc := gob.NewEncoder(fp);
	if err := enc.Encode(m); err != nil {
		panic("cant encode");
	}
}

// Load map m from filePath string using gob Decoder
func loadIndexMapToFile(filePath string) (m map[uint32]map[uint16]int) {
	fp, err := os.Open(filePath);
	if err != nil {
		panic("cant open file");
	}
	defer fp.Close();
	
	enc := gob.NewDecoder(fp);
	if err := enc.Decode(&m); err != nil {
		panic("cant decode");
	}
	return m;
}

// read contigStr from a .fasta file
func ReadFASTA(sequence_file string) []byte {
    f, err := os.Open(sequence_file)

    if err != nil {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)
    byte_array := bytes.Buffer{}
    var isPrefix bool
    _, isPrefix, err = br.ReadLine()
    
    if err != nil || isPrefix {
        fmt.Printf("%v\n", err)
        os.Exit(1)
    }

    var line []byte
    
    for {
        line, isPrefix, err = br.ReadLine()
        if err != nil || isPrefix {
            break
        } else {
            byte_array.Write(line)
        }
    }

    return []byte(byte_array.String())
}

// TODO: in_dir and out_dir cmd line args
//			.fastq processing
func main() {
	
	gkcMap = make(map[uint32]map[uint16]int);
	
	var k_value = flag.String("k", "", "value of k in k-mer")
	var spaced_kmer = flag.String("s", "", "use spaced kmer")
	//var in_dir = flag.String("i", "", "input directory")
	//var out_dir = flag.String("o", "", "output directory")
	flag.Parse();
	
	var input_dir = "../input/";
	var output_dir = "../output/";
	var k = uint32(4);
	var k1 = uint32(4);
	var sp = uint32(4);
	var k2 = uint32(4);
	var spaced = false;
	
	if(*k_value != "") {		
		// catch multi-value return
		catchk, err := strconv.ParseUint(*k_value, 10, 32);
		check(err);
						
		// assign k value to k
		k = uint32(catchk);
	}
	if(*spaced_kmer != "") {
		spaced = true;
		temp := strings.Split(*spaced_kmer, ",");
				
		catchk1, err := strconv.ParseUint(temp[0], 10, 32);
		check(err);
		catchsp, err := strconv.ParseUint(temp[1], 10, 32);
		check(err);
		catchk2, err := strconv.ParseUint(temp[2], 10, 32);
		check(err);
		
		k1 = uint32(catchk1);
		sp = uint32(catchsp);
		k2 = uint32(catchk2);
	}
	
	// Read list of input files into FileInfo slice named files
	files, err := ioutil.ReadDir(input_dir);
	check(err);
	
	if(spaced == false) {
		// Iterate through files, locating *.fasta
		for i, f := range files {
			if(filepath.Ext(f.Name()) == ".fasta") {
				contigStr := ReadFASTA(input_dir+f.Name());
			
				storeContigMappings(uint16(i), contigStr, k);
			}
		}
	} else {
		for i, f := range files {
			if(filepath.Ext(f.Name()) == ".fasta") {
				contigStr := ReadFASTA(input_dir+f.Name());

				storeContigMappingsSpaced(uint16(i), contigStr, k1, sp, k2);
			}
		}		
	}
	
	saveIndexMapToFile(output_dir+"index.map", gkcMap);
}
