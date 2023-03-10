package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func swap(array []int, indexA, indexB int) {
	t := array[indexA]
	array[indexA] = array[indexB]
	array[indexB] = t
}

type SorterQuick struct{}

func (s *SorterQuick) split(array []int, L, R int) int {
	P := array[R]
	m := L - 1
	for i := L; i <= R; i++ {
		if array[i] <= P {
			m++
			swap(array, m, i)
		}
	}
	return m
}

func (s *SorterQuick) qsort(array []int, L, R int) {
	if L >= R {
		return
	}
	M := s.split(array, L, R)
	s.qsort(array, L, M-1)
	s.qsort(array, M+1, R)
}

func (s *SorterQuick) Sort(in []int) []int {
	N := len(in)
	s.qsort(in, 0, N-1)
	return in
}

func (s *SorterQuick) Name() string {
	return "Quick"
}

type SorterMerge struct{}

func (s *SorterMerge) merge(array []int, L, M, R int) {
	T := make([]int, R-L+1)
	a := L
	b := M + 1
	t := 0

	for a <= M && b <= R {
		if array[a] <= array[b] {
			T[t] = array[a]
			a++
		} else {
			T[t] = array[b]
			b++
		}
		t++
	}
	for a <= M {
		T[t] = array[a]
		t++
		a++
	}
	for b <= R {
		T[t] = array[b]
		t++
		b++
	}
	for i := L; i <= R; i++ {
		array[i] = T[i-L]
	}
}
func (s *SorterMerge) msort(array []int, L, R int) {
	if L >= R {
		return
	}
	M := (L + R) / 2
	s.msort(array, L, M)
	s.msort(array, M+1, R)
	s.merge(array, L, M, R)
}

func (s *SorterMerge) Sort(in []int) []int {
	N := len(in)
	s.msort(in, 0, N-1)
	return in
}

func (s *SorterMerge) Name() string {
	return "Merge"
}

//................External sort..........................

type SorterExternal struct {
	testFile       *os.File
	limitArraySize int
	fileA          *os.File
	fileB          *os.File
	fileC          *os.File
	fileD          *os.File
}

func (s *SorterExternal) testFileGenerate(numLine int, maxNumber int) *os.File {
	file, err := ioutil.TempFile("", "otus")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < numLine; i++ {
		file.WriteString(strconv.Itoa(rand.Intn(maxNumber)) + "\n")
	}
	return file
}

func (s *SorterExternal) Name() string {
	return "External"
}

func (s *SorterExternal) GenerateTestData(numLine int, maxNumber int) {
	s.testFile = s.testFileGenerate(numLine, maxNumber)
}

func (s *SorterExternal) RemoveTestData() {
	os.Remove(s.testFile.Name())
}

func (s *SorterExternal) RemoveTempFile() {
	os.Remove(s.fileA.Name())
	os.Remove(s.fileB.Name())
	os.Remove(s.fileC.Name())
	os.Remove(s.fileD.Name())
}

func (s *SorterExternal) CreateTempFile() {
	var err error
	s.fileA, err = ioutil.TempFile("", "otus")
	if err != nil {
		log.Fatal(err)
	}
	s.fileB, err = ioutil.TempFile("", "otus")
	if err != nil {
		log.Fatal(err)
	}
	s.fileC, err = ioutil.TempFile("", "otus")
	if err != nil {
		log.Fatal(err)
	}
	s.fileD, err = ioutil.TempFile("", "otus")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SorterExternal) fileToArray(reader *bufio.Reader) ([]int, bool) {
	result := make([]int, 0, s.limitArraySize)

	for i := 0; i < s.limitArraySize; i++ {
		l, _, err := reader.ReadLine()
		if err != nil {
			return result, true
		}
		intVar, err := strconv.Atoi(string(l))
		if err != nil {
			return result, true
		}
		result = append(result, intVar)
	}
	return result, false
}

