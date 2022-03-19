package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type DictInfo struct {
	text     []string
	chrMap   map[string]int
	wordList PairList
}

type Pair struct {
	Key   string
	Value int
}

type Letter struct {
	pos     int
	match   string
	include []string
	exclude []string
}

type LetterList []Letter

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func countLetterFreqInSlice(textSlice []string) map[string]int {
	//charMap := PairList{}
	tempMap := make(map[string]int)
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	for _, word := range textSlice {
		word = removeDups(word)
		letters := stringToSlice(word)
		for _, letter := range letters {
			if isAlpha(letter) {
				tempMap[letter]++
			}
		}
	}

	return tempMap
}

func countDictionary(textSlice []string) DictInfo {
	if len(textSlice) <= 0 {
		textSlice = readDictionaryToSlice()
	}

	//fmt.Println("textslice: ", textSlice)
	chrMap := countLetterFreqInSlice(textSlice)

	scores := scoreSlice(textSlice, chrMap)
	sorted := sortMapByValue(scores)

	return DictInfo{text: textSlice, chrMap: chrMap, wordList: sorted}
}

func splitStringByFiveLetterWords(text string) []string {
	re := regexp.MustCompile(`[a-z]{5}`)
	return re.FindAllString(text, -1)
}

func readDictionaryToSlice() []string {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("new-dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	return splitStringByFiveLetterWords(text)
}

func score(text string, chrMap map[string]int) int {
	chars := getUniqueSliceFromText(text)
	score := int(0)
	for _, chr := range chars {
		score += chrMap[chr]
	}
	return score
}

func getUniqueSliceFromText(text string) []string {
	uniques := removeDups(text)
	return stringToSlice(uniques)
}

func stringToSlice(text string) []string {
	return strings.Split(text, "")
}

func removeDups(s string) string {
	var out bytes.Buffer
	letters := make(map[string]int)
	chars := strings.Split(s, "")
	for _, chr := range chars {
		_, ok := letters[chr]
		if !ok {
			letters[chr] = 1
			out.WriteString(chr)
		}
	}
	return out.String()
}

func scoreSlice(s []string, chrMap map[string]int) map[string]int {
	scores := make(map[string]int)
	//fmt.Println("strings:", s)
	for _, text := range s {
		scores[text] = score(text, chrMap)
		//fmt.Println("text: ", text, " score: ", scores[text])
	}
	return scores
}

func sortMapByValue(inMap map[string]int) PairList {
	pl := make(PairList, len(inMap))
	i := 0
	for k, v := range inMap {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))

	return pl
}

func findMatch(textSlice []string, letter string, pos int) []string {
	var matches []string
	for _, word := range textSlice {
		letters := stringToSlice(word)
		if letters[pos] == letter {
			matches = append(matches, word)
		}
	}
	return matches
}

func findInclude(textSlice []string, letter string, pos int) []string {
	var matches []string
	for _, word := range textSlice {
		i := strings.Index(word, letter)
		if strings.Contains(word, letter) && i != pos {
			matches = append(matches, word)
		}
	}
	return matches
}

func findExclude(textSlice []string, letter string) []string {
	var matches []string
	for _, v := range textSlice {
		if !strings.Contains(v, letter) {
			matches = append(matches, v)
		}
	}
	return matches
}

func printGuess(guessSlice []string, word map[int]string, pos int) string {
	var guessBuffer bytes.Buffer
	var carrotBuffer bytes.Buffer
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	guessBuffer.WriteString("[ ")
	carrotBuffer.WriteString("  ")

	i := 0
	for _, v := range guessSlice {
		if word[i] == v {
			guessBuffer.WriteString(colorGreen)
			guessBuffer.WriteString(v)
			guessBuffer.WriteString(colorReset)
		} else {
			guessBuffer.WriteString(v)

		}

		if i == pos {
			carrotBuffer.WriteString("^ ")
		} else {
			carrotBuffer.WriteString("  ")
		}

		guessBuffer.WriteString(" ")
		i++
	}

	guessBuffer.WriteString(" ]")
	return guessBuffer.String() + "\n" + carrotBuffer.String()
}

//func newLetterList() LetterList {
//	return LetterList{
//		Letter{pos: 0, match: ""},
//		Letter{pos: 1, match: ""},
//		Letter{pos: 2, match: ""},
//		Letter{pos: 3, match: ""},
//		Letter{pos: 4, match: ""},
//	}
//}

