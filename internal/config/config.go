package config

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type CustomBadge struct {
	Label string `yaml:"label"`
	Image string `yaml:"image"`
	URL   string `yaml:"url"`
}

type CIBadgeConfig struct {
	Enabled  bool
	Workflow string
	Branch   string
}

func (c *CIBadgeConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		var b bool
		if err := value.Decode(&b); err != nil {
			return err
		}
		c.Enabled = b
		c.Workflow = "ci.yml"
		c.Branch = "main"
		return nil
	}
	type rawCI struct {
		Workflow string `yaml:"workflow"`
		Branch   string `yaml:"branch"`
	}
	var raw rawCI
	if err := value.Decode(&raw); err != nil {
		return err
	}
	c.Enabled = true
	c.Workflow = raw.Workflow
	if c.Workflow == "" {
		c.Workflow = "ci.yml"
	}
	c.Branch = raw.Branch
	if c.Branch == "" {
		c.Branch = "main"
	}
	return nil
}

type BadgesConfig struct {
	CI      CIBadgeConfig `yaml:"ci"`
	Release bool          `yaml:"release"`
	License bool          `yaml:"license"`
	Custom  []CustomBadge `yaml:"custom"`
}

type DocEntry struct {
	Title       string `yaml:"title"`
	Path        string `yaml:"path"`
	Description string `yaml:"description"`
}

func (d *DocEntry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		d.Description = strings.TrimSpace(value.Value)
		return nil
	}
	type rawDoc struct {
		Title       string `yaml:"title"`
		Path        string `yaml:"path"`
		Description string `yaml:"description"`
	}
	var raw rawDoc
	dec := yaml.NewDecoder(strings.NewReader(nodeToYAML(value)))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	d.Title = raw.Title
	d.Path = raw.Path
	d.Description = raw.Description
	return nil
}

type DocumentationConfig struct {
	Architecture DocEntry `yaml:"architecture"`
	Development  DocEntry `yaml:"development"`
	Deployment   DocEntry `yaml:"deployment"`
}

func (doc *DocumentationConfig) UnmarshalYAML(value *yaml.Node) error {
	var m map[string]yaml.Node
	if err := value.Decode(&m); err != nil {
		return err
	}
	for k := range m {
		if k != "architecture" && k != "development" && k != "deployment" {
			return fmt.Errorf("unknown documentation entry: '%s'", k)
		}
	}
	if n, ok := m["architecture"]; ok {
		if err := n.Decode(&doc.Architecture); err != nil {
			return err
		}
	} else {
		return errors.New("missing required documentation entry: 'architecture'")
	}
	if n, ok := m["development"]; ok {
		if err := n.Decode(&doc.Development); err != nil {
			return err
		}
	} else {
		return errors.New("missing required documentation entry: 'development'")
	}
	if n, ok := m["deployment"]; ok {
		if err := n.Decode(&doc.Deployment); err != nil {
			return err
		}
	} else {
		return errors.New("missing required documentation entry: 'deployment'")
	}
	return nil
}

type LicenseConfig struct {
	Unlicensed bool
	Types      []string
	Holder     any
	Text       string
	Note       string
	Filename   string
}

func (l *LicenseConfig) UnmarshalYAML(value *yaml.Node) error {
	var m map[string]yaml.Node
	if err := value.Decode(&m); err != nil {
		return err
	}
	for k := range m {
		if k == "year" {
			return errors.New("field 'license.year' is removed; current year is always used automatically")
		}
		if k != "type" && k != "holder" && k != "text" && k != "note" && k != "filename" {
			return fmt.Errorf("unknown field in license: '%s'", k)
		}
	}
	typeNode, ok := m["type"]
	if !ok {
		return errors.New("missing required field 'license.type'")
	}
	if typeNode.Kind == yaml.ScalarNode {
		if typeNode.Value == "false" {
			l.Unlicensed = true
		} else {
			return errors.New("field 'license.type' must be false or a list containing 'MIT' and/or 'APACHE-2.0'")
		}
	} else if typeNode.Kind == yaml.SequenceNode {
		var types []string
		if err := typeNode.Decode(&types); err != nil {
			return err
		}
		if len(types) == 0 {
			return errors.New("field 'license.type' list cannot be empty")
		}
		for _, t := range types {
			norm := strings.ToUpper(strings.TrimSpace(t))
			if norm != "MIT" && norm != "APACHE-2.0" {
				return fmt.Errorf("unsupported license '%s'; currently only 'MIT' and 'APACHE-2.0' are supported", t)
			}
			l.Types = append(l.Types, norm)
		}
	} else {
		return errors.New("field 'license.type' must be false or a list containing 'MIT' and/or 'APACHE-2.0'")
	}

	if holderNode, ok := m["holder"]; ok {
		var h any
		if err := holderNode.Decode(&h); err != nil {
			return err
		}
		l.Holder = h
	}
	if textNode, ok := m["text"]; ok {
		l.Text = textNode.Value
	}
	if noteNode, ok := m["note"]; ok {
		l.Note = noteNode.Value
	}
	if fileNode, ok := m["filename"]; ok {
		l.Filename = fileNode.Value
	}
	return nil
}

type Config struct {
	Schema        string              `yaml:"$schema"`
	Name          string              `yaml:"name"`
	Repo          string              `yaml:"repo"`
	Tagline       string              `yaml:"tagline"`
	Badges        BadgesConfig        `yaml:"badges"`
	QuickStart    string              `yaml:"quick_start"`
	Highlights    []string            `yaml:"highlights"`
	Documentation DocumentationConfig `yaml:"documentation"`
	License       LicenseConfig       `yaml:"license"`
}

func nodeToYAML(node *yaml.Node) string {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	_ = enc.Encode(node)
	_ = enc.Close()
	return buf.String()
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("missing required field: 'name'")
	}
	if strings.TrimSpace(c.Repo) == "" {
		return errors.New("missing required field: 'repo'")
	}
	if strings.TrimSpace(c.Tagline) == "" {
		return errors.New("missing required field: 'tagline'")
	}
	if strings.TrimSpace(c.QuickStart) == "" {
		return errors.New("missing required field: 'quick_start'")
	}
	if len(c.Highlights) == 0 {
		return errors.New("field 'highlights' must be a non-empty list of strings")
	}
	return nil
}

func Parse(data []byte) (*Config, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("YAML validation error: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
