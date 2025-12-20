package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/jsonc"
)

const (
	localConfigFile    = "./config.json"
	localConfigFileAlt = "./config.jsonc"
)

type configManager struct {
	goValidator   *validator.Validate
	decoderConfig mapstructure.DecoderConfig
}

func GetConfig(ctx context.Context) *Application {
	var cfg Application
	config := &configManager{
		goValidator: validator.New(),
		decoderConfig: mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				stringToSlogLevelHookFunc(),
			),
		},
	}

	err := config.readLocal(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func (c *configManager) readLocal(ctx context.Context, destination any) error {
	// attempt to read from .jsonc to so that JSON files can have docs, too:
	useJSONc := true
	localConfRaw, err := os.ReadFile(localConfigFileAlt)
	if err != nil {
		// .jsonc was unreadable; read from normal .json file:
		localConfRaw, err = os.ReadFile(localConfigFile)
		if err != nil {
			return err
		}
		useJSONc = false
	}

	// read JSON into map:
	var localMap map[string]any
	if useJSONc {
		localConfRaw = jsonc.ToJSON(localConfRaw)
	}

	if err = json.Unmarshal(localConfRaw, &localMap); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// decode Static configurations
	if err = c.decodeAndValidate(ctx, localMap, destination); err != nil {
		return err
	}

	return nil
}

func (c *configManager) decodeAndValidate(ctx context.Context, input any, destination any) error {
	if err := c.decode(input, &destination); err != nil {
		return fmt.Errorf("failed to decode values into %T: %w", destination, err)
	}

	if err := c.goValidator.StructCtx(ctx, destination); err != nil {
		result, errJson := json.Marshal(destination)
		if errJson != nil {
			return fmt.Errorf("failed json marshal: %w", err)
		}
		return fmt.Errorf("failed to validate '%s' in %T: %w", string(result), destination, err)
	}
	return nil
}

func (c *configManager) decode(input any, destination any) (err error) {
	defer func() {
		if panicked := recover(); panicked != nil {
			err = fmt.Errorf("panicked when decoding config: %v", panicked)
		}
	}()
	dc := c.decoderConfig
	dc.Result = destination
	mdc, err := mapstructure.NewDecoder(&dc)
	if err != nil {
		return err
	}
	return mdc.Decode(input)
}

func stringToSlogLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(
		s reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if s.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(slog.Level(0)) {
			return data, nil
		}

		switch data.(string) {
		case "DEBUG":
			return slog.LevelDebug, nil
		case "INFO":
			return slog.LevelInfo, nil
		case "WARN":
			return slog.LevelWarn, nil
		case "ERROR":
			return slog.LevelError, nil
		default:
			return slog.LevelDebug, fmt.Errorf("unrecognized slog.Level '%s', valid values are: DEBUG, INFO, WARN, ERROR", data)
		}
	}
}