//func updateLetterList(list LetterList, letter string, pos int, mode int) LetterList {
//	for key, _ := range list {
//		switch mode {
//		//match
//		case 1:
//			if key == pos {
//				list[key].match = letter
//			}
//		//include
//		case 2:
//			if key == pos {
//				list[key].exclude = append(list[key].exclude, letter)
//			}
//			if key != pos {
//				list[key].include = append(list[key].include, letter)
//			}
//		//exclude
//		case 3:
//			list[key].exclude = append(list[key].exclude, letter)
//		}
//	}
//
//	return LetterList{}
//}

func printWordList(p PairList) {
	i := 0
	highest := 0
	for range p {
		if p[i].Value >= highest || len(p) <= 10 {
			fmt.Printf("%s (%v)\n", p[i].Key, p[i].Value)
			highest = p[i].Value
		}
		i++
	}
}

//func printLetterList(ll LetterList) {
//	i := 0
//	var match string
//	//maxIncludes := 0
//	//maxExcludes := 0
//	var out bytes.Buffer
//	out.WriteString("[ ")
//	for range ll {
//		fmt.Printf("Pos: %v Match: %s (%v)\n", ll[i].pos, ll[i].match)
//		//fmt.Println("includes: ")
//		//fmt.Println(ll[i].include)
//		//fmt.Println("excludes: ")
//		//fmt.Println(ll[i].exclude)
//		match = ll[i].match
//		if match != "" {
//			out.WriteString(match + " ")
//		} else {
//			out.WriteString("_ ")
//		}
//
//		i++
//	}
//	out.WriteString("]")
//	println(out.String())
//}

func main() {
	info := countDictionary([]string{})
	var output string
	//ll := newLetterList()
	text := info.text
	answerSet := []int{}

	//reader := bufio.NewReader(os.Stdin)
	word := make(map[int]string)
	word[0] = ""
	word[1] = ""
	word[2] = ""
	word[3] = ""
	word[4] = ""
	colorReset := "\033[0m"
	//
	//colorRed := "\033[31m"
	colorGreen := "\033[32m"
	//colorYellow := "\033[33m"
	colorBlue := "\033[34m"
	//colorPurple := "\033[35m"
	//colorCyan := "\033[36m"
	//colorWhite := "\033[37m"

	for {

		printWordList(info.wordList)
		fmt.Println("total words: ", len(info.wordList))
		fmt.Println("recommended guess:", info.wordList[0].Key, " (", info.wordList[0].Value, ")")
		fmt.Println("---------------------")

		fmt.Println(colorBlue, "Make a guess: ", colorReset)
		fmt.Print("-> ")
		var guess string
		fmt.Scan(&guess)
		// convert CRLF to LF
		guess = strings.ToLower(strings.Replace(guess, "\n", "", -1))
		// Make a Regex to say we only want letters and numbers
		reg, err := regexp.Compile("[^a-z]+")
		if err != nil {
			log.Fatal(err)
		}
		guess = reg.ReplaceAllString(guess, "")

		guessLetters := stringToSlice(guess)
		fmt.Println(guess)
		i := 0
		for _, v := range guessLetters {
			//fmt.Println(guessLetters)
			//fmt.Println(word)
			output = printGuess(guessLetters, word, i)
			output += "\n 1 - match, 2 - include, 3 - exclude"
			output += "\n-> "

			if word[i] == v {
				output += "\n" + colorGreen + "we should be skipping this Letter" + colorReset
				fmt.Print(output)
			} else {
				var input int
				fmt.Println(output)
				fmt.Scan(&input)

				switch input {
				case 1:
					text = findMatch(text, v, i)
					word[i] = v
				case 2:
					text = findInclude(text, v, i)
				case 3:
					text = findExclude(text, v)
				default:
					output += "\nyou suck, that wasn't a 1, 2 or a 3...start over"
				}

				//updateLetterList(ll, v, i, input)
			}

			fmt.Println("words remaining:", len(text))
			output = ""
			i++
		}

		//fmt.Println(text)
		//fmt.Println(ll)
		answerSet = append(answerSet, len(text))
		fmt.Println(answerSet)
		//printLetterList(ll)

		info = countDictionary(text)
		text = info.text

		if len(text) <= 1 {
			fmt.Println("recommended guess:", info.wordList[0].Key, " (", info.wordList[0].Value, ")")
			os.Exit(0)
		}
		guess = ""
	}
}
