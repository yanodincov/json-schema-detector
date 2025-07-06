package analyze

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/analyzer"
)

var (
	outputFile string
	autoCommit bool
)

// Cmd представляет команду analyze
var Cmd = &cobra.Command{
	Use:   "analyze [input.json]",
	Short: "Анализирует JSON файл и создает схему",
	Long: `Анализирует структуру JSON файла и генерирует соответствующую 
JSON Schema с автоматическим определением типов и структур.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	Cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Выходной файл для схемы")
	Cmd.Flags().BoolVarP(&autoCommit, "auto-commit", "a", false, "Автоматический коммит изменений схемы")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Проверяем существование входного файла
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("входной файл не найден: %s", inputFile)
	}

	// Если выходной файл не указан, создаем его на основе входного
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		outputFile = inputFile[:len(inputFile)-len(ext)] + ".schema.json"
	}

	fmt.Printf("Анализ файла: %s\n", inputFile)
	fmt.Printf("Выходной файл: %s\n", outputFile)

	// Создаем анализатор
	analyzer := analyzer.New()

	// Анализируем файл
	result, err := analyzer.AnalyzeFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка анализа: %w", err)
	}

	// Сохраняем результат
	if err := analyzer.SaveSchema(result, outputFile); err != nil {
		return fmt.Errorf("ошибка сохранения схемы: %w", err)
	}

	fmt.Printf("Схема успешно создана: %s\n", outputFile)
	fmt.Printf("Проанализировано объектов: %d\n", result.Statistics.TotalObjects)
	fmt.Printf("Уникальных структур: %d\n", result.Statistics.UniqueStructures)

	// Автоматический коммит если флаг установлен
	if autoCommit {
		if err := commitSchemaChanges(outputFile, "analyze"); err != nil {
			fmt.Printf("⚠️ Ошибка автоматического коммита: %v\n", err)
		} else {
			fmt.Printf("✅ Изменения схемы закоммичены\n")
		}
	}

	return nil
}

// commitSchemaChanges выполняет автоматический коммит изменений схемы
func commitSchemaChanges(schemaFile, operation string) error {
	// Проверяем, что мы в git репозитории
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git не найден")
	}

	// Добавляем файл схемы в git
	cmd := exec.Command("git", "add", schemaFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка git add: %w", err)
	}

	// Создаем коммит
	commitMessage := fmt.Sprintf("schema: %s %s", operation, filepath.Base(schemaFile))
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка git commit: %w", err)
	}

	return nil
}