func (s *SorterExternal) arrayToFile(array []int, file *os.File) {
	for _, el := range array {
		file.WriteString(strconv.Itoa(el) + "\n")
	}
}

func (s *SorterExternal) splitSrcDataByTwoFileAB() {
	isEnd := false
	var ar []int
	index := 0
	s.testFile.Seek(0, 0)
	os.Truncate(s.fileA.Name(), 0)
	s.fileA.Seek(0, 0)
	os.Truncate(s.fileB.Name(), 0)
	s.fileB.Seek(0, 0)
	reader := bufio.NewReader(s.testFile)
	for !isEnd {
		ar, isEnd = s.fileToArray(reader)
		if len(ar) == 0 {
			break
		}
		sort.Ints(ar)
		if index%2 == 0 {
			s.arrayToFile(ar, s.fileA)
		} else {
			s.arrayToFile(ar, s.fileB)
		}
		index++
	}
}

func (s *SorterExternal) oneMerge(srcFile1, srcFile2, dstFile1, dstFile2 *os.File) {
	readerOne := bufio.NewReader(srcFile1)
	readerTwo := bufio.NewReader(srcFile2)
	writerOne := bufio.NewWriter(dstFile1)
	writerTwo := bufio.NewWriter(dstFile2)
	srcFile1.Seek(0, 0)
	srcFile2.Seek(0, 0)
	dstFile1.Seek(0, 0)
	dstFile2.Seek(0, 0)
	var valA, valB, oldValA, oldValB int
	isNextA := true
	isNextB := true
	count := 0
	outputWriter := writerOne
	currentWriterIsOne := true
	isBegin := true
	isFileAEnd := false
	isFileBEnd := false
	isChangeA := false
	isChangeB := false
	for {
		if isNextA && !isFileAEnd {
			strA, _, err := readerOne.ReadLine()
			if err != nil {
				isFileAEnd = true
			}
			valA, _ = strconv.Atoi(string(strA))
		}
		if isNextB && !isFileBEnd {
			strB, _, err := readerTwo.ReadLine()
			if err != nil {
				isFileBEnd = true
			}
			valB, _ = strconv.Atoi(string(strB))
		}

		if isFileBEnd && isFileAEnd {
			break
		}

		if oldValA > valA {
			isChangeA = true
		}
		if oldValB > valB {
			isChangeB = true
		}

		if isChangeA && isChangeB && !isBegin {
			if currentWriterIsOne {
				outputWriter = writerTwo
			} else {
				outputWriter = writerOne
			}
			currentWriterIsOne = !currentWriterIsOne
			if !isFileBEnd && !isFileAEnd {
				isChangeA = false
				isChangeB = false
			} else if isFileBEnd && !isFileAEnd {
				isChangeA = false
			} else if !isFileBEnd && isFileAEnd {
				isChangeB = false
			}
		}

		if !isChangeA && !isChangeB {
			if valA <= valB {
				outputWriter.WriteString(strconv.Itoa(valA) + "\n")
				isNextA = true
				isNextB = false
			} else {
				outputWriter.WriteString(strconv.Itoa(valB) + "\n")
				isNextA = false
				isNextB = true
			}
		} else if isChangeA && !isChangeB {
			outputWriter.WriteString(strconv.Itoa(valB) + "\n")
			isNextA = false
			isNextB = true
		} else if isChangeB && !isChangeA {
			outputWriter.WriteString(strconv.Itoa(valA) + "\n")
			isNextA = true
			isNextB = false
		} else if isChangeA && isChangeB {
			panic("Error")
		}

		oldValA = valA
		oldValB = valB
		isBegin = false
		count++
	}
	writerOne.Flush()
	writerTwo.Flush()
	os.Truncate(srcFile1.Name(), 0)
	os.Truncate(srcFile2.Name(), 0)
}

