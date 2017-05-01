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
	if(err != nil){
		return nil, err
	}
	err = cfg.readFromFile(configFolder + "/" + environment + "." + configFile, true)
	if(err != nil){
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

func (c *Config) GetString(key string, defaultVal string) string {
	val, _ := c.configData[key]

	if len(val) == 0 {
		val = defaultVal
	}

	return val
}

func (c *Config) GetInt(key string, defaultVal int64) int64 {
	val, _ := c.configData[key]

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return intVal
}

func (c *Config) GetFloat(key string, defaultVal float64) float64 {
	val, _ := c.configData[key]

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return floatVal
}

func (c *Config) GetUInt(key string, defaultVal uint64) uint64 {
	val, _ := c.configData[key]

	uiVal, err := strconv.ParseUint(val, 10,64)
	if err != nil {
		return defaultVal
	}
	return uiVal
}

func (c *Config) GetBool(key string, defaultVal bool) bool {
	val, _ := c.configData[key]

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return boolVal
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
