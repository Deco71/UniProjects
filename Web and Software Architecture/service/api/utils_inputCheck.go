package api

import (
	"errors"
	"net/http"
	"strconv"
)

func CheckStringLenght(line string) error {
	if len(line) < 3 || len(line) > 16 {
		return errors.New("string too long")
	}
	return nil
}

func getOffset(r *http.Request) int {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	return offset
}
