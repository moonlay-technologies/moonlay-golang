package helper

import "fmt"

func GenerateClause(id int, status string) string {
	return fmt.Sprintf("id = %d AND status = '%s'", id, status)
}
