package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ConfigPath stores the path to the config file for saving
var ConfigPath string

type Config struct {
	DataDir    string           `toml:"data_dir"` // session store directory, default ~/.cc-connect
	Agent      AgentConfig     `toml:"agent"`    // global agent config (new format)
	Platforms  []PlatformConfig `toml:"platforms"` // platform configs (new format)
	AllowUsers []AllowUser     `toml:"allow_users"`
	Log        LogConfig       `toml:"log"`
	Language   string          `toml:"language"` // "en" or "zh", default is "en"
	Speech     SpeechConfig    `toml:"speech"`

	// Legacy: support old projects-based config
	Projects []ProjectConfig `toml:"projects"`
}

// SpeechConfig configures speech-to-text for voice messages.
type SpeechConfig struct {
	Enabled  bool   `toml:"enabled"`
	Provider string `toml:"provider"` // "openai" | "groq"
	Language string `toml:"language"` // e.g. "zh", "en"; empty = auto-detect
	OpenAI   struct {
		APIKey  string `toml:"api_key"`
		BaseURL string `toml:"base_url"`
		Model   string `toml:"model"`
	} `toml:"openai"`
	Groq struct {
		APIKey string `toml:"api_key"`
		Model  string `toml:"model"`
	} `toml:"groq"`
}

// ProjectConfig binds one agent (with a specific work_dir) to one or more platforms.
type ProjectConfig struct {
	Name       string           `toml:"name"`
	Agent      AgentConfig      `toml:"agent"`
	Platforms []PlatformConfig `toml:"platforms"`
	AllowUsers []AllowUser     `toml:"allow_users"`
}

// AllowUser defines a user that is allowed to access the bot.
type AllowUser struct {
	Platform string `toml:"platform"` // "feishu", "telegram", "qq", etc.
	UserID   string `toml:"user_id"`  // platform-specific user ID
}

type AgentConfig struct {
	Type      string           `toml:"type"`
	Options   map[string]any   `toml:"options"`
	Providers []ProviderConfig `toml:"providers"`
}

type ProviderConfig struct {
	Name    string            `toml:"name"`
	APIKey  string            `toml:"api_key"`
	BaseURL string            `toml:"base_url,omitempty"`
	Model   string            `toml:"model,omitempty"`
	Env     map[string]string `toml:"env,omitempty"`
}

