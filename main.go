package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var sampleSize int

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readHeader(reader *bufio.Reader) string {
	header, err := reader.ReadString('\n')
	check(err)
	return strings.TrimSpace(header)
}

func sampleNotes(inputPath, outputPath string, checkFunc func([]string) bool) map[string]bool {
	// Open notes file
	notesFile, err := os.Open(inputPath)
	check(err)
	defer notesFile.Close()

	// First, count ratings per note
	ratingCounts := make(map[string]int)

	// Count ratings across all rating files
	for i := 0; i < 16; i++ {
		ratingsPath := fmt.Sprintf("input/ratings/ratings-%05d.tsv", i)
		ratingsFile, err := os.Open(ratingsPath)
		check(err)

		reader := bufio.NewReader(ratingsFile)
		readHeader(reader) // Skip header

		// Count ratings for each note
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fields := strings.Split(strings.TrimSpace(line), "\t")
			noteId := fields[0]
			ratingCounts[noteId]++
		}

		ratingsFile.Close()
	}

	// Create a slice of note IDs sorted by rating count
	type noteCount struct {
		noteId string
		count  int
	}
	var sortedNotes []noteCount
	for noteId, count := range ratingCounts {
		sortedNotes = append(sortedNotes, noteCount{noteId, count})
	}

	// Sort by count in descending order
	sort.Slice(sortedNotes, func(i, j int) bool {
		return sortedNotes[i].count > sortedNotes[j].count
	})

	// Select top N notes
	selectedNotes := make(map[string]bool)
	for i := 0; i < sampleSize && i < len(sortedNotes); i++ {
		selectedNotes[sortedNotes[i].noteId] = true
	}

	// Write selected notes to output file
	reader := bufio.NewReader(notesFile)
	header := readHeader(reader)

	outFile, err := os.Create(outputPath)
	check(err)
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Write header
	fmt.Fprintln(writer, header)

	// Read all noteIds first
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fields := strings.Split(strings.TrimSpace(line), "\t")
		if selectedNotes[fields[0]] {
			checkFunc(fields)
			fmt.Fprint(writer, line)
		}
	}

	return selectedNotes
}

func filterFile(inputPath, outputPath string, checkFunc func([]string) bool) {
	inFile, err := os.Open(inputPath)
	check(err)
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	check(err)
	defer outFile.Close()

	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Copy header
	header := readHeader(reader)
	fmt.Fprintln(writer, header)

	// Filter records
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fields := strings.Split(strings.TrimSpace(line), "\t")
		if checkFunc(fields) {
			fmt.Fprint(writer, line)
		}
	}
}

func main() {
	// Add flag parsing at start of main
	var userSampleSize int
	flag.IntVar(&sampleSize, "sample", 1000, "number of notes to sample")
	flag.IntVar(&userSampleSize, "users", 1000, "number of top users to include")
	flag.Parse()

	// create output directories
	os.MkdirAll("input-sample/ratings", 0755)

	// Track participating users
	participants := make(map[string]bool)

	selectedNotes := sampleNotes("input/notes-00000.tsv", "input-sample/notes-00000.tsv", func(fields []string) bool {
		participants[fields[1]] = true
		return true
	})

	// Count ratings per user for selected notes
	userRatingCounts := make(map[string]int)

	// First pass: count ratings per user
	for i := 0; i < 16; i++ {
		inputPath := fmt.Sprintf("input/ratings/ratings-%05d.tsv", i)
		file, err := os.Open(inputPath)
		check(err)

		reader := bufio.NewReader(file)
		readHeader(reader)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fields := strings.Split(strings.TrimSpace(line), "\t")
			if selectedNotes[fields[0]] {
				userId := fields[1]
				userRatingCounts[userId]++
			}
		}
		file.Close()
	}

	// Select top N users
	type userCount struct {
		userId string
		count  int
	}
	var sortedUsers []userCount
	for userId, count := range userRatingCounts {
		sortedUsers = append(sortedUsers, userCount{userId, count})
	}

	sort.Slice(sortedUsers, func(i, j int) bool {
		return sortedUsers[i].count > sortedUsers[j].count
	})

	// Create map of selected users
	selectedUsers := make(map[string]bool)
	for i := 0; i < userSampleSize && i < len(sortedUsers); i++ {
		selectedUsers[sortedUsers[i].userId] = true
	}

	// Filter ratings
	// First declare the function variable
	var filterRatings = func(fields []string) bool {
		if selectedNotes[fields[0]] && selectedUsers[fields[1]] {
			participants[fields[1]] = true
			return true
		}
		return false
	}

	// Second pass: write filtered ratings
	for i := 0; i < 16; i++ {
		inputPath := fmt.Sprintf("input/ratings/ratings-%05d.tsv", i)
		outputPath := fmt.Sprintf("input-sample/ratings/ratings-%05d.tsv", i)
		filterFile(inputPath, outputPath, filterRatings)
	}

	// Filter noteStatusHistory
	filterFile("input/noteStatusHistory-00000.tsv", "input-sample/noteStatusHistory-00000.tsv", func(fields []string) bool {
		return selectedNotes[fields[0]]
	})

	// Filter userEnrollment (only keep participants from ratings)
	filterFile("input/userEnrollment-00000.tsv", "input-sample/userEnrollment-00000.tsv", func(fields []string) bool {
		return selectedUsers[fields[0]]
	})
}

// go run main.go -sample 10000 -users 1000
