package mysql

func Apply(database *Database, host, user, pass, dbname string) {
	conn := NewConnection(host, user, pass, dbname)

	database.Apply(&conn)
}
