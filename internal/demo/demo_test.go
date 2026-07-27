package demo

import (
	"context"
	"testing"
)

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"REST APIs":            "rest-apis",
		"  Machine Learning  ": "machine-learning",
		"C++ & Go!":            "c-go",
		"":                     "mini-aula",
		"---":                  "mini-aula",
	}
	for in, want := range cases {
		if got := slug(in); got != want {
			t.Errorf("slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLanguageName(t *testing.T) {
	cases := map[string]string{
		"":       "English", // default
		"en":     "English",
		"pt":     "Portuguese",
		"PT":     "Portuguese",
		"es":     "Spanish",
		"fr":     "French",
		"German": "German", // pass-through for unknown codes
	}
	for in, want := range cases {
		if got := languageName(in); got != want {
			t.Errorf("languageName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRunRejectsEmptyTopic(t *testing.T) {
	_, err := Run(context.Background(), Options{Topic: "   ", ProjectDir: t.TempDir()})
	if err == nil {
		t.Fatal("expected error for empty topic")
	}
}