func (s *SorterExternal) Sort() []int {
	s.CreateTempFile()
	s.splitSrcDataByTwoFileAB()
	fileReadOne := s.fileA
	fileReadTwo := s.fileB
	fileWriteOne := s.fileC
	fileWriteTwo := s.fileD

	iterateCount := 0
	for {
		s.oneMerge(fileReadOne, fileReadTwo, fileWriteOne, fileWriteTwo)
		iterateCount++
		tempReadOne := fileReadOne
		fileReadOne = fileWriteOne
		fileWriteOne = tempReadOne
		tempReadTwo := fileReadTwo
		fileReadTwo = fileWriteTwo
		fileWriteTwo = tempReadTwo
		f1, _ := os.Stat(fileReadOne.Name())
		f2, _ := os.Stat(fileReadTwo.Name())
		if f1.Size() == 0 || f2.Size() == 0 {
			break
		}
	}
	fmt.Println("Iterate count = ", iterateCount)

	return nil
}

func (s *SorterExternal) PrintSrcFiles() {
	countSrc := 0
	s.testFile.Seek(0, 0)
	scanner := bufio.NewScanner(s.testFile)
	for scanner.Scan() {
		countSrc++
		fmt.Printf(string(scanner.Text()) + " ")
	}
	fmt.Printf("\n")
	fmt.Println("Count in Src File", countSrc)
}

func (s *SorterExternal) PrintDstFiles() {
	fmt.Println("File A")
	s.fileA.Seek(0, 0)
	scanner := bufio.NewScanner(s.fileA)
	countA := 0
	for scanner.Scan() {
		countA++
		fmt.Printf(string(scanner.Text()) + " ")
	}
	fmt.Printf("\n")
	fmt.Println("Count in File A", countA)

	fmt.Println("File B")
	s.fileB.Seek(0, 0)
	scanner = bufio.NewScanner(s.fileB)
	countB := 0
	for scanner.Scan() {
		countB++
		fmt.Printf(string(scanner.Text()) + " ")
	}
	fmt.Printf("\n")
	fmt.Println("Count in File B", countB)

	fmt.Println("File C")
	s.fileC.Seek(0, 0)
	scanner = bufio.NewScanner(s.fileC)
	countC := 0
	for scanner.Scan() {
		countC++
		fmt.Printf(string(scanner.Text()) + " ")
	}
	fmt.Printf("\n")
	fmt.Println("Count in File C", countC)
	fmt.Println("File D")
	s.fileD.Seek(0, 0)
	scanner = bufio.NewScanner(s.fileD)
	countD := 0
	for scanner.Scan() {
		countD++
		fmt.Printf(string(scanner.Text()) + " ")
	}
	fmt.Printf("\n")
	fmt.Println("Count in File D", countD)
}

//.......................................................

var dir string
var isExternal bool

func init() {
	flag.StringVar(&dir, "dir", "", "dir tests")
	flag.BoolVar(&isExternal, "e", false, "is external sort")
}

type TestData struct {
	input, output []int
	sizeArray     int
}

func main() {
	flag.Parse()
	if isExternal {
		sortExt := &SorterExternal{limitArraySize: 100}
		sortExt.GenerateTestData(300, 10)
		defer sortExt.RemoveTestData()
		defer sortExt.RemoveTempFile()
		fmt.Println("Print src array")
		sortExt.PrintSrcFiles()
		start := time.Now()
		sortExt.Sort()
		elapsed := time.Since(start)
		fmt.Println("\n Time - ", elapsed)
		fmt.Println("\n Print result array")
		sortExt.PrintDstFiles()
	} else {
		listFolder := []string{"0.random", "1.digits", "2.sorted", "3.revers"}

		listSorterAlgo := []ISorter{&SorterMerge{}}
		for _, lf := range listFolder {
			log.Printf("Test folder - %s \n", lf)
			testData, err := readTestData(dir+"\\"+lf, 7)
			if err != nil {
				panic(err)
			}
			for _, alg := range listSorterAlgo {
				err = runTests(testData, alg)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
