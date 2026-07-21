package maps

type User struct {
	ID     int
	Name   string
	City   string
	Active bool
	Tags   []string
}

type Product struct {
	SKU      string
	Quantity int
	Price    int
}

func CountWords(words []string) map[string]int {
	// TODO: посчитать, сколько раз встречается каждое слово.
	// Вернуть новую готовую к записи map. Регистр и пустая строка имеют значение.
	return nil
}

func GetOrDefault(values map[string]int, key string, fallback int) int {
	// TODO: вернуть значение существующего ключа или fallback, если ключ отсутствует.
	// Нулевое значение по существующему ключу должно возвращаться как обычное значение.
	return 0
}

func DeleteAndReport(values map[string]int, key string) bool {
	// TODO: удалить ключ и сообщить, существовал ли он до вызова.
	// Для отсутствующего ключа, пустой и nil map вернуть false без panic.
	return false
}

func MergeCounters(left, right map[string]int) map[string]int {
	// TODO: объединить два счётчика в новой map, складывая значения общих ключей.
	// Исходные map после вызова должны остаться без изменений.
	return nil
}

func KeysByValueAtLeast(values map[string]int, min int) []string {
	// TODO: вернуть ключи, значения которых не меньше min.
	// Результат должен иметь стабильный алфавитный порядок.
	return nil
}

func MaxKeyByValue(values map[string]int) (string, bool) {
	// TODO: найти ключ с максимальным значением.
	// Для пустого входа вернуть "", false; при равенстве выбрать меньший ключ по алфавиту.
	return "", false
}

func BuildUserIndex(users []User) map[int]User {
	// TODO: построить индекс пользователей по ID.
	// Если один ID встречается несколько раз, в результате должен остаться последний пользователь.
	return nil
}

func GroupActiveUsersByCity(users []User) map[string][]User {
	// TODO: сгруппировать только активных пользователей по городу.
	// Порядок пользователей внутри каждого города должен совпадать с исходным слайсом.
	return nil
}

func CountTags(users []User) map[string]int {
	// TODO: посчитать все появления тегов у всех пользователей.
	// Повторяющийся тег у одного пользователя также считается отдельным появлением.
	return nil
}

func BuildInventory(products []Product) map[string]Product {
	// TODO: построить склад по SKU.
	// При повторном SKU оставить последнюю запись из входного слайса.
	return nil
}

func ReserveStock(inventory map[string]Product, sku string, count int) bool {
	// TODO: уменьшить остаток существующего товара на положительное count, если товара хватает.
	// При неуспехе вернуть false и не менять склад.
	return false
}

func Restock(inventory map[string]*Product, sku string, count int) bool {
	// TODO: увеличить остаток товара, хранящегося в map как указатель.
	// Ключ должен существовать, указатель не должен быть nil, count должен быть положительным.
	return false
}

func LowStockSKUs(inventory map[string]Product, limit int) []string {
	// TODO: вернуть SKU товаров с остатком не больше limit.
	// Результат должен иметь стабильный алфавитный порядок.
	return nil
}

func InventoryValue(inventory map[string]Product) int {
	// TODO: посчитать общую стоимость склада как сумму Quantity * Price.
	// Пустой и nil склад имеют стоимость 0.
	return 0
}

func ApplyPriceUpdates(inventory map[string]Product, updates map[string]int) int {
	// TODO: применить положительные новые цены только к существующим товарам.
	// Вернуть число реально обновлённых товаров; неизвестные SKU и неположительные цены пропустить.
	return 0
}
