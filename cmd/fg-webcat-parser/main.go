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
type FGGroup struct {
	ID         string
	Name       string
	Categories []FGCategory
}

type FGCategory struct {
	ID    int
	Name  string
	GrpID string
	UTM   string
}

type UTMAction struct {
	Block        string
	Allow        string
	Monitor      string
	Warning      string
	Authenticate string
}

// --------------------------- Main -----------------------------
func main() {
	// init
	txtContent := readTextFile(filepath)
	fgGroupMap, fgCategoryMap := initMapsFromtxt(txtContent)
	utm := UTMAction{
		Block:        "block",
		Allow:        "allow",
		Monitor:      "monitor",
		Warning:      "warning",
		Authenticate: "authenticate",
	}

	// start UI
	fmt.Print(welcomeMessage)
	fmt.Println(inputMessage)

	fgWebConf := readUserInput()
	bcIDs := parseConfig(fgWebConf, fgCategoryMap, utm)

	fmt.Print(blockedCategoriesMessage)

	printCategoryStatus(fgGroupMap, bcIDs, utm.Allow)

	fmt.Print(exitMessage)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// --------------------------- Functions -----------------------------

func initMapsFromtxt(txt []string) (map[string]FGGroup, map[int]FGCategory) {
	// init maps
	mGroup := make(map[string]FGGroup)
	mCategory := make(map[int]FGCategory)

	// define regex
	reGroup := regexp.MustCompile(`^(g\d{2})\s+(.*):$`)
	reCategory := regexp.MustCompile(`^(\d*)(\s.*)$`)

	var currentGroup *FGGroup

	for _, l := range txt {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		if match := reGroup.FindStringSubmatch(l); match != nil {
			if currentGroup != nil {
				mGroup[currentGroup.ID] = *currentGroup
			}

			currentGroup = &FGGroup{
				ID:         match[1],
				Name:       match[2],
				Categories: []FGCategory{},
			}
		} else if match := reCategory.FindStringSubmatch(l); match != nil {
			if currentGroup != nil {
				i, _ := strconv.Atoi(match[1])
				category := FGCategory{
					ID:    i,
					Name:  match[2],
					GrpID: currentGroup.ID,
				}
				mCategory[i] = category
				currentGroup.Categories = append(currentGroup.Categories, category)
			}
		}
	}
	// Add last group after exiting loop
	if currentGroup != nil {
		mGroup[currentGroup.ID] = *currentGroup
	}

	return mGroup, mCategory
}

func printCategoryStatus(mGroup map[string]FGGroup, categories map[int]FGCategory, status string) {
	// make temp map to group categories by group ID
	gc := make(map[string][]FGCategory)
	for _, c := range categories {
		if c.UTM == status {
			gc[c.GrpID] = append(gc[c.GrpID], c)
		}
	}

	if len(gc) == 0 {
		fmt.Println("No categories of this status found")
	} else {
		// Print groups and categories
		for gID, cs := range gc {
			if g, ok := mGroup[gID]; ok {
				fmt.Println(g.Name)
				for _, c := range cs {
					fmt.Println("    ", c.Name)
				}
			}
		}
	}
}

func parseConfig(conf []string, mCategory map[int]FGCategory, utm UTMAction) map[int]FGCategory {
	var cID int
	var lastLine string
	var action string
	var a bool

	cs := make(map[int]FGCategory)

	// Look through config and find what utm action is set on them
	for _, l := range conf {
		l := strings.TrimSpace(l)

		if strings.HasPrefix(l, "set category") {
			c := strings.Fields(l)
			cID, _ = strconv.Atoi(c[len(c)-1])
		} else if strings.HasPrefix(l, "set action") {
			action = strings.Fields(l)[2]
			a = true
		} else if strings.HasPrefix(lastLine, "set category") && (l == "next" || lastLine == "set log disable") {
			action = utm.Monitor
			a = true
		}

		if a {
			c := FGCategory{
				ID:    cID,
				Name:  mCategory[cID].Name,
				GrpID: mCategory[cID].GrpID,
				UTM:   action,
			}
			cs[cID] = c
		}

		lastLine = l
		a = false
	}

	// find all categories that was not found in config. These are set to "allow"
	for i, category := range mCategory {
		if _, ok := cs[i]; !ok {
			c := FGCategory{
				ID:    i,
				Name:  category.Name,
				GrpID: category.GrpID,
				UTM:   utm.Allow,
			}
			cs[i] = c
		}
	}

	return cs
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
