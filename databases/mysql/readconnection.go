package mysql

func ReadConnection(host, user, pass, dbname string, readData bool) *Database {
	database := new(Database)

	conn := NewConnection(host, user, pass, dbname)
	database.ReadConnection(&conn, readData)

	return database
}
