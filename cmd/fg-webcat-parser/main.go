package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// --------------------------- Consts -----------------------------
const fileIn string = "fg-webcat-parser/fg-category-list.txt"
const assetsPathEnv string = "GO_ASSETS"

// --------------------------- Structs -----------------------------
type FGGroup struct {
	ID         int
	Name       string
	Categories []FGCategory
}

type FGCategory struct {
	ID    int
	Name  string
	GrpID int
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
	assetsPath, err := initEnv(assetsPathEnv)
	if err != nil {
		fmt.Println(err, "\nExiting program...")
		os.Exit(0)
	}
	catListPath := filepath.Join(assetsPath, fileIn)
	txtContent := readTextFile(catListPath)
	fgGroupMap, fgCategoryMap := initMapsFromtxt(txtContent)
	utm := UTMAction{
		Block:        "block",
		Allow:        "allow",
		Monitor:      "monitor",
		Warning:      "warning",
		Authenticate: "authenticate",
	}

	// Welcome message
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf("------------------------- FortiGate webfilter parser -------------------------\n")
	fmt.Printf("------------------------------------------------------------------------------\n")
	fmt.Printf("Program will take the configuration snippet of a FortiGate webfilter profile,\n")
	fmt.Printf("parse it and return a view of the UTM status of available categories.\n")
	fmt.Printf("------------------------------------------------------------------------------\n\n")

	fmt.Println("Enter the configuration snippet and press Ctrl+D (or Ctrl+Z if using Windows).")

	// Read and parse configuration from user input
	confWebFilterProfile := readUserInput()
	confedCategories := parseConfig(confWebFilterProfile, fgCategoryMap, utm)

	// Read UTM status from user input and show result
	fmt.Print(`
------------------------------------------------------------
Config parsed.
`)

	// loop the selection until user exits
	for {
		fmt.Print(`
Select which UTM status:
1 - `, utm.Allow, `
2 - `, utm.Block, `
3 - `, utm.Monitor, `
4 - `, utm.Warning, `
5 - `, utm.Authenticate, `
0 - Exit program
`)

		utmStatus := readUserInputSingle()
		switch utmStatus {
		case "1":
			printCategoryStatus(fgGroupMap, confedCategories, utm.Allow)
		case "2":
			printCategoryStatus(fgGroupMap, confedCategories, utm.Block)
		case "3":
			printCategoryStatus(fgGroupMap, confedCategories, utm.Monitor)
		case "4":
			printCategoryStatus(fgGroupMap, confedCategories, utm.Warning)
		case "5":
			printCategoryStatus(fgGroupMap, confedCategories, utm.Authenticate)
		case "0":
			fmt.Println("Good bye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}

		fmt.Printf("------------------------------------------------------------\nPress enter to go back to UTM selection.\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

// --------------------------- Functions -----------------------------

func initEnv(envName string) (string, error) {
	env := os.Getenv(envName)
	var err error = nil
	if env == "" {
		err = fmt.Errorf("environment variable: \"%s\" is not set", envName)
	}
	return env, err
}

func initMapsFromtxt(txt []string) (map[int]FGGroup, map[int]FGCategory) {
	// init maps
	mGroup := make(map[int]FGGroup)
	mCategory := make(map[int]FGCategory)

	// define regex
	reGroup := regexp.MustCompile(`^(g\d{2})\s+(.*)$`)
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
			i, _ := strconv.Atoi(strings.Split(match[1], "g")[1])
			currentGroup = &FGGroup{
				ID:         i,
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

func printCategoryStatus(mGroup map[int]FGGroup, categories map[int]FGCategory, status string) {
	fmt.Print(`
------------------------------------------------------------
Format is as below:
Category Group
	Category Name

The `, status, ` categories are: 
------------------------------------------------------------
`)

	// make temp map to group categories by group ID
	gc := make(map[int][]FGCategory)
	for _, c := range categories {
		if c.UTM == status {
			gc[c.GrpID] = append(gc[c.GrpID], c)
		}
	}

	// return early if no categories found
	if len(gc) == 0 {
		fmt.Println("No categories of this status found")
		return
	}

	// sort map keys so we can print category groups in order
	keys := make([]int, 0, len(gc))
	for k := range gc {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Print groups and categories
	for _, k := range keys {
		if g, exist := mGroup[k]; exist {
			fmt.Println(g.Name)
			for _, c := range gc[k] {
				fmt.Println("    ", c.Name)
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
		if !s.Scan() {
			break
		}
		lines = append(lines, s.Text())
	}

	err := s.Err()
	if err != nil {
		log.Fatal(err)
	}

	return lines
}

func readUserInputSingle() string {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	ln := s.Text()
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return ln
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
