package validator

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/qnhqn1/file-validator/internal/cache"
	"github.com/qnhqn1/file-validator/internal/storage/pgstorage"
)

type Service interface {
	ValidateAndStore(ctx context.Context, key string, payload []byte) error
}

type service struct {
	storage pgstorage.StorageInterface
	cache   cache.Cache
}

func New(storage pgstorage.StorageInterface, cache cache.Cache) Service {
	return &service{storage: storage, cache: cache}
}

func (s *service) ValidateAndStore(ctx context.Context, key string, payload []byte) error {

	if err := validateDOCX(payload); err != nil {
		return fmt.Errorf("Валидация DOCX не удалась: %w", err)
	}

	_ = s.cache.Set(ctx, "validated:"+key, []byte("1"), 5*time.Minute)

	if err := s.storage.InsertEvent(ctx, key, payload); err != nil {
		return fmt.Errorf("сохранить событие: %w", err)
	}
	return nil
}

func validateDOCX(data []byte) error {

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("не является допустимым ZIP: %w", err)
	}

	requiredFiles := map[string]bool{
		"[Content_Types].xml": false,
		"_rels/.rels":         false,
		"word/document.xml":   false,
	}

	for _, file := range reader.File {
		if _, ok := requiredFiles[file.Name]; ok {
			requiredFiles[file.Name] = true
		}

		if strings.HasPrefix(file.Name, "word/") {

			if strings.Contains(file.Name, "..") {
				return fmt.Errorf("подозрительный путь в ZIP: %s", file.Name)
			}
		}
	}

	for name, present := range requiredFiles {
		if !present {
			return fmt.Errorf("отсутствует обязательный файл: %s", name)
		}
	}

	docFile, err := reader.Open("word/document.xml")
	if err != nil {
		return fmt.Errorf("невозможно открыть document.xml: %w", err)
	}
	defer docFile.Close()

	xmlContent, err := ioutil.ReadAll(docFile)
	if err != nil {
		return fmt.Errorf("невозможно прочитать document.xml: %w", err)
	}
	content := string(xmlContent)
	if !strings.HasPrefix(strings.TrimSpace(content), "<?xml") && !strings.Contains(content, "<w:document") {
		return fmt.Errorf("document.xml не выглядит как допустимый XML")
	}

	text := extractTextFromDOCX(content)

	if err := validateCyrillicPercentage(text); err != nil {
		return fmt.Errorf("Валидация кириллицы не удалась: %w", err)
	}

	if err := validateDates(text); err != nil {
		return fmt.Errorf("валидация даты не удалась: %w", err)
	}

	return nil
}

func extractTextFromDOCX(xmlContent string) string {

	re := regexp.MustCompile(`<w:t[^>]*>(.*?)</w:t>`)
	matches := re.FindAllStringSubmatch(xmlContent, -1)
	var text strings.Builder
	for _, match := range matches {
		if len(match) > 1 {
			text.WriteString(match[1])
		}
	}
	return text.String()
}

func validateCyrillicPercentage(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("в документе не найден текст")
	}
	cyrillicCount := 0
	totalChars := 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			totalChars++
			if unicode.Is(unicode.Cyrillic, r) {
				cyrillicCount++
			}
		}
	}
	if totalChars == 0 {
		return fmt.Errorf("в документе не найдены буквы")
	}
	percentage := float64(cyrillicCount) / float64(totalChars) * 100
	if percentage < 90 {
		return fmt.Errorf("документ содержит только %.2f%% кириллицы, требуется 90%%", percentage)
	}
	return nil
}

func validateDates(text string) error {

	datePatterns := []*regexp.Regexp{
		regexp.MustCompile(`\b\d{1,2}\.\d{1,2}\.\d{4}\b`), // ДД.ММ.ГГГГ
		regexp.MustCompile(`\b\d{4}-\d{1,2}-\d{1,2}\b`),   // ГГГГ-ММ-ДД
		regexp.MustCompile(`\b\d{1,2}/\d{1,2}/\d{4}\b`),   // ДД/ММ/ГГГГ или ММ/ДД/ГГГГ
	}

	var validDates []time.Time
	for _, pattern := range datePatterns {
		matches := pattern.FindAllString(text, -1)
		for _, match := range matches {
			if parsed, ok := parseDate(match); ok {
				validDates = append(validDates, parsed)
			}
		}
	}

	if len(validDates) == 0 {
		return fmt.Errorf("в документе не найдены допустимые даты")
	}

	minDate := validDates[0]
	maxDate := validDates[0]
	for _, d := range validDates {
		if d.Before(minDate) {
			minDate = d
		}
		if d.After(maxDate) {
			maxDate = d
		}
	}

	diff := maxDate.Sub(minDate)
	threeYears := 3 * 365 * 24 * time.Hour
	if diff > threeYears {
		return fmt.Errorf("даты в документе отличаются более чем на 3 года (min: %s, max: %s)", minDate.Format("02.01.2006"), maxDate.Format("02.01.2006"))
	}

	return nil
}

func parseDate(dateStr string) (time.Time, bool) {
	formats := []string{
		"02.01.2006", // ДД.MM.ГГГГ
		"2006-01-02", // ГГГГ-ММ-ДД
		"01/02/2006", // ММ/ДД/ГГГГ
		"02/01/2006", // ДД/MM/ГГГГ
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
