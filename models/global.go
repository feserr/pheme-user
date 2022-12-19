package models

import (
	"github.com/feserr/pheme-user/database"
)

// Db global db var
var Db = database.Connect()
