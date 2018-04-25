package service

import "strings"

func validateCorsOrigin(origin, cors string) (allow bool) {
	if cors == "*" {
		return true
	}

	allowedHosts := strings.Split(cors, ",")

	for _, host := range allowedHosts {
		if origin == host {
			return true
		}
	}

	return false
}
