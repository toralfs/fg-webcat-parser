package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// --------------------------- Consts -----------------------------
const welcomeMessage string = `
---------------------------------------------------------
Welcome to FG webcategory parser
		
Please enter the webfilter profile configuration
and the program will extract the blocked categories. 

File path needs to be: ../../assets/fg-category-list.txt
---------------------------------------------------------
`
const inputMessage string = "Enter the configuration snippet and press enter."
const blockedCategoriesMessage string = `
---------------------------------------------------------
Format is as below:
Category Group
    Category Name

The blocked categories are: 
---------------------------------------------------------
`
const exitMessage string = `
---------------------------------------------------------
Press enter to close the program
`
const filepath string = "../../assets/fg-category-list.txt"

// --------------------------- Structs -----------------------------
type fgCatGroup struct {
	GrpID   string
	GrpName string
	Cats    map[int]fgCategory
}

type fgCategory struct {
	CatID   int
	CatName string
}

// --------------------------- Main -----------------------------
func main() {
	fmt.Print(welcomeMessage)
	fmt.Println(inputMessage)
	fgWebConf := readUserInput()
	fgWebCats := readTextFile(filepath)

	cm := makeCategoryMap(fgWebCats)
	bcIDs := findBlockedCategoryIDs(fgWebConf)

	fmt.Print(blockedCategoriesMessage)

	printBlockedCats(cm, bcIDs)

	fmt.Print(exitMessage)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// --------------------------- Functions -----------------------------
func printBlockedCats(fgCatMap []fgCatGroup, blockedCats []int) {
	for _, group := range fgCatMap {
		var cNames []string
		for _, cat := range blockedCats {
			n, ok := group.Cats[cat]
			if ok {
				cNames = append(cNames, n.CatName)
			}
		}
		if len(cNames) > 0 {
			fmt.Println(group.GrpName)
			for _, n := range cNames {
				fmt.Println("    ", n)
			}
		}
	}
}

func makeCategoryMap(fgc []string) []fgCatGroup {
	// Create a regex to match the "gXX Group Name" format
	groupRegex := regexp.MustCompile(`^(g\d{2})\s+(.*):$`)

	var groups []fgCatGroup
	var currentGroup *fgCatGroup

	for _, line := range fgc {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if match := groupRegex.FindStringSubmatch(line); match != nil {
			// If we have an existing group, save it before starting a new one
			if currentGroup != nil {
				groups = append(groups, *currentGroup)
			}

			// Create a new group with GrpID and GrpName
			currentGroup = &fgCatGroup{
				GrpID:   match[1],
				GrpName: match[2],
				Cats:    map[int]fgCategory{},
			}
		} else if currentGroup != nil {
			fields := strings.Fields(line)
			catID, _ := strconv.Atoi(fields[0])
			catName := strings.Join(fields[1:], " ")
			category := fgCategory{
				CatID:   catID,
				CatName: catName,
			}

			currentGroup.Cats[catID] = category
		}
	}

	// Append the last group after exiting the loop
	if currentGroup != nil {
		groups = append(groups, *currentGroup)
	}

	return groups
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
