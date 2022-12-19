package models

import (
	"os"

	"github.com/feserr/pheme-user/database"
)

// Db global db var
var Db = database.Connect(os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"))

func removeFromList[T comparable](list []T, item T) []T {
	for pos, other := range list {
		if other == item {
			return append(list[:pos], list[pos+1:]...)
		}
	}
	return list
}
