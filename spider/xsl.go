package spider

import (
	"github.com/tealeg/xlsx/v3"
	"sync"
)

// Структура для работы с файлом с результатом
type ResultFile struct {
	mu   sync.Mutex
	file *xlsx.Sheet
}

//Записывает строку в файл с локом через Mutex
func (rf *ResultFile) WriteLine(line []string) {
	rf.mu.Lock()
	row := rf.file.AddRow()
	for _, field := range line {
		cell := row.AddCell()
		cell.Value = field
	}
	defer rf.mu.Unlock()
}

// Возвращает новый ResultFile
func NewResultFile(file *xlsx.Sheet) ResultFile {
	return ResultFile{
		file: file,
	}
}
