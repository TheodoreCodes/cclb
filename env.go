package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Env struct {
	Port       int    `env:"PORT" default:"8080"` // Port where to listen for incoming connections
	Timeout    string `env:"TIMEOUT" default:"30s"`
	LogLevel   string `json:"LOG_LEVEL" default:"DEBUG"`
	ConfigFile string `json:"CONFIG_FILE" default:"./conf.json"`
}

func (rv *Env) Load() error {
	st := reflect.TypeOf(*rv)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		envName := field.Tag.Get("env")
		if envVal := os.Getenv(envName); envVal != "" {
			rv.SetVal(i, envVal)
		} else if envVal = field.Tag.Get("default"); envVal != "" {
			rv.SetVal(i, envVal)
		} else {
			return errors.New(fmt.Sprintf("no value available for `%s`", field.Name))
		}
	}

	return nil
}

func (rv *Env) SetVal(idx int, val string) {
	f := reflect.ValueOf(rv).Elem().Field(idx)
	switch f.Type().Name() {
	case "string":
		f.SetString(val)
		break
	case "int":
		intVal, _ := strconv.Atoi(val)
		f.SetInt(int64(intVal))
	}
}

func LoadEnv() (*Env, error) {
	env := Env{}
	err := env.Load()

	if err != nil {
		return nil, err
	}

	return &env, nil
}
