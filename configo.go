package configo


import (
	"os"
	"bufio"
	"strconv"
	"strings"
	"errors"
	"regexp"
)

const (
	separator = "="
	configFolder = "./config"
	configFile = "app.config"
)

var re *regexp.Regexp


type Config struct {
	environment string
	configData map[string]string
	configKeys []string
}

//noinspection GoUnusedExportedFunction
func New(environment string) (*Config, error) {
	re = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9\\._-]*$")

	cfg := Config{environment:environment, configData:make(map[string]string) }

	err := cfg.readFromFile(configFolder + "/" + configFile, false)
	if err != nil {
		return nil, err
	}
	err = cfg.readFromFile(configFolder + "/" + environment + "." + configFile, true)
	if err != nil {
		return nil, err
	}

	for k := range cfg.configData {
		cfg.configKeys = append(cfg.configKeys, k)
	}

	return &cfg, nil
}

func (c *Config) Environment() string {
	return c.environment
}

func (c *Config) GetKeys() []string {
	return c.configKeys
}

func (c *Config) GetString(key string) (string, error) {
	val, ok := c.configData[key]

	if !ok {
		return "", errors.New("Config entry with key '" + key + "' does not exist")
	}

	return val, nil
}

func (c *Config) GetStringOrDefault(key string, defaultVal string) string {
	val, err := c.GetString(key)

	if err != nil {
		val = defaultVal
	}

	return val
}

func (c *Config) GetInt(key string) (int64, error) {
	val, ok := c.configData[key]

	if !ok {
		return 0, errors.New("Config entry with key '" + key + "' does not exist")
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.New("Config value '" + val + "' for key '" + key + "' is not integer")
	}
	return intVal, nil
}

func (c *Config) GetIntOrDefault(key string, defaultVal int64) int64 {
	val, err := c.GetInt(key)

	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetFloat(key string) (float64, error) {
	val, ok := c.configData[key]

	if !ok {
		return 0, errors.New("Config entry with key '" + key + "' does not exist")
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New("Config value '" + val + "' for key '" + key + "' is not float")
	}
	return floatVal, nil
}

func (c *Config) GetFloatOrDefault(key string, defaultVal float64) float64 {
	val, err := c.GetFloat(key)

	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetUInt(key string) (uint64, error) {
	val, ok := c.configData[key]

	if !ok {
		return 0, errors.New("Config entry with key '" + key + "' does not exist")
	}

	uiVal, err := strconv.ParseUint(val, 10,64)
	if err != nil {
		return 0, errors.New("Config value '" + val + "' for key '" + key + "' is not unsigned integer")
	}
	return uiVal, nil
}

func (c *Config) GetUIntOrDefault(key string, defaultVal uint64) uint64 {
	val, err := c.GetUInt(key)

	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetBool(key string) (bool, error) {
	val, ok := c.configData[key]

	if !ok {
		return false, errors.New("Config entry with key '" + key + "' does not exist")
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, errors.New("Config value '" + val + "' for key '" + key + "' is not boolean")
	}
	return boolVal, nil
}

func (c *Config) GetBoolOrDefault(key string, defaultVal bool) bool {
	val, err := c.GetBool(key)

	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) readFromFile(configFile string, ignoreFileError bool) error {
	file, err := os.Open(configFile)
	defer file.Close()

	if !ignoreFileError {
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		key, value, err := c.parseLine(line)
		if err != nil {
			return err
		}

		c.configData[key] = value
	}

	return nil
}

func (c *Config) parseLine(line string) (key string, value string, err error) {
	// Extract key part
	pos := strings.Index(line, separator)
	if pos == -1 {
		return "", "", errors.New("Invalid config entry. " + line)
	}

	key = strings.Trim(strings.TrimSpace(line[:pos]), "\"")
	if len(key) == 0 || !re.MatchString(key) {
		return "", "", errors.New("Invalid key. " + line)
	}

	value = ""

	// Extract value part
	adjustedPos := pos + len(separator)
	if adjustedPos < len(line) {
		value = strings.Trim(strings.TrimSpace(line[adjustedPos:]), "\"")
	}

	return key, value, nil
}
