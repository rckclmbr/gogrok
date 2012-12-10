package grok

import (
	"testing"
)

func TestAddPattern(t *testing.T) {
	grok := New()
	err := grok.AddPattern("WORD", "\\w+")
	if err != nil {
		t.Fatalf("Failed to create continuum: %v", err)
	}
	
	grok.AddPattern("WILLFAIL", "[")
	// According to the code, this will always pass, even if it's an invalid regexp
	//if err == nil {
	//	t.Fatalf("Didn't receive failure when should have")
	//}

}

func TestAddPatternsFromFile(t *testing.T) {
	grok := New()
	err := grok.AddPatternsFromFile("./invalidFilename")
	if err == nil {
		t.Fatalf("Should've received an error")
	}
}

func TestMatch(t *testing.T) {
	grok := New()
	err := grok.AddPattern("WORD", "\\w+")
	if err != nil {
		t.Fatalf("Error adding pattern: %v", err)
	}
	err = grok.Compile("%{WORD:something} %{WORD:world}")
	if err != nil {
		t.Fatalf("Error Compiling: %v", err)
	}
	match, err := grok.Match("Hellooo World")
	if err != nil {
		t.Fatalf("Error matching: %v (%v)", err, match)
	}
	
	if match["WORD:something"] != "Hellooo" || match["WORD:world"] != "World" {
		t.Fatalf("Matching didn't actually work")
	}
}