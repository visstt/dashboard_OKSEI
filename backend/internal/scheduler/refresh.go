package scheduler

import (
	"fmt"
	"log"
	"os"
	"time"

	"dashboard/internal/converter"
)

type Scheduler struct {
	projectRoot      string
	attendanceInput  string
	attendanceOutput string
	statementInput   string
	statementOutput  string
	pythonScript     string
	// Кэш времени последнего изменения файлов для оптимизации
	lastModified map[string]time.Time
}

func NewScheduler(projectRoot, attendanceInput, attendanceOutput, statementInput, statementOutput, pythonScript string) *Scheduler {
	return &Scheduler{
		projectRoot:      projectRoot,
		attendanceInput:  attendanceInput,
		attendanceOutput: attendanceOutput,
		statementInput:   statementInput,
		statementOutput:  statementOutput,
		pythonScript:     pythonScript,
		lastModified:     make(map[string]time.Time),
	}
}

// RefreshData обновляет данные, запуская оба конвертера
// Проверяет изменения файлов перед конвертацией (оптимизация)
func (s *Scheduler) RefreshData() error {
	log.Println("[Scheduler] Начало обновления данных...")

	// Проверяем наличие входных файлов и их изменения
	if shouldUpdate, err := s.shouldUpdateFile(s.attendanceInput, s.attendanceOutput); err != nil {
		log.Printf("[Scheduler] Предупреждение: %v", err)
	} else if shouldUpdate {
		// Конвертируем посещаемость
		log.Println("[Scheduler] Конвертация посещаемости...")
		if err := converter.ConvertAttendance(s.attendanceInput, s.attendanceOutput); err != nil {
			return fmt.Errorf("ошибка конвертации посещаемости: %v", err)
		}
		// Обновляем время последнего изменения
		if info, err := os.Stat(s.attendanceInput); err == nil {
			s.lastModified[s.attendanceInput] = info.ModTime()
		}
		log.Println("[Scheduler] Посещаемость обновлена")
	} else {
		log.Println("[Scheduler] Посещаемость не изменилась, пропускаем")
	}

	// Проверяем наличие файла ведомости и его изменения
	if shouldUpdate, err := s.shouldUpdateFile(s.statementInput, s.statementOutput); err != nil {
		log.Printf("[Scheduler] Предупреждение: %v", err)
	} else if shouldUpdate {
		// Конвертируем ведомость
		log.Println("[Scheduler] Конвертация ведомости...")
		if err := converter.ConvertStatement(s.statementInput, s.statementOutput, s.pythonScript); err != nil {
			return fmt.Errorf("ошибка конвертации ведомости: %v", err)
		}
		// Обновляем время последнего изменения
		if info, err := os.Stat(s.statementInput); err == nil {
			s.lastModified[s.statementInput] = info.ModTime()
		}
		log.Println("[Scheduler] Ведомость обновлена")
	} else {
		log.Println("[Scheduler] Ведомость не изменилась, пропускаем")
	}

	log.Println("[Scheduler] Обновление данных завершено успешно!")
	return nil
}

// shouldUpdateFile проверяет, нужно ли обновлять файл
// Возвращает true, если входной файл новее выходного или выходного файла нет
func (s *Scheduler) shouldUpdateFile(inputFile, outputFile string) (bool, error) {
	// Проверяем наличие входного файла
	inputInfo, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("входной файл не найден: %s", inputFile)
	}
	if err != nil {
		return false, fmt.Errorf("ошибка проверки входного файла: %v", err)
	}

	// Проверяем наличие выходного файла
	outputInfo, err := os.Stat(outputFile)
	if os.IsNotExist(err) {
		// Выходного файла нет - нужно обновить
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("ошибка проверки выходного файла: %v", err)
	}

	// Сравниваем время изменения
	// Если входной файл новее выходного - нужно обновить
	if inputInfo.ModTime().After(outputInfo.ModTime()) {
		return true, nil
	}

	// Проверяем кэш (если файл уже обрабатывался в этой сессии)
	if lastMod, exists := s.lastModified[inputFile]; exists {
		if inputInfo.ModTime().After(lastMod) {
			return true, nil
		}
	}

	return false, nil
}
