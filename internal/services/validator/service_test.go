package validator

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/qnhqn1/file-validator/internal/services/validator/mocks"
)

func createValidDOCXPayload() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Пример текста на кириллице с датой 27.12.2025</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createInvalidCyrillicDOCXPayload() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>This is English text without Cyrillic and no date</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithoutDates() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Пример текста на кириллице без даты</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithLowCyrillic() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Text with some Cyrillic буквы but mostly English words and date 27.12.2025</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithInvalidDate() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст на кириллице с невалидной датой 32.13.2025</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createEmptyDOCX() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithNumbersOnly() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>1234567890</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXMissingRequiredFile() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{

		"_rels/.rels":       `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml": `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithSuspiciousPath() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст</w:t></w:r></w:p></w:body></w:document>`,
		"word/../evil.txt":    `malicious content`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithInvalidXML() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `This is not XML at all, just plain text`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithXMLButNoDocument() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><root>Not a document</root>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithEmptyDocumentXML() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   ``,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithMultipleDates() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст с датами 27.12.2025 и 2025-12-27</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithDatesTooFarApart() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст с датами 01.01.2020 и 01.01.2025</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

func createDOCXWithDateMMDDYYYY() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8"?>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8"?>`,
		"word/document.xml":   `<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Текст с датой 12/27/2025</w:t></w:r></w:p></w:body></w:document>`,
	}

	for name, content := range files {
		w, _ := zw.Create(name)
		w.Write([]byte(content))
	}
	zw.Close()
	return buf.Bytes()
}

type ValidatorServiceSuite struct {
	suite.Suite
	ctx     context.Context
	cache   *mocks.MockCache
	storage *mocks.MockStorageInterface
	svc     Service
}

func (s *ValidatorServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.cache = &mocks.MockCache{}
	s.storage = &mocks.MockStorageInterface{}
	s.svc = New(s.storage, s.cache)
}

func (s *ValidatorServiceSuite) TestValidateAndStore_Success() {
	key := "test-key"
	payload := createValidDOCXPayload()

	s.cache.On("Set", s.ctx, "validated:"+key, []byte("1"), 5*time.Minute).Return(nil)
	s.storage.On("InsertEvent", s.ctx, key, payload).Return(nil)

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().NoError(err)
}

func (s *ValidatorServiceSuite) TestValidateAndStore_InvalidDOCX() {
	key := "test-key"
	payload := []byte("invalid")

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "Валидация DOCX не удалась")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_InvalidCyrillic() {
	key := "test-key"
	payload := createInvalidCyrillicDOCXPayload()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "Валидация кириллицы не удалась")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_NoDates() {
	key := "test-key"
	payload := createDOCXWithoutDates()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "валидация даты не удалась")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_LowCyrillic() {
	key := "test-key"
	payload := createDOCXWithLowCyrillic()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "Валидация кириллицы не удалась")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_InvalidDate() {
	key := "test-key"
	payload := createDOCXWithInvalidDate()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "валидация даты не удалась")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_EmptyDocument() {
	key := "test-key"
	payload := createEmptyDOCX()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "в документе не найден текст")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_NumbersOnly() {
	key := "test-key"
	payload := createDOCXWithNumbersOnly()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "в документе не найдены буквы")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_MissingRequiredFile() {
	key := "test-key"
	payload := createDOCXMissingRequiredFile()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "отсутствует обязательный файл")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_SuspiciousPath() {
	key := "test-key"
	payload := createDOCXWithSuspiciousPath()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "подозрительный путь в ZIP")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_InvalidXML() {
	key := "test-key"
	payload := createDOCXWithInvalidXML()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "не выглядит как допустимый XML")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_XMLButNoDocument() {
	key := "test-key"
	payload := createDOCXWithXMLButNoDocument()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "в документе не найден текст")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_DateMMDDYYYY() {
	key := "test-key"
	payload := createDOCXWithDateMMDDYYYY()

	s.cache.On("Set", s.ctx, "validated:"+key, []byte("1"), 5*time.Minute).Return(nil)
	s.storage.On("InsertEvent", s.ctx, key, payload).Return(nil)

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().NoError(err)
}

func (s *ValidatorServiceSuite) TestValidateAndStore_EmptyDocumentXML() {
	key := "test-key"
	payload := createDOCXWithEmptyDocumentXML()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "не выглядит как допустимый XML")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_StorageError() {
	key := "test-key"
	payload := createValidDOCXPayload()

	s.cache.On("Set", s.ctx, "validated:"+key, []byte("1"), 5*time.Minute).Return(nil)
	s.storage.On("InsertEvent", s.ctx, key, payload).Return(fmt.Errorf("storage error"))

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "сохранить событие")
}

func (s *ValidatorServiceSuite) TestValidateAndStore_MultipleDates() {
	key := "test-key"
	payload := createDOCXWithMultipleDates()

	s.cache.On("Set", s.ctx, "validated:"+key, []byte("1"), 5*time.Minute).Return(nil)
	s.storage.On("InsertEvent", s.ctx, key, payload).Return(nil)

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().NoError(err)
}

func (s *ValidatorServiceSuite) TestValidateAndStore_DatesTooFarApart() {
	key := "test-key"
	payload := createDOCXWithDatesTooFarApart()

	err := s.svc.ValidateAndStore(s.ctx, key, payload)
	s.Require().Error(err)
	assert.Contains(s.T(), err.Error(), "даты в документе отличаются более чем на 3 года")
}

func TestValidatorServiceSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceSuite))
}
