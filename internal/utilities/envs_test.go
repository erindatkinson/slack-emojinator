package utilities

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/neko"
)

func TestCheckEnvs(t *testing.T) {
	tests := neko.Modern(t)

	var reqs map[string]string = map[string]string{
		"SLACK_TEAM":        "test",
		"SLACK_COOKIE":      "some-cookie;yes",
		"SLACK_TOKEN":       "xoxc-1234",
		"SLACK_CONCURRENCY": "2"}

	tests.Setup(func() {
		os.Clearenv()
		for env, val := range reqs {
			os.Setenv(env, val)
		}
	})

	tests.Cleanup(func() {
		os.Clearenv()
	})

	tests.It("collects env data when correctly set", func(t *testing.T) {
		InitConfig()
		err := CheckEnvs()
		assert.Nil(t, err)
	})

	tests.Run()
}
