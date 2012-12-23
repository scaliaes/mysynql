package log

func Log(message string) {
	if STANDARD <= currentLevel() {
		log("LOG", message)
	}
}
