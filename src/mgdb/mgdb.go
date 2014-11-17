package mgdb

import (
	"fmt"
//	"math"
	"io/ioutil"
	"os"
	"encoding/gob"
	"bufio"
	"bytes"
//	"flag"
	"path/filepath"
	"strconv"
//	"strings"
)

var IsMatrix bool = true;
var GkcMap map[uint32]map[uint16]int;
var GkcMatrix [][]int;

var K = uint32(8);
var K1 = uint32(4);
//var sp = uint32(10);
var Sp = uint32(0);
var K2 = uint32(4);
var IsSpaced = false;


func Check(e error) {
    if e != nil {
        panic(e)
    }
}

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

func store(genomeId uint16, kmerId uint32) {
	if IsMatrix {
		//fmt.Printf("Store: %d,%d\n", kmerId, genomeId);
		GkcMatrix[kmerId][genomeId]++;
	} else {
		gMap, ok := GkcMap[kmerId];
		if !ok {
			gMap = make(map[uint16]int);
			GkcMap[kmerId] = gMap;
		}
		gMap[genomeId]++;
	}
}

func storeContigMappings(genomeId uint16, contigStr []byte, k uint32) {
	fmt.Printf("genome=%d, Storing contig (k=%d), length=%d\n", genomeId, k, len(contigStr));
	
	if len(contigStr)< int(k) {
		return;
	}
	
	var i uint32;
	var kmerVal uint32 = 0;
	for i=0; i<(k-1); i++ {
		kmerVal = (kmerVal<<2) + getBaseId(contigStr[i]);
		//fmt.Printf("%d\n", kmerVal);
	}
	bitMask := uint32((1<<(k<<1))-1);
	length := uint32(len(contigStr));
	for i=k-1; i<length; i++ {
		kmerVal = (kmerVal<<2 & bitMask) + getBaseId(contigStr[i]);
		store(genomeId, kmerVal);
		//fmt.Printf("%d %c %d %d\n", i, contigStr[i], getBaseId(contigStr[i]), kmerVal);
	}
}

func storeContigMappingsSpaced(genomeId uint16, contigStr []byte, k1 uint32, s uint32, k2 uint32) {
	
	fmt.Printf("genome=%d, Storing contig (k=%d-%d-%d), length=%d\n", genomeId, k1, s, k2, len(contigStr));
	
	if len(contigStr)< int(k1+s+k2) {
		return;
	}
	
	var i,j uint32;
	var kmerValFirst uint32 = 0;
	var kmerValSecond uint32 = 0;
	
	for i=0; i<(k1-1); i++ {
		kmerValFirst = (kmerValFirst<<2) + getBaseId(contigStr[i]);
	}
	for i=k1+s; i<(k1+s+k2-1); i++ {
		kmerValSecond = (kmerValSecond<<2) + getBaseId(contigStr[i]);
	}
	
	bitMask1 := uint32((1<<(k1<<1))-1);
	bitMask2 := uint32((1<<(k2<<1))-1);
	length := uint32(len(contigStr));
	var kmerVal uint32;
	i=k1-1;
	for j=k1+s+k2-1; j<length; j++ {
		kmerValFirst = (kmerValFirst<<2 & bitMask1) + getBaseId(contigStr[i]);
		kmerValSecond = (kmerValSecond<<2 & bitMask2) + getBaseId(contigStr[j]);
		kmerVal = kmerValFirst<<(k2<<1) + kmerValSecond;
		store(genomeId, kmerVal);
		//fmt.Printf("%d %c %d %d\n", j, contigStr[j], getBaseId(contigStr[j]), kmerVal);
		i++;
	}
}

