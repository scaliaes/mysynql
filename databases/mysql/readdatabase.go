package mysql

func ReadDatabase(host, user, pass, dbname string) *Database {
	database := new(Database)

	conn := NewConnection(host, user, pass, dbname)
	database.ReadConnection(&conn)

	return database
}
