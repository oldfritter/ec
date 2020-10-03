package models

func MainMigrations() {
	mainDB := MainDbBegin()
	defer mainDB.DbRollback()

	// config
	mainDB.AutoMigrate(&Config{})

	// device
	mainDB.AutoMigrate(&Device{})

	// friend_ship
	mainDB.AutoMigrate(&FriendShip{})
	mainDB.Model(&FriendShip{}).AddUniqueIndex("index_friend_ships_on_owner_id_and_friend_id", "owner_id", "friend_id")

	// group
	mainDB.AutoMigrate(&Group{})

	// group_member
	mainDB.AutoMigrate(&GroupMember{})
	mainDB.Model(&GroupMember{}).AddUniqueIndex("index_group_members_on_member_id_and_group_id", "member_id", "group_id")

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
