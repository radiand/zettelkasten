package common

import "time"
import "os"

// Now returns current time in UTC.
func Now() time.Time {
	return time.Now().UTC()
}

// ModificationTime return last modification time of a file.
func ModificationTime(path string) (time.Time, error) {
	fstat, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fstat.ModTime(), nil
}

// Delta returns duration in seconds.
func Delta(before, after time.Time) int64 {
	return int64(after.Sub(before).Abs().Seconds())
}
