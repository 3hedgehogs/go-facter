package formatter

import (
	"fmt"
	y "github.com/ghodss/yaml"
)

// YAMLFormatter prints-out facts in YAML format
type YAMLFormatter struct {
}

// NewYAMLFormatter returns new YAML formatter
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{}
}

// Print prints-out facts in YAML format
func (yf *YAMLFormatter) Print(facts map[string]interface{}) error {
	d, err := y.Marshal(&facts)
	if err != nil {
		return err
	}
	fmt.Printf("---\n%s\n", string(d))
	return nil
}
