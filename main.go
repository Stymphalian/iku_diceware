package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	kNumRandomBytes = 120000
)

var (
	FnumPhrases = flag.Int("num-phrases", 1, "Number of phrases to generate")
	Fn          = flag.Int("n", 6, "number words to generate")
	Fwordlist   = flag.String("word_list", "wordlist.txt", "Filepath to the wordlist")
	Fjustrandom = flag.Bool("just_random", false,
		"whether just to output random digits. Generates --n number of digits")
	FnumSpecials = flag.Int("num-specials", 0,
		"How many special characters to insert to the passphrase")
	FnumUppers = flag.Int("num-uppers", 0,
		"How many many words to promote to upper case")
)

func makeMapFromWordlist(filepath string) map[int]string {
	m := make(map[int]string)

	fin, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanWords)

	state := 0
	lastNum := 0
	for scanner.Scan() {
		if state == 0 {
			lastNum, err = strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			state = 1
		} else {
			m[lastNum] = scanner.Text()
			state = 0
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

func makeSpecialCharMap() map[int]string {
	m := make(map[int]string)
	chars := `~!#$%^&*()-=+[]\{}:;"'<>?/0123456789`
	for i := 0; i < len(chars); i++ {
		col := i%6 + 1
		row := (i / 6) + 1
		m[row*10+col] = string(chars[i])
	}
	return m
}

func RandRange(lowerInc int, upperExc int) int {
	p, err := rand.Int(rand.Reader, big.NewInt(int64(upperExc-lowerInc)))
	if err != nil {
		log.Fatal(err)
	}
	p.Add(p, big.NewInt(int64(lowerInc)))
	return int(p.Int64())
}

func GetNDigitRandNum(n int) int {
	num := 0
	for j := 0; j < n; j++ {
		p := RandRange(1, 7)
		num = 10*num + p
	}
	return num
}

func main() {
	flag.Parse()

	if *Fjustrandom {
		for i := 0; i < *Fn; i++ {
			p := RandRange(1, 7)
			fmt.Print(p)
		}
		fmt.Println()
		return
	}

	// Read the wordlist into internal map
	wordlist := makeMapFromWordlist(*Fwordlist)
	speciallist := makeSpecialCharMap()

	for phraseIndex := 0; phraseIndex < *FnumPhrases; phraseIndex++ {
		passphrase := make([]string, 0)
		for i := 0; i < *Fn; i++ {
			// get 5 digit number
			num := GetNDigitRandNum(5)

			// lookup 5 digit number in wordlist
			passphrase = append(passphrase, wordlist[num])
		}

		// how many of the words to title case. no guaranteee that the same word
		// might be picked again
		for i := 0; i < *FnumUppers; i++ {
			wordIndex := RandRange(0, len(passphrase))
			passphrase[wordIndex] = strings.Title(passphrase[wordIndex])
		}

		// add special characters if required
		for i := 0; i < *FnumSpecials; i++ {
			wordIndex := RandRange(0, len(passphrase))
			charIndex := RandRange(0, len(passphrase[wordIndex]))
			specialChar := speciallist[GetNDigitRandNum(2)]

			passphrase[wordIndex] = (passphrase[wordIndex][0:charIndex] +
				specialChar +
				passphrase[wordIndex][charIndex:len(passphrase[wordIndex])])
		}

		// Print the passphrase
		for _, s := range passphrase {
			fmt.Print(s, " ")
		}
		fmt.Println()
	}
}
