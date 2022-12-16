package models

import "github.com/feserr/pheme-user/database"

var Db = database.Connect()

func removeFromList[T comparable](list []T, item T) []T {
	for pos, other := range list {
		if other == item {
			return append(list[:pos], list[pos+1:]...)
		}
	}
	return list
}
