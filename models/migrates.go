package models

func MainMigrations() {
	mainDB := MainDbBegin()
	defer mainDB.DbRollback()

	// config
	mainDB.AutoMigrate(&Config{})

	// device
	mainDB.AutoMigrate(&Device{})

	// identity
	mainDB.AutoMigrate(&Identity{})
	mainDB.Model(&Identity{}).AddUniqueIndex("index_identity_on_source_and_symbol", "source", "symbol")

	// public_key
	mainDB.AutoMigrate(&PublicKey{})
	mainDB.Model(&PublicKey{}).AddUniqueIndex("index_public_key_on_user_sn_and_index", "user_sn", "index")

	// token
	mainDB.AutoMigrate(&Token{})

	// user
	mainDB.AutoMigrate(&User{})
}

func LogMigrations() {
	logDB := LogDbBegin()
	defer logDB.DbRollback()
}
