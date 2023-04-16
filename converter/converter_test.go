package converter

import (
	"fmt"
	"testing"
)

func TestNewHexConverter(t *testing.T) {
	type testData struct {
		inputData    string
		expectedData string
	}

	tests := []testData{
		{
			inputData:    "123456789abcdef",
			expectedData: "81985529216486895",
		},
		{
			inputData:    "F",
			expectedData: "15",
		},
		{
			inputData:    "AaFf",
			expectedData: "43775",
		},
		{
			inputData:    "123",
			expectedData: "291",
		},
	}

	converter := NewHexConverter()
	for _, test := range tests {
		result := converter.Convert(test.inputData)

		if result != test.expectedData {
			err := fmt.Sprintf("input: %s expected: %s; got: %s",
				test.inputData, test.expectedData, result)
			t.Error(err)
		}
	}

}
