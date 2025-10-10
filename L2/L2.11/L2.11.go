package main

import (
	"fmt"
	"sort"
	"strings"
)

// SetAnagram группирует слова по множествам анаграмм
func SetAnagram(words []string) map[string][]string {
	groups := make(map[string][]string)
	seen := make(map[string]bool) // для удаления дубликатов

	//проходимся по словам
	for _, word := range words {
		//приводим слово к нижнему регистру
		w := strings.ToLower(word)
		// если слово уже есть в мапе, тогда переходим к следующему слову
		if seen[w] {
			continue
		}
		//добавляем в мапу для проверки дубликатов
		seen[w] = true

		// Разделяем слово на буквы для сортировки этих букв
		letters := strings.Split(w, "")
		//сортируем буквы
		sort.Strings(letters)
		//объединяем отсортированные буквы и приравниаем к ключу
		key := strings.Join(letters, "")

		//добавляем слова в мапу по ключу(сортированные буквы), добавляются слова
		groups[key] = append(groups[key], w)
	}

	// Формируем финальную map
	result := make(map[string][]string)
	//проходимся по мапе
	for _, group := range groups {
		//если длина слайса значений больше одного(т.к. если меньше означает, что слово не имеет анаграммы)
		//тогда сортируем слайс слов и берем первое из слов как ключ к результируюзей мапе
		if len(group) > 1 {
			//сортируем слайс слов
			sort.Strings(group)
			// ключ — первое слово из исходного списка (до сортировки)
			result[group[0]] = group
		}
	}
	//возвращаем результирующую мапу
	return result
}

func main() {
	//создаем слайс слов
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	//результат равен возвращаемой мапе из функции
	res := SetAnagram(words)

	//проходимся по мапе и выводим результаты
	for key, group := range res {
		fmt.Printf("%s: %v\n", key, group)
	}
}
