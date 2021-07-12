package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// package const
const (
	servicePortENV = "SERVICE_PORT_ENV"
	versionAPIENV  = "VERSION_API_ENV"

	dbNameENV       = "DB_NAME_ENV"
	dbHostENV       = "DB_HOST_ENV"
	dbPortENV       = "DB_PORT_ENV"
	dbUserNmaeENV   = "DB_USER_NAME_ENV"
	dbPasswordENV   = "DB_PASSWORD_ENV"
	dbURLPatternENV = "DB_URL_PATTERN_ENV"
	dbSSLModeENV    = "DB_SSLMODE_ENV"

	loggerPathENV      = "LOGGER_PATH_ENV"
	loggerVerboseENV   = "LOGGER_VERBOSE_ENV"
	loggerSystemLogENV = "LOGGER_SYSTEM_LOG_ENV"

	notifierCreateConsumersENV = "NOTIFIER_CONSUMERS_CREATE_ENV"
	notifierUpdateConsumersENV = "NOTIFIER_CONSUMERS_UPDATE_ENV"
	notifierDeleteConsumersENV = "NOTIFIER_CONSUMERS_DELETE_ENV"
	notifierTimeOutENV         = "NOTIFIER_TIMEOUT_ENV"
	notifierClientMaxRetryENV  = "NOTIFIER_CLIENT_MAX_RETRY_ENV"
	notifierTimeoutIncreaceENV = "NOTIFIER_TIMEOUT_INCREACE_ENV"

	queueSizeENV      = "QUEUE_SIZE_ENV"
	goRoutinesSizeENV = "GO_ROUTINE_SIZE_ENV"
)

var (
	queueSizeDefault      = 100
	goRoutinesSizeDefault = 100
)

// package errors
var (
	errEmptyConfiguration = errors.New("empty configuration")
)

// Service is a struct with service configuration
type Service struct {
	Port       string
	VersionAPI string
}

// DB is a struct with database cofiguration
type DB struct {
	Name       string
	Host       string
	Port       string
	UserName   string
	Password   string
	PatternURL string
	SSLMode    string
}

// Logger is a struct with logger configuration
type Logger struct {
	FilePath  string
	Verbose   bool
	SystemLog bool
}

// Notifier is a struct with notifier configuration
type Notifier struct {
	Consumers             Consumers
	Timeout               int
	ClientMaxRetry        int
	ClientTimeoutIncrease int
}

// Queue is a queue config struct
type Queue struct {
	QueueSize      int
	GoRoutinesSize int
}

// OnCreate returnes a list of consumers to notify on Create action
func (n Notifier) OnCreate() []string {
	return n.Consumers.OnCreate
}

// OnUpdate returnes a list of consumers to notify on Update action
func (n Notifier) OnUpdate() []string {
	return n.Consumers.OnUpdate
}

// OnDelete returnes a list of consumers to notify on Delete action
func (n Notifier) OnDelete() []string {
	return n.Consumers.OnDelete
}

type Consumers struct {
	OnCreate []string
	OnUpdate []string
	OnDelete []string
}

// Config is a struct with concurent safe public method to access a config
type Config struct {
	mu       *sync.RWMutex
	service  Service
	db       DB
	logger   Logger
	notifier Notifier
	queue    Queue
}

// New initiates a new Configuration instance
func New() (*Config, error) {
	cfg := &Config{
		mu: &sync.RWMutex{},
	}

	err := cfg.setService()
	if err != nil {
		return nil, fmt.Errorf("failed to create config, error %s", err.Error())
	}

	err = cfg.setDB()
	if err != nil {
		return nil, fmt.Errorf("failed to create config, error %s", err.Error())
	}

	err = cfg.setLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create config, error %s", err.Error())
	}

	err = cfg.setNotifier()
	if err != nil {
		return nil, fmt.Errorf("failed to create config, error %s", err.Error())
	}

	cfg.setQueue()

	return cfg, nil
}

// Service returns a copy of Service config
func (c *Config) Service() Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Service{
		Port:       c.service.Port,
		VersionAPI: c.service.VersionAPI,
	}
}

// DB returns a copy of DB config
func (c *Config) DB() DB {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return DB{
		Name:       c.db.Name,
		Host:       c.db.Host,
		Port:       c.db.Port,
		UserName:   c.db.UserName,
		Password:   c.db.Password,
		PatternURL: c.db.PatternURL,
		SSLMode:    c.db.SSLMode,
	}
}

