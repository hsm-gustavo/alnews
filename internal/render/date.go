package render

import "time"

func FormatDate(date string) string {
	t, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		// fail-safe, just return default if there's an error
		return date
	}

	formattedDate := t.Format("Jan 02, 2006 @ 15:04")
	
	return formattedDate
}