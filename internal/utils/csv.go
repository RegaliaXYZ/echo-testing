package utils

import (
	"bp-echo-test/internal/models"
	"encoding/csv"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
)

func CorpusToCSV(task string, payload []models.Payload, corpus_type string, w *storage.Writer) error {
	// write csv header
	csvWriter := csv.NewWriter(w)
	if task == "nlu" {
		header := []string{"utterance", "intent"}
		writeErr := csvWriter.Write(header)
		if writeErr != nil {
			return writeErr
		}

		for _, row := range payload {
			if strings.Contains(row.Type, strings.ToLower(corpus_type)) {
				record := []string{row.Utterance, row.Intent}
				writeErr := csvWriter.Write(record)
				if writeErr != nil {
					return writeErr
				}
			}

		}

	} else if task == "ner" {
		//TODO: test ner upload content with QC integration
		header := []string{"tags", "utterance", "span_offsets"}
		writeErr := csvWriter.Write(header)
		if writeErr != nil {
			return writeErr
		}
		for _, row := range payload {
			if strings.Contains(row.Type, strings.ToLower(corpus_type)) {
				var tags []string
				var offsets []string
				for index, tag := range row.Tags {
					tags = append(tags, tag.Name)
					offsets = append(offsets, strconv.Itoa(row.Tags[index].StartOffset)+"-"+strconv.Itoa(row.Tags[index].EndOffset))
				}
				formattedTags := strings.Join(tags, " ")
				formattedOffsets := strings.Join(offsets, " ")
				record := []string{formattedTags, row.Utterance, formattedOffsets}
				writeErr := csvWriter.Write(record)
				if writeErr != nil {
					return writeErr
				}
			}
		}
	}
	csvWriter.Flush()
	return nil
}
