package kernel

import (
	"fmt"
	"sort"
	"strings"
)

// Provider is a deterministic capability implementation: the same inputs always
// produce the same output bytes. No network, no clock, no randomness.
type Provider func(inputs map[string]any) (string, error)

// providers is the built-in deterministic registry for the walking skeleton.
var providers = map[string]Provider{
	"text.constant":  textConstant,
	"text.uppercase": textUppercase,
	"text.concat":    textConcat,
	"text.template":  textTemplate,
}

// Capabilities returns the registered capability names (for diagnostics).
func Capabilities() []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}

func textConstant(inputs map[string]any) (string, error) {
	v, ok := inputs["value"].(string)
	if !ok {
		return "", fmt.Errorf("text.constant requires a string input %q", "value")
	}
	return v, nil
}

func textUppercase(inputs map[string]any) (string, error) {
	v, ok := inputs["text"].(string)
	if !ok {
		return "", fmt.Errorf("text.uppercase requires a string input %q", "text")
	}
	return strings.ToUpper(v), nil
}

func textConcat(inputs map[string]any) (string, error) {
	raw, ok := inputs["parts"].([]any)
	if !ok {
		return "", fmt.Errorf("text.concat requires an array input %q", "parts")
	}
	parts := make([]string, 0, len(raw))
	for i, p := range raw {
		s, ok := p.(string)
		if !ok {
			return "", fmt.Errorf("text.concat part %d is not a string", i)
		}
		parts = append(parts, s)
	}
	sep, _ := inputs["separator"].(string)
	return strings.Join(parts, sep), nil
}

func textTemplate(inputs map[string]any) (string, error) {
	tmpl, ok := inputs["template"].(string)
	if !ok {
		return "", fmt.Errorf("text.template requires a string input %q", "template")
	}
	result := tmpl
	if vars, ok := inputs["vars"].(map[string]any); ok {
		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys) 
		for _, k := range keys {
			result = strings.ReplaceAll(result, "{{"+k+"}}", fmt.Sprintf("%v", vars[k]))
		}
	}
	return result, nil
}
