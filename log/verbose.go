package log

func Verbose(message string) {
	if VERBOSE <= currentLevel() {
		log("VERBOSE", message)
	}
}
