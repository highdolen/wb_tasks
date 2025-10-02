package main

import (
	"fmt"
	"sort"
	"strings"
)

// Тип для сортировки рун
type runeSlice []rune

func (r runeSlice) Len() int           { return len(r) }
func (r runeSlice) Less(i, j int) bool { return r[i] < r[j] }
func (r runeSlice) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// Функция для сортировки букв в слове
func sortRunes(s string) string {
	runes := runeSlice([]rune(s))
	sort.Sort(runes)
	return string(runes)
}

// Функция для поиска множеств анаграмм
func FindAnagramSets(words []string) map[string][]string {
	anagramGroups := make(map[string][]string)

	// Приводим все слова к нижнему регистру и группируем по "сигнатуре" (отсортированные буквы)
	for _, word := range words {
		word = strings.ToLower(word)
		sig := sortRunes(word)
		anagramGroups[sig] = append(anagramGroups[sig], word)
	}

	result := make(map[string][]string)

	for _, group := range anagramGroups {
		if len(group) > 1 { // Игнорируем одиночные слова
			// Убираем дубликаты
			unique := make(map[string]struct{})
			for _, w := range group {
				unique[w] = struct{}{}
			}

			cleanGroup := make([]string, 0, len(unique))
			for w := range unique {
				cleanGroup = append(cleanGroup, w)
			}

			// Сортируем по алфавиту
			sort.Strings(cleanGroup)

			// Ключом делаем первое слово в отсортированном списке
			result[cleanGroup[0]] = cleanGroup
		}
	}

	return result
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	anagramSets := FindAnagramSets(words)

	for k, v := range anagramSets {
		fmt.Printf("%s: %v\n", k, v)
	}
}
