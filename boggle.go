package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

type coordinate struct {
	x int
	y int
}

type p struct {
	Results template.HTML
}

const (
	//staticDir is the place for static html
	staticDir = "./web/"
	//PORT set to regular http 8080
	port = "8080"
)

var dictionary = make(map[string]bool)
var boggleWords = make(map[string]bool)

func getAllwords(boggleBoard [4][4]string, visited map[coordinate]bool, row int, col int, word *string) {
	bCoordinate := coordinate{row, col}
	visited[bCoordinate] = true
	*word = *word + boggleBoard[row][col]
	checkDictionary(*word)
	for boggleBoardRow := row - 1; boggleBoardRow <= row+1 && boggleBoardRow < 4; boggleBoardRow++ {
		for boggleBoardCol := col - 1; boggleBoardCol <= col+1 && boggleBoardCol < 4; boggleBoardCol++ {
			boggleWordCoordinate := coordinate{boggleBoardRow, boggleBoardCol}
			if boggleBoardRow >= 0 && boggleBoardCol >= 0 && !visited[boggleWordCoordinate] {
				getAllwords(boggleBoard, visited, boggleBoardRow, boggleBoardCol, word)
			}
		}
	}
	temp := *word
	*word = temp[0 : len(temp)-1]
	visited[bCoordinate] = false
}

//permutation combination of all possibilities in the board from left to right downwards
func bogglePnC(boggleBoard [4][4]string) {
	visited := make(map[coordinate]bool)
	var word string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			getAllwords(boggleBoard, visited, i, j, &word)
		}
	}

}

func createDictionary() {
	file, err := os.Open(staticDir + "words_alpha.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := (scanner.Text())
		if !dictionary[word] {
			dictionary[word] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkDictionary(word string) {
	if dictionary[word] && len(word) > 3 && !boggleWords[word] {
		boggleWords[word] = true
	}
}

//GetWordsHandler is a handler for results
func GetWordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//clear boggleWords for next board.
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var boggleBoard [4][4]string
		var entry = 1
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				b := r.FormValue(strconv.Itoa(entry))
				boggleBoard[i][j] = b
				entry++
			}
		}
		bogglePnC(boggleBoard)
		var words []string
		for w := range boggleWords {
			words = append(words, w)
		}
		sort.Strings(words)
		var results string
		for _, word := range words {
			results = results + "<li class=\"list-group-item\">" + word + "</li>"
		}
		response := p{template.HTML(results)}
		t, err := template.ParseFiles(staticDir + "results.html")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("clear the set of words")
		boggleWords = make(map[string]bool)
		if err := t.ExecuteTemplate(w, "results.html", response); err != nil {
			fmt.Fprintf(w, "ExecuteTemplate() err: %v", err)
		}
	}
}

func main() {
	createDictionary()
	fmt.Println("Starting web server...")
	r := mux.NewRouter()
	r.HandleFunc("/getwords", GetWordsHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))
	http.ListenAndServe(":"+port, r)
}
