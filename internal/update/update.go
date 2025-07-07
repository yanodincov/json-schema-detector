package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-schema-detector/pkg/analyzer"
)

var (
	inputFile  string
	autoCommit bool
)

// Cmd представляет команду update
var Cmd = &cobra.Command{
	Use:   "update [schema.json]",
	Short: "Обновляет существующую схему новыми данными",
	Long: `Обновляет существующую JSON Schema новыми данными из JSON файла,
сохраняя существующие описания и комментарии.`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	Cmd.Flags().StringVarP(&inputFile, "input", "i", "", "JSON файл с новыми данными")
	Cmd.Flags().BoolVarP(&autoCommit, "auto-commit", "a", false, "Автоматический коммит изменений схемы")
	Cmd.MarkFlagRequired("input")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	schemaFile := args[0]

	// Проверяем существование файлов
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("файл схемы не найден: %s", schemaFile)
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("входной файл не найден: %s", inputFile)
	}

	fmt.Printf("Обновление схемы: %s\n", schemaFile)
	fmt.Printf("Новые данные: %s\n", inputFile)

	// Создаем анализатор
	analyzer := analyzer.New()

	// Загружаем существующую схему
	existingSchema, err := analyzer.LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки схемы: %w", err)
	}

	// Анализируем новые данные
	newResult, err := analyzer.AnalyzeFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка анализа новых данных: %w", err)
	}

	// Объединяем схемы
	mergedResult, err := analyzer.MergeResults(existingSchema, newResult)
	if err != nil {
		return fmt.Errorf("ошибка объединения схем: %w", err)
	}

	// Сохраняем обновленную схему
	if err := analyzer.SaveSchema(mergedResult, schemaFile); err != nil {
		return fmt.Errorf("ошибка сохранения схемы: %w", err)
	}

	fmt.Printf("Схема успешно обновлена: %s\n", schemaFile)
	fmt.Printf("Добавлено новых объектов: %d\n", newResult.Statistics.TotalObjects)

	// Автоматический коммит если флаг установлен
	if autoCommit {
		if err := commitSchemaChanges(schemaFile, "update"); err != nil {
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
