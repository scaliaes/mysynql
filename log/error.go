package log

func Error(message string) {
	if ERROR <= currentLevel() {
		log("ERROR", message)
	}
}
