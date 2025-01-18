package main

import (
	"testing"
)

func TestDarkClassRemoval(t *testing.T) {

	inputs := []string{
		`class="normal-class dark:bg-gray-800 other-class"`,   // Single space before/after
		`class="normal-class  dark:bg-gray-800  other-class"`, //  Multiple spaces before/after
		`class="normal-class dark:bg-gray-800"`,               //  Space before, nothing after
		`class="dark:bg-gray-800 other-class"`,                //  Nothing before, space after
		`class="dark:bg-gray-800
		other-class"`, //  new line after class
		`class="dark:bg-gray-800"`, // only dark class
	}

	expected := []string{
		`class="normal-class other-class"`,  // Single space before/after
		`class="normal-class  other-class"`, //  Multiple spaces before/after
		`class="normal-class"`,              //  Space before, nothing after
		`class="other-class"`,               //  Nothing before, space after
		`class="other-class"`,               //  new line after class
		`class=""`,
	}

	for i, input := range inputs {
		output := RemoveDarkClass(input)
		if output != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], output)
		}
	}
}