//func saveIndexMapToFile(filePath string, m map[uint32]map[uint16]int) {
func SaveIndexMapToFile(filePath string, m interface{}) {
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

func loadIndexMapFromFile(filePath string) (m map[uint32]map[uint16]int) {
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
func LoadIndexArrayFromFile(filePath string) (m [][]int) {
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

func storeSequenceWithN(seq []byte, genomeId uint16) {
	fromIndex := 0;
	length := len(seq);
	for i:=0; i<length; i++ {
		if seq[i] == 'N' || i==(length-1) {
			if fromIndex == i {
				fromIndex = i+1;
				continue;
			}
			contigStr := seq[fromIndex:i];
			if IsSpaced {
				storeContigMappingsSpaced(genomeId, contigStr, K1, Sp, K2);
			} else {
				storeContigMappings(genomeId, contigStr, K);
			}
			fromIndex = i+1;
		}
	}
}

func ParseReadFnaStore(filePath string, genomeId uint16) {
	
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer f.Close()
	br := bufio.NewReader(f)
	var isPrefix bool

	var line []byte
	
	readStr := bytes.Buffer{}
	
	//counter := 0;
	for {
		line, isPrefix, err = br.ReadLine()
		
		if err != nil || isPrefix {
			//fmt.Printf("%d - %d\n",counter,len(readStr.String()));
			//counter++;
			storeSequenceWithN([]byte(readStr.String()), genomeId);
			break;
		}
		if len(line)>0 && line[0]=='>' {
			//fmt.Printf("%d - %d\n",counter,len(readStr.String()));
			//counter++;
			//fmt.Println(readStr.String());
			//fmt.Println(len(readStr.String()));
			storeSequenceWithN([]byte(readStr.String()), genomeId);
			
			readStr = bytes.Buffer{};
		} else {
			readStr.Write(line);
		}
	}

	

	//return []byte(byte_array.String())
}


func InitializeMatrix(genomeCount int) {
	if IsMatrix {
		var rowCount int;
		if IsSpaced {
			rowCount = (1<<((K1+K2)<<1));//-1; // 2^k-1
		} else {
			rowCount = (1<<(K<<1));//-1; // 2^k-1
		}
		colCount := genomeCount;
		fmt.Println("Row Count=", rowCount, ", Column Count=", colCount);
		
		GkcMatrix = make([][]int, rowCount);
		for i :=0; i<rowCount; i++ {
			GkcMatrix[i] = make([]int, colCount);
		}
	} else {
		GkcMap = make(map[uint32]map[uint16]int);
	}
}

func GetRowCount() int {
	if IsSpaced {
		return (1<<((K1+K2)<<1));//-1; // 2^(k1+k2)
	} else {
		return (1<<(K<<1));//-1; // 2^k
	}
}

func FnaToSequenceDump(filePath string, outputPath string) {
	
	
	fpOut, errOut := os.Create(outputPath);
	Check(errOut);
	wOut := bufio.NewWriter(fpOut);

	
	
	fp, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer fp.Close()
	br := bufio.NewReader(fp)
	var isPrefix bool

	var line []byte
	
	line, isPrefix, err = br.ReadLine()
	if err != nil || isPrefix {
		return;
	}
	readStr := bytes.Buffer{}
	
	for {
		line, isPrefix, err = br.ReadLine()
		
		if err != nil || isPrefix {
			wOut.WriteString(readStr.String());
			wOut.WriteString("\n");
			break;
		}
		if len(line)>0 && line[0]=='>' {
			//fmt.Println(readStr.String());
			//fmt.Println(len(readStr.String()));
			//storeSequenceWithN([]byte(readStr.String()), genomeId);
			//fmt.Println(readStr.String());
			wOut.WriteString(readStr.String());
			wOut.WriteString("\n");
			
			readStr = bytes.Buffer{};
		} else {
			readStr.Write(line);
		}
	}

	fp.Close();
	
	wOut.Flush();
	fpOut.Close();

	//return []byte(byte_array.String())
}

func GetExtFileCount(dirS string, extS string) int {
	files, err := ioutil.ReadDir(dirS);
	Check(err);
	
	count := 0;
	for _, f := range files {
		if(filepath.Ext(f.Name()) == extS) {
			//fmt.Printf("%d --- %s\n", i, f.Name());
			count++;
		}
	}
	return count;
}

func GetExtFileNames(dirS string, extS string) []string {
	files, err := ioutil.ReadDir(dirS);
	Check(err);
	
	fileCount := GetExtFileCount(dirS, extS);
	
	fileNames := make([]string, fileCount);
	
	index := 0;
	for _, f := range files {
		if(filepath.Ext(f.Name()) == extS) {
			//fmt.Printf("%d --- %s\n", index, f.Name());
			fileNames[index] = f.Name();
			index++;
		}
	}
	return fileNames;
}

func GetGenomeIdToFileNameMap(inputDir string) (map[uint16]string, int){
	genomeFileMap := make(map[uint16]string);
	genomeCount := 0;
	var fileName string;
	if IsMatrix {
		fileName = fmt.Sprintf("index.%d-%d-%d.matrix.meta", K1,Sp,K2);
	} else {
		fileName = fmt.Sprintf("index.%d-%d-%d.map.meta", K1,Sp,K2);
	}
	
	fp, err := os.Open(inputDir + "/" + fileName);
	Check(err);
	br := bufio.NewReader(fp);
	
	for {
		line, isPrefix, err := br.ReadLine();
		if err != nil || isPrefix {
			break;
		}
		var i int;
		for i=0; i<len(line); i++ {
			if(line[i]==',') {
				break;
			}
		}
		if i!=len(line) {
			tmp, err := strconv.ParseUint(string(line[0:i]), 10, 32);
			Check(err);
			genomeId := uint16(tmp);
			genomeFileName := string(line[(i+1):len(line)]);
			genomeFileMap[genomeId] = genomeFileName;
			genomeCount++;
			fmt.Printf("%d --- %s\n", genomeId, genomeFileName);
		}
    }
	
	fp.Close();
	return genomeFileMap, genomeCount;
}

func PrintHashMap() {
	for kmerId, gMap := range GkcMap {
		fmt.Print("kmerId:", kmerId);
		for genomeId, count := range gMap {
			fmt.Println(" genomeId:", genomeId, " count:", count);
			//fmt.Println("count:", gkcMap[kmerId][genomeId]);
		}
	}
}