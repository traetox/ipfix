package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type nfv9Record struct {
	Name      string
	ElementId string
	Length    string
}

func createV9Registry(recs []nfv9Record) {
	dictFile, err := os.Create("builtin-v9-dictionary.go")
	if err != nil {
		log.Fatalln(err)
	}
	defer dictFile.Close()

	dictFile.WriteString("package ipfix\n\n")
	dictFile.WriteString("// Autogenerated " + time.Now().Format(time.UnixDate) + "\n")
	dictFile.WriteString("var builtinNetflowV9Dictionary = fieldDictionary{\n")
	for _, r := range recs {

		if r.Name == "" || r.Length == "" || r.ElementId == "" {
			continue
		}

		dictFile.WriteString(fmt.Sprintf("\tdictionaryKey{0, %3s}: DictionaryEntry{FieldID: %3s, Name: \"%s\", Type: FieldTypes[\"%s\"]},\n",
			strings.TrimSpace(r.ElementId),
			strings.TrimSpace(r.ElementId),
			strings.TrimSpace(r.Name),
			strings.TrimSpace(r.Length)))
	}
	dictFile.WriteString("}\n")
}

func decodeV9Records() (nfrecs []nfv9Record) {
	f, err := os.Open("nfv9-fields.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	rdr := csv.NewReader(f)

	records, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range records {
		if len(v) != 3 {
			log.Fatalf("read bad column: %v", v)
		}
		rec := nfv9Record{
			Name:      v[0],
			ElementId: v[1],
			Length:    v[2],
		}
		nfrecs = append(nfrecs, rec)
	}
	return
}

func main() {
	recs := decodeV9Records()
	createV9Registry(recs)
}