package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const welcomeMessage string = `
---------------------------------------------------------
Welcome to FG webcategory parser
		
Please enter the webfilter profile configuration
and the program will extract the blocked categories. 

File path needs to be: ../../assets/fg-category-list.txt
---------------------------------------------------------
`

const inputMessage string = "Enter the configuration snippet and press enter."

func main() {
	fmt.Print(welcomeMessage)
	fmt.Println(inputMessage)
	fgWebConf := readUserInput()
	fgWebCats := readTextFile("../../assets/fg-category-list.txt")

	cm := makeCategoryMap(fgWebCats)
	bcIDs := findBlockedCategoryIDs(fgWebConf)
	allBlockedCats := getBlockedCats(cm, bcIDs)

	fmt.Println("---------------------------------------------------------")
	fmt.Println("The blocked categories are: ")
	fmt.Println("---------------------------------------------------------")
	for _, catname := range allBlockedCats {
		fmt.Println(catname)
	}

	fmt.Printf("\n Press enter to close the program \n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getBlockedCats(fgCatMap map[int]string, blockedCats []int) []string {
	var catnames []string
	for _, cat := range blockedCats {
		catnames = append(catnames, fgCatMap[cat])
	}

	return catnames
}

func makeCategoryMap(fgc []string) map[int]string {
	m := make(map[int]string)

	for _, cat := range fgc {
		f := strings.SplitN(cat, " ", 2)

		i, err := strconv.Atoi(f[0])
		if err == nil {
			m[i] = f[1]
		} else {
			fmt.Println("Failed to convert: ", i, "to int", err)
		}
	}
	return m
}

func findBlockedCategoryIDs(fgWebConf []string) []int {
	var bc []int
	var catnum int
	var err error
	for _, l := range fgWebConf {
		lTrimmed := strings.TrimSpace(l)
		if strings.HasPrefix(lTrimmed, "set category") {
			cat := strings.Fields(lTrimmed)
			catnum, err = strconv.Atoi(cat[len(cat)-1])
			if err != nil {
				fmt.Println("Failed to conver to int", err)
			}
		}

		if strings.HasPrefix(lTrimmed, "set action block") {
			bc = append(bc, catnum)
		}
	}
	return bc
}

func readUserInput() []string {
	s := bufio.NewScanner(os.Stdin)

	var lines []string
	for {
		s.Scan()
		l := s.Text()
		if len(l) == 0 {
			break
		}
		lines = append(lines, l)
	}

	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}

	return lines
}

func readTextFile(path string) []string {
	var fileContent []string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(file)

	for s.Scan() {
		fileContent = append(fileContent, s.Text())
	}

	return fileContent
}
