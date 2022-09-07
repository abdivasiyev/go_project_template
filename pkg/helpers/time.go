package helpers

import "time"

func TimeToString(t time.Time, layout string, valid bool) string {
	if !valid {
		return ""
	}

	return t.Format(layout)
}
