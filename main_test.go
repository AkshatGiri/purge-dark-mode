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
		`class="dark:"`, // technically invalid class, but still should be handled properly
		`dabc`, // less than 5 characters, starts with d
		``, // empty string
	}

	expected := []string{
		`class="normal-class other-class"`,  // the middle dark class should be removed, only leaving 1 space
		`class="normal-class  other-class"`, //  middle class should be removed, 2 spaces remaing in the middle
		`class="normal-class"`,              //  removed dark class at the end
		`class="other-class"`,               //  removed dark class in the beginning
		`class="other-class"`,               //  removed dark class and the new line after it
		`class=""`,                          // removed the only dark class
		`class=""`, // removed the invalid class.
		`dabc`, // should remain the same and not crash the program
		``, // should remain empty string and not crash
	}

	for i, input := range inputs {
		output, _ := RemoveDarkClasses(input)

		if output != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], output)
		}
	}
}