type PlatformConfig struct {
	Type    string         `toml:"type"`
	Options map[string]any `toml:"options"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{
		Log: LogConfig{Level: "info"},
	}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.DataDir == "" {
		if home, err := os.UserHomeDir(); err == nil {
			cfg.DataDir = filepath.Join(home, ".cc-connect")
		} else {
			cfg.DataDir = ".cc-connect"
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	// Support new single-agent config
	if c.Agent.Type != "" && len(c.Platforms) > 0 {
		for j, p := range c.Platforms {
			if p.Type == "" {
				return fmt.Errorf("config: platforms[%d].type is required", j)
			}
		}
		return nil
	}

	// Legacy: support old projects-based config
	if len(c.Projects) > 0 {
		for i, proj := range c.Projects {
			prefix := fmt.Sprintf("projects[%d]", i)
			if proj.Name == "" {
				return fmt.Errorf("config: %s.name is required", prefix)
			}
			if proj.Agent.Type == "" {
				return fmt.Errorf("config: %s.agent.type is required", prefix)
			}
			if len(proj.Platforms) == 0 {
				return fmt.Errorf("config: %s needs at least one [[projects.platforms]]", prefix)
			}
			for j, p := range proj.Platforms {
				if p.Type == "" {
					return fmt.Errorf("config: %s.platforms[%d].type is required", prefix, j)
				}
			}
		}
		return nil
	}

	// Neither new nor legacy config provided
	return fmt.Errorf("config: either [agent] + [platforms] or [[projects]] is required")
}

// SaveActiveProvider persists the active provider name.
func SaveActiveProvider(providerName string) error {
	if ConfigPath == "" {
		return fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Support new global config
	if cfg.Agent.Type != "" {
		if cfg.Agent.Options == nil {
			cfg.Agent.Options = make(map[string]any)
		}
		cfg.Agent.Options["provider"] = providerName
		return saveConfig(cfg)
	}

	// Legacy: support old projects config
	for i := range cfg.Projects {
		if cfg.Projects[i].Name == "default" {
			if cfg.Projects[i].Agent.Options == nil {
				cfg.Projects[i].Agent.Options = make(map[string]any)
			}
			cfg.Projects[i].Agent.Options["provider"] = providerName
			break
		}
	}
	return saveConfig(cfg)
}

// AddProviderToConfig adds a provider to the agent config and saves.
func AddProviderToConfig(provider ProviderConfig) error {
	if ConfigPath == "" {
		return fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Support new global config
	if cfg.Agent.Type != "" {
		for _, existing := range cfg.Agent.Providers {
			if existing.Name == provider.Name {
				return fmt.Errorf("provider %q already exists", provider.Name)
			}
		}
		cfg.Agent.Providers = append(cfg.Agent.Providers, provider)
		return saveConfig(cfg)
	}

	// Legacy: support old projects config
	found := false
	for i := range cfg.Projects {
		if cfg.Projects[i].Name == "default" {
			for _, existing := range cfg.Projects[i].Agent.Providers {
				if existing.Name == provider.Name {
					return fmt.Errorf("provider %q already exists", provider.Name)
				}
			}
			cfg.Projects[i].Agent.Providers = append(cfg.Projects[i].Agent.Providers, provider)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("no valid config found")
	}
	return saveConfig(cfg)
}

// RemoveProviderFromConfig removes a provider from the agent config and saves.
func RemoveProviderFromConfig(providerName string) error {
	if ConfigPath == "" {
		return fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Support new global config
	if cfg.Agent.Type != "" {
		providers := cfg.Agent.Providers
		for j := range providers {
			if providers[j].Name == providerName {
				cfg.Agent.Providers = append(providers[:j], providers[j+1:]...)
				return saveConfig(cfg)
			}
		}
		return fmt.Errorf("provider %q not found", providerName)
	}

	// Legacy: support old projects config
	found := false
	for i := range cfg.Projects {
		if cfg.Projects[i].Name == "default" {
			providers := cfg.Projects[i].Agent.Providers
			for j := range providers {
				if providers[j].Name == providerName {
					cfg.Projects[i].Agent.Providers = append(providers[:j], providers[j+1:]...)
					found = true
					break
				}
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("provider %q not found", providerName)
	}
	return saveConfig(cfg)
}

func saveConfig(cfg *Config) error {
	f, err := os.Create(ConfigPath)
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// SaveLanguage saves the language setting to the config file.
func SaveLanguage(lang string) error {
	if ConfigPath == "" {
		return fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	cfg.Language = lang
	return saveConfig(cfg)
}

// ListProjects returns project names from the config file.
// Deprecated: use GetAgentProviders instead.
func ListProjects() ([]string, error) {
	if ConfigPath == "" {
		return nil, fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// New global config
	if cfg.Agent.Type != "" {
		return []string{"default"}, nil
	}

	// Legacy projects config
	var names []string
	for _, p := range cfg.Projects {
		names = append(names, p.Name)
	}
	return names, nil
}

// GetAgentProviders returns providers for the agent.
func GetAgentProviders() ([]ProviderConfig, string, error) {
	if ConfigPath == "" {
		return nil, "", fmt.Errorf("config path not set")
	}
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, "", fmt.Errorf("read config: %w", err)
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, "", fmt.Errorf("parse config: %w", err)
	}

	// New global config
	if cfg.Agent.Type != "" {
		active, _ := cfg.Agent.Options["provider"].(string)
		return cfg.Agent.Providers, active, nil
	}

	// Legacy projects config - use first project
	if len(cfg.Projects) > 0 {
		p := cfg.Projects[0]
		active, _ := p.Agent.Options["provider"].(string)
		return p.Agent.Providers, active, nil
	}

	return nil, "", fmt.Errorf("no agent config found")
}

// GetProjectProviders returns providers for a given project.
// Deprecated: use GetAgentProviders instead.
func GetProjectProviders(projectName string) ([]ProviderConfig, string, error) {
	return GetAgentProviders()
}
