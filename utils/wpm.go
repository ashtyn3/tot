package utils

func WPM(characters int, time int) int {
	mins := float64(time) / 60
	wpm := int((float64(characters) / 5.0) / mins)
	return wpm
}
