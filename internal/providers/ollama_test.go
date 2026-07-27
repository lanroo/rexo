package providers

import (
	"reflect"
	"testing"
)

func TestParseOllamaList(t *testing.T) {
	out := `NAME                ID              SIZE      MODIFIED
gpt-oss:20b         17052f91a42e    13 GB     2 days ago
qwen3:14b           bdbd181c33f2    9.3 GB    2 days ago
deepseek-r1:1.5b    e0979632db5a    1.1 GB    3 days ago
`
	got := parseOllamaList(out)
	want := []string{"gpt-oss:20b", "qwen3:14b", "deepseek-r1:1.5b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseOllamaList = %v, want %v", got, want)
	}
}

func TestParseOllamaListEmpty(t *testing.T) {
	if got := parseOllamaList("NAME  ID  SIZE  MODIFIED\n"); len(got) != 0 {
		t.Errorf("expected no models, got %v", got)
	}
}

func TestExplicitModelWins(t *testing.T) {
	o := &Ollama{cli: cliRunner{bin: "ollama"}, model: "qwen3:14b"}
	got, err := o.resolveModel(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "qwen3:14b" {
		t.Errorf("explicit model ignored: got %q", got)
	}
}
