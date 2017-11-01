package utils

func ABTagAdd(oldTag string, newTag string) string {
	if newTag == "" {
		return oldTag
	}

	if oldTag == "" {
		return newTag
	} else {
		return oldTag + "|" + newTag
	}
}
