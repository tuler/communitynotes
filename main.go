package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
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

    reader := bufio.NewReader(notesFile)
    header := readHeader(reader)

    // Write sampled notes
    outFile, err := os.Create(outputPath)
    check(err)
    defer outFile.Close()
    writer := bufio.NewWriter(outFile)
    defer writer.Flush()

    // Write header
    fmt.Fprintln(writer, header)

    // Read all noteIds first
    var noteIds []string
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        fields := strings.Split(strings.TrimSpace(line), "\t")
        noteIds = append(noteIds, fields[0])
    }

    // Randomly sample noteIds
    selectedNotes := make(map[string]bool)
    for i := 0; i < sampleSize && i < len(noteIds); i++ {
        idx := rand.Intn(len(noteIds))
        selectedNotes[noteIds[idx]] = true
        noteIds = append(noteIds[:idx], noteIds[idx+1:]...)
    }

    // Reopen file to read full records
    notesFile.Seek(0, 0)
    reader = bufio.NewReader(notesFile)
    readHeader(reader) // Skip header

    // Write selected records
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
    flag.IntVar(&sampleSize, "sample", 1000, "number of notes to sample")
    flag.Parse()
    
    // create output directories
    os.MkdirAll("input-sample/ratings", 0755)
    
    // Track participating users
    participants := make(map[string]bool)

    selectedNotes := sampleNotes("input/notes-00000.tsv", "input-sample/notes-00000.tsv", func(fields []string) bool {
        participants[fields[1]] = true
        return true
    })
    
    // Filter ratings
    // First declare the function variable
    var filterRatings = func(fields []string) bool {
        if selectedNotes[fields[0]] {
            participants[fields[1]] = true
            return true
        }
        return false
    }
    
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
        return participants[fields[0]]
    })
}

// go run main.go -sample 10000
