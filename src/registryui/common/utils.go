package common

import (
	"fmt"
)

// PrettySize format bytes in more readable units.
func PrettySize(size float64) string {
	units := []string{"B", "KB", "MB", "GB"}
	i := 0
	for size > 1024 && i < len(units) {
		size = size / 1024
		i = i + 1
	}
	// Format decimals as follow: 0 B, 0 KB, 0.0 MB, 0.00 GB
	decimals := i - 1
	if decimals < 0 {
		decimals = 0
	}
	return fmt.Sprintf("%.*f %s", decimals, size, units[i])
}