// Logger returns a copy of Logger config
func (c *Config) Logger() Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Logger{
		FilePath:  c.logger.FilePath,
		Verbose:   c.logger.Verbose,
		SystemLog: c.logger.SystemLog,
	}
}

// Notifier returns a copy of Notifier config
func (c *Config) Notifier() Notifier {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Notifier{
		Consumers: c.notifier.Consumers,
		Timeout:   c.notifier.Timeout,
	}
}

func (c *Config) Queue() Queue {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Queue{
		QueueSize:      c.queue.QueueSize,
		GoRoutinesSize: c.queue.GoRoutinesSize,
	}
}

// setService sets Service config
func (c *Config) setService() error {
	port, err := getENV(servicePortENV)
	if err != nil {
		return err
	}

	versionAPI, err := getENV(versionAPIENV)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.service = Service{
		Port:       port,
		VersionAPI: versionAPI,
	}

	return nil
}

// setDB sets DB config
func (c *Config) setDB() error {
	name, err := getENV(dbNameENV)
	if err != nil {
		return err
	}

	host, err := getENV(dbHostENV)
	if err != nil {
		return err
	}

	port, err := getENV(dbPortENV)
	if err != nil {
		return err
	}

	urlPattern, err := getENV(dbURLPatternENV)
	if err != nil {
		return err
	}

	userName, err := getENV(dbUserNmaeENV)
	if err != nil {
		return err
	}

	password, err := getENV(dbPasswordENV)
	if err != nil {
		return err
	}

	sslMode, err := getENV(dbSSLModeENV)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.db = DB{
		Name:       name,
		Host:       host,
		Port:       port,
		PatternURL: urlPattern,
		UserName:   userName,
		Password:   password,
		SSLMode:    sslMode,
	}

	return nil
}

// setLogger sets Logger config
func (c *Config) setLogger() error {
	path, err := getENV(loggerPathENV)
	if err != nil {
		return err
	}

	verbose, err := getBoolENV(loggerVerboseENV)
	if err != nil {
		return err
	}

	systemLog, err := getBoolENV(loggerSystemLogENV)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger = Logger{
		FilePath:  path,
		Verbose:   verbose,
		SystemLog: systemLog,
	}

	return nil
}

// setNotifier sets Notifier config
func (c *Config) setNotifier() error {
	consumers := Consumers{
		OnCreate: getStringSliceENV(notifierCreateConsumersENV),
		OnUpdate: getStringSliceENV(notifierUpdateConsumersENV),
		OnDelete: getStringSliceENV(notifierDeleteConsumersENV),
	}

	timeOut, err := getIntENV(notifierTimeOutENV)
	if err != nil {
		return err
	}

	retry, err := getIntENV(notifierClientMaxRetryENV)
	if err != nil {
		return err
	}

	timeOutIncreace, err := getIntENV(notifierTimeoutIncreaceENV)
	if err != nil {
		return err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.notifier = Notifier{
		Consumers:             consumers,
		Timeout:               timeOut,
		ClientMaxRetry:        retry,
		ClientTimeoutIncrease: timeOutIncreace,
	}

	return nil
}

func (c *Config) setQueue() {
	queueSize, err := getIntENV(queueSizeENV)
	if err != nil {
		queueSize = queueSizeDefault
	}

	goRoutineSize, err := getIntENV(goRoutinesSizeENV)
	if err != nil {
		goRoutineSize = goRoutinesSizeDefault
	}

	c.queue = Queue{
		QueueSize:      queueSize,
		GoRoutinesSize: goRoutineSize,
	}
}

func getENV(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("%s, %w", name, errEmptyConfiguration)
	}

	return v, nil
}

func getBoolENV(name string) (bool, error) {
	v := os.Getenv(name)
	if v == "" {
		return false, fmt.Errorf("%s, %w", name, errEmptyConfiguration)
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s, %w", name, err)
	}

	return b, nil
}

func getIntENV(name string) (int, error) {
	v := os.Getenv(name)
	if v == "" {
		return 0, fmt.Errorf("%s, %w", name, errEmptyConfiguration)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", name, err)
	}

	return i, nil
}

func getStringSliceENV(name string) []string {
	return strings.Split(os.Getenv(name), ",")
}
