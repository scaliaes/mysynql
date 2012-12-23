package log

func Debug(message string) {
	if DEBUG <= currentLevel() {
		log("DEBUG", message)
	}
}
