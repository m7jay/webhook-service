package workers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"

	"github.com/m7jay/webhook-service/internal/services"
)

type FileProcessor struct {
	service *services.WebhookService
}

func NewFileProcessor(service *services.WebhookService) *FileProcessor {
	return &FileProcessor{
		service: service,
	}
}

func (fp *FileProcessor) ProcessFile(filename string, eventID uint) error {
	ext := filepath.Ext(filename)
	switch ext {
	case ".csv":
		return fp.processCSV(filename, eventID)
	case ".xlsx":
		return fp.processExcel(filename, eventID)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

func (fp *FileProcessor) processCSV(filename string, eventID uint) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue
		}

		payload := make(map[string]string)
		for i, value := range record {
			payload[headers[i]] = value
		}

		if err := fp.service.TriggerWebhook(eventID, payload); err != nil {
			log.Printf("Error triggering webhook for CSV record: %v", err)
		}
	}

	return nil
}

func (fp *FileProcessor) processExcel(filename string, eventID uint) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return fmt.Errorf("Excel file must have at least two rows (headers and data)")
	}

	headers := rows[0]

	for _, row := range rows[1:] {
		payload := make(map[string]string)
		for i, value := range row {
			if i < len(headers) {
				payload[headers[i]] = value
			}
		}

		if err := fp.service.TriggerWebhook(eventID, payload); err != nil {
			log.Printf("Error triggering webhook for Excel row: %v", err)
		}
	}

	return nil
}
