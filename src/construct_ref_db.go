package main 

import (
	"fmt"
	"os"
	"bufio"
//	"io/ioutil"
//	"encoding/gob"
//	"bytes"
	"flag"
//	"path/filepath"
	"strconv"
	"strings"
	
	"mgdb"

)

func main() {
	
	
	var k_value = flag.String("k", "", "value of k in k-mer")
	var spaced_kmer = flag.String("sk", "", "use spaced kmer")
	var input_dir = flag.String("i", "../Bioinformatics-Project-Genomes/Dataset1_3Genomes", "input directory")
	var output_dir = flag.String("o", "../output", "output directory")
	flag.Parse();
	
	//var input_dir = "C:\\dev\\project\\assignments\\05.2014.Fall\\AlgoBioInform\\ksmer\\input\\fasta"; //"../input/";
	
	fmt.Println("Input Directory = ", *input_dir);
	fmt.Println("Output Directory = ", *output_dir);
	
	if(*k_value != "") {
		// catch multi-value return
		catchk, err := strconv.ParseUint(*k_value, 10, 32);
		mgdb.Check(err);
						
		// assign k value to k
		mgdb.K = uint32(catchk);
		mgdb.IsSpaced = false;
		
		fmt.Printf("k=%d\n", mgdb.K);
	}
	if(*spaced_kmer != "") {
		mgdb.IsSpaced = true;
		temp := strings.Split(*spaced_kmer, ",");
				
		catchk1, err := strconv.ParseUint(temp[0], 10, 32);
		mgdb.Check(err);
		catchsp, err := strconv.ParseUint(temp[1], 10, 32);
		mgdb.Check(err);
		catchk2, err := strconv.ParseUint(temp[2], 10, 32);
		mgdb.Check(err);
		
		mgdb.K1 = uint32(catchk1);
		mgdb.Sp = uint32(catchsp);
		mgdb.K2 = uint32(catchk2);
		fmt.Printf("k1=%d, space=%d, k2=%d\n", mgdb.K1, mgdb.Sp, mgdb.K2);
	}
	
	
	
	// Read list of input files into FileInfo slice named files
	//files, err := ioutil.ReadDir(*input_dir);
	//mgdb.Check(err);
	files := mgdb.GetExtFileNames(*input_dir, ".fasta");
	
	mgdb.InitializeMatrix(len(files));
	
	for i, fileName := range files {
		fmt.Printf("%d --- %s\n", i, fileName);
		mgdb.ParseReadFnaStore(*input_dir + "/" +fileName, uint16(i));
	}
	
	var dbSuffixFileName = "";
	if mgdb.IsSpaced {
		dbSuffixFileName = fmt.Sprintf(".%d-%d-%d", mgdb.K1, mgdb.Sp, mgdb.K2);
	} else {
		dbSuffixFileName = fmt.Sprintf(".%d-0-0", mgdb.K);
	}
	
	var indexFileName string;
	if mgdb.IsMatrix {
		indexFileName = "index"+dbSuffixFileName+".matrix";
		mgdb.SaveIndexMapToFile(*output_dir + "/" + indexFileName, mgdb.GkcMatrix);
	} else {
		indexFileName = "index"+dbSuffixFileName+".map";
		mgdb.SaveIndexMapToFile(*output_dir + "/" + indexFileName, mgdb.GkcMap);
	}
	
	fpMeta, errOut := os.Create(*output_dir + "/" + indexFileName + ".meta");
	mgdb.Check(errOut);
	wMeta := bufio.NewWriter(fpMeta);
	for i, fileName := range files {
		wMeta.WriteString(fmt.Sprintf("%d,%s\n", i, fileName));
	}
	wMeta.Flush();
	fpMeta.Close();

}

