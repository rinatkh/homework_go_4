package structs

import "unsafe"

type User struct {
	ID         int
	Name       string
	Active     bool
	LoginCount int
	Tags       []string
}

func UpdateUserName(users map[int]User, id int, name string) bool {
	// TODO: изменить имя существующего пользователя, хранящегося в map как значение.
	// Если пользователя нет, вернуть false и не добавлять новый ключ.
	return false
}

func ActivateUser(users map[int]User, id int) bool {
	// TODO: сделать существующего пользователя активным.
	// Если пользователя нет, вернуть false.
	return false
}

func IncrementLoginCount(users map[int]User, id int) bool {
	// TODO: увеличить LoginCount существующего пользователя на единицу.
	// Если пользователя нет, вернуть false.
	return false
}

func CopyUsers(users map[int]User) map[int]User {
	// TODO: вернуть независимую копию map и вложенных слайсов Tags.
	// Изменения любой части результата не должны менять исходные данные.
	return nil
}

func BuildPointerIndex(users []User) map[int]*User {
	// TODO: построить индекс ID -> *User из независимых копий входных пользователей.
	// Указатели разных ID не должны вести к одному объекту или менять исходный слайс.
	return nil
}

func RenameThroughPointer(users map[int]*User, id int, name string) bool {
	// TODO: изменить имя объекта по существующему ненулевому указателю.
	// Для отсутствующего ключа или nil-значения вернуть false.
	return false
}

func (u User) DisplayName() string {
	// TODO: вернуть строку "ID:Name".
	return ""
}

type Admin struct {
	User
	Permissions []string
}

func (a Admin) DisplayName() string {
	// TODO: вернуть строку "admin:ID:Name", перекрывая одноимённый метод встроенного User.
	return ""
}

func (a Admin) HasPermission(permission string) bool {
	// TODO: проверить точное наличие permission в списке администратора.
	// Регистр, пробелы и повторения не нормализуются.
	return false
}

func RenameEmbeddedUser(admin *Admin, name string) {
	// TODO: изменить Name у встроенного User.
	// Для nil admin функция должна безопасно завершиться.
}

type Account struct {
	ID      int
	Balance int
}

type PremiumAccount struct {
	Account
	Bonus int
}

func (p PremiumAccount) Total() int {
	// TODO: вернуть сумму основного баланса и бонуса.
	return 0
}

func PremiumLabel(p PremiumAccount) string {
	// TODO: вернуть строку "ID:Balance total=Total".
	return ""
}

type BadLayout struct {
	Active   bool
	Amount   int64
	Verified bool
	Retries  int32
}

type GoodLayout struct {
	Amount   int64
	Retries  int32
	Active   bool
	Verified bool
}

func UserValueSize(user User) uintptr {
	// TODO: вернуть фактический размер переданного значения User в байтах.
	return 0
}

func UserPointerSize(user *User) uintptr {
	// TODO: вернуть фактический размер переменной-указателя на User в байтах.
	// Результат не зависит от того, nil указатель или нет.
	return 0
}

func LayoutSizes() (bad, good, saved uintptr) {
	// TODO: вернуть размеры BadLayout и GoodLayout, затем число сэкономленных байт.
	return 0, 0, 0
}

var _ uintptr = unsafe.Sizeof(User{})
