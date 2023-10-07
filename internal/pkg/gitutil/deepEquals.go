package gitutil

import (
	"fmt"
	"reflect"
)

func DeepEqualCompare(arrSet1 [][]interface{}, arrSet2 [][]interface{}) {
	for i := 0; i < len(arrSet1); i++ {
		difference := CompareDeepEquals(arrSet1[i], arrSet2[i])
		if difference != "" {
			fmt.Printf("Difference in array %d: %s\n", i, difference)
		}
	}
}

func CompareDeepEquals(a, b interface{}) string {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return fmt.Sprintf("Different types: %T vs %T", a, b)
	}

	if !reflect.DeepEqual(a, b) {
		return fmt.Sprintf("Different values: %v vs %v", a, b)
	}

	return ""
}
