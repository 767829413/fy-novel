package chapter

import (
	"fmt"
	"testing"
)

func TestTxtConvert(t *testing.T) {
	content := ``
	res := txtConvert("Title", content)
	fmt.Println(res)
}
