package utilities

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var Envs []string = []string{
	"team",
	"cookie",
	"token",
	"concurrency",
}

func CheckEnvs() error {
	errors := []string{}
	for _, env := range Envs {
		if viper.Get(env) == nil {
			errors = append(
				errors,
				strings.ToUpper(
					fmt.Sprintf(
						"%s_%s",
						viper.GetEnvPrefix(),
						strings.ReplaceAll(env, "-", "_"))))
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("envs required: [%s]", strings.Join(errors, ", "))
	}

	return nil
}
