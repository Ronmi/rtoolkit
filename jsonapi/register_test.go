package jsonapi

import "testing"

func TestCamelConverting(t *testing.T) {
	// [method name, expect]
	cases := [][2]string{
		{"CamelCase", "camel_case"},
		{"URL123", "url123"},
		{"TestIdGetter", "test_id_getter"},
		{"TestIDGetter", "test_id_getter"},
		{"TestID3Getter", "test_id3_getter"},
	}

	for _, c := range cases {
		t.Run(c[0], func(t *testing.T) {
			if actual := convertCamelTo_(c[0]); c[1] != actual {
				t.Fatalf("expected %s, got %s", c[1], actual)
			}
		})
	}
}
