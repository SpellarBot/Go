package utils

import (
	"strings"
)

func GetLocation(location, country string) string {
	loc := strings.ToUpper(location)
	cou := strings.ToUpper(country)

	if loc != "" {
		return getLocationId(loc)
	} else if cou != "" {
		return getLocationId(cou)
	} else {
		return getLocationId("ALL")
	}
}

func getLocationId(code string) string {
	loc := "US"
	if code != "" {
		switch strings.TrimSpace(code) {
		case "SPECIAL":
			loc = "SPECIAL"
		case "NP", "IN", "PK", "BD":
			loc = "IN"
		case "ID":
			loc = "ID"
		case "BR":
			loc = "BR"
		case "MENA", "AE", "SA", "EG", "IQ", "DZ", "OM", "QA":
			loc = "AE"
		case "US":
		case "ALL":
			loc = "US"
		}
	}

	return loc
}
