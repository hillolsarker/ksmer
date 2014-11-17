package main 

import (
	"fmt"
//	"math"
//	"io/ioutil"
	"os"
//	"encoding/gob"
	"bufio"
//	"bytes"
	"flag"
//	"path/filepath"
	"strconv"
	"strings"
	
	"mgdb"
)




func main() {
	
	
	//mgdb.FnaToSequenceDump("C:\\dev\\project\\assignments\\05.2014.Fall\\AlgoBioInform\\ksmer\\dataset1_sim\\output\\output-454.fna", "C:\\dev\\project\\assignments\\05.2014.Fall\\AlgoBioInform\\ksmer\\dataset1_sim\\output\\output-454.txt");
	//return;
	var k_value = flag.String("k", "", "value of k in k-mer")
	var spaced_kmer = flag.String("sk", "", "use spaced kmer")
	var input_dir = flag.String("i", "../dataset1_sim/output", "input directory")
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
	
	genomeFileMap, genomeCount := mgdb.GetGenomeIdToFileNameMap(*output_dir);
	//_, genomeCount := mgdb.GetGenomeIdToFileNameMap(*output_dir);
	
	kmerCount := mgdb.GetRowCount();
	
	if !mgdb.IsMatrix {
		fmt.Printf("Map file is not supported right now. Need minor change though\n");
		return;
	}
	
	kmerStr := fmt.Sprintf("%d-%d-%d",mgdb.K1, mgdb.Sp, mgdb.K2);
	rdbFileName := "index." + kmerStr + ".matrix";
	rdb := mgdb.LoadIndexArrayFromFile(*output_dir + "/" + rdbFileName);
	files := mgdb.GetExtFileNames(*input_dir, ".fna");
	for _, fileName := range files {
		mgdb.InitializeMatrix(1);
		mgdb.ParseReadFnaStore(*input_dir + "/" + fileName, 0);
		
		fp, err := os.Create(*output_dir + "/" + fileName + "." + kmerStr + ".csv");
		mgdb.Check(err);
		
		w := bufio.NewWriter(fp);
		w.WriteString("K-mer");
		for i:=0; i<genomeCount; i++ {
			w.WriteString(",");
			w.WriteString(genomeFileMap[uint16(i)]);
		}
		w.WriteString(",b\n");
		
		for k:=0; k<kmerCount; k++ {
			bi := mgdb.GkcMatrix[k][0];
			w.WriteString(fmt.Sprintf("%d", k));
			for genomeId:=0; genomeId<genomeCount; genomeId++ {
				w.WriteString(fmt.Sprintf(",%d", rdb[k][genomeId]));
			}
			w.WriteString(fmt.Sprintf(",%d\n", bi));
		}
		w.Flush();
		fp.Close();
	}
}

