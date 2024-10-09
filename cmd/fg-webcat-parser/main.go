package main

import (
	"fmt"
	"strconv"
	"strings"
)

const fgCategories string = `1 Drug Abuse
3 Hacking
4 Illegal or Unethical
5 Discrimination
6 Explicit Violence
12 Extremist Groups
59 Proxy Avoidance
62 Plagiarism
83 Child Sexual Abuse
96 Terrorism
98 Crypto Mining
99 Potentially Unwanted Program
2 Alternative Beliefs
7 Abortion
8 Other Adult Materials
9 Advocacy Organizations
11 Gambling
13 Nudity and Risque
14 Pornography
15 Dating
16 Weapons (Sales)
57 Marijuana
63 Sex Education
64 Alcohol
65 Tobacco
66 Lingerie and Swimsuit
67 Sports Hunting and War Games
19 Freeware and Software Downloads
24 File Sharing and Storage
25 Streaming Media and Download
72 Peer-to-peer File Sharing
75 Internet Radio and TV
76 Internet Telephony
26 Malicious Websites
61 Phishing
86 Spam URLs
88 Dynamic DNS
90 Newly Observed Domain
91 Newly Registered Domain
17 Advertising
18 Brokerage and Trading
20 Games
23 Web-based Email
28 Entertainment
29 Arts and Culture
30 Education
33 Health and Wellness
34 Job Search
35 Medicine
36 News and Media
37 Social Networking
38 Political Organizations
39 Reference
40 Global Religion
42 Shopping
44 Society and Lifestyles
46 Sports
47 Travel
48 Personal Vehicles
54 Dynamic Content
55 Meaningless Content
58 Folklore
68 Web Chat
69 Instant Messaging
70 Newsgroups and Message Boards
71 Digital Postcards
77 Child Education
78 Real Estate
79 Restaurant and Dining
80 Personal Websites and Blogs
82 Content Servers
85 Domain Parking
87 Personal Privacy
89 Auction
31 Finance and Banking
41 Search Engines and Portals
43 General Organizations
49 Business
50 Information and Computer Security
51 Government and Legal Organizations
52 Information Technology
53 Armed Forces
56 Web Hosting
81 Secure Websites
84 Web-based Applications
92 Charitable Organizations
93 Remote Access
94 Web Analytics
95 Online Meeting
97 URL Shortening
100 Artificial Intelligence Technology
101 Cryptocurrency`

const fgWebConf string = `config webfilter profile
    edit "CF_SOCIAL"
        config ftgd-wf
            set options error-allow
            config filters
                edit 24
                    set category 24
                next
                edit 25
                    set category 25
                    set action block
                next
                edit 26
                    set category 26
                    set action block
                next
                edit 28
                    set category 28
                next
                edit 29
                    set category 29
                next
                edit 6
                    set category 6
                    set action block
                next
                edit 30
                    set category 30
                next
                edit 31
                    set category 31
                next
                edit 33
                    set category 33
                next
                edit 2
                    set category 2
                next
                edit 4
                    set category 4
                    set action block
                next
                edit 14
                    set category 14
                    set action block
                next
                edit 1
                    set category 1
                    set action block
                next
                edit 9
                    set category 9
                next
                edit 5
                    set category 5
                    set action block
                next
                edit 16
                    set category 16
                next
                edit 12
                    set category 12
                    set action block
                next
                edit 7
                    set category 7
                next
                edit 3
                    set category 3
                    set action block
                next
                edit 8
                    set category 8
                    set action block
                next
                edit 13
                    set category 13
                    set action block
                next
                edit 17
                    set category 17
                    set action block
                next
                edit 15
                    set category 15
                    set action block
                next
                edit 34
                    set category 34
                next
                edit 35
                    set category 35
                next
                edit 36
                    set category 36
                next
                edit 37
                    set category 37
                next
                edit 38
                    set category 38
                next
                edit 39
                    set category 39
                next
                edit 18
                    set category 18
                next
                edit 19
                    set category 19
                    set action block
                next
                edit 20
                    set category 20
                    set action block
                next
                edit 40
                    set category 40
                next
                edit 41
                    set category 41
                next
                edit 43
                    set category 43
                next
                edit 44
                    set category 44
                next
                edit 46
                    set category 46
                next
                edit 47
                    set category 47
                next
                edit 48
                    set category 48
                next
                edit 49
                    set category 49
                next
                edit 50
                    set category 50
                next
                edit 51
                    set category 51
                next
                edit 52
                    set category 52
                next
                edit 53
                    set category 53
                next
                edit 54
                    set category 54
                next
                edit 55
                    set category 55
                    set action block
                next
                edit 56
                    set category 56
                next
                edit 57
                    set category 57
                    set action block
                next
                edit 58
                    set category 58
                    set action block
                next
                edit 59
                    set category 59
                    set action block
                next
                edit 61
                    set category 61
                    set action block
                next
                edit 62
                    set category 62
                    set action block
                next
                edit 63
                    set category 63
                    set action block
                next
                edit 64
                    set category 64
                    set action block
                next
                edit 65
                    set category 65
                next
                edit 66
                    set category 66
                    set action block
                next
                edit 67
                    set category 67
                next
                edit 68
                    set category 68
                next
                edit 69
                    set category 69
                next
                edit 70
                    set category 70
                next
                edit 71
                    set category 71
                next
                edit 72
                    set category 72
                    set action block
                next
                edit 75
                    set category 75
                    set action block
                next
                edit 76
                    set category 76
                next
                edit 78
                    set category 78
                next
                edit 79
                    set category 79
                next
                edit 80
                    set category 80
                next
                edit 81
                    set category 81
                next
                edit 82
                    set category 82
                next
                edit 85
                    set category 85
                next
                edit 86
                    set category 86
                    set action block
                next
                edit 89
                    set category 89
                next
                edit 83
                    set category 83
                    set action block
                next
                edit 11
                    set category 11
                    set action block
                next
                edit 23
                    set category 23
                next
                edit 42
                    set category 42
                next
                edit 77
                    set category 77
                next
                edit 87
                    set category 87
                next
                edit 84
                    set category 84
                next
                edit 88
                next
                edit 140
                    set category 140
                next
                edit 141
                    set category 141
                next
            end
        end
    next
end`

func main() {
	cm := makeCategoryMap(fgCategories)
	bcIDs := findBlockedCategoryIDs(fgWebConf)
	allBlockedCats := getBlockedCats(cm, bcIDs)

	for _, catname := range allBlockedCats {
		fmt.Println(catname)
	}
}

func getBlockedCats(fgCatMap map[int]string, blockedCats []int) []string {
	var catnames []string
	for _, cat := range blockedCats {
		catnames = append(catnames, fgCatMap[cat])
	}

	return catnames
}

func makeCategoryMap(fgc string) map[int]string {
	m := make(map[int]string)

	for _, cat := range strings.Split(fgc, "\n") {
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

func findBlockedCategoryIDs(fgWebConf string) []int {
	var bc []int
	var catnum int
	var err error
	for _, l := range strings.Split(fgWebConf, "\n") {
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
