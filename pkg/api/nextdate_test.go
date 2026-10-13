package api

import (
	"testing"
	"time"
)

func TestNextDateSimple(t *testing.T) {

	now, _ := time.Parse(dateFormat, "20240210")

	result, err := NextDate(now, "20240210", "m 1,-1 2,8")

	if err != nil {

		t.Fatal(err)

	}

	t.Log(result)

}
