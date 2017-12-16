package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func countSettingsKeys(t *testing.T, m map[string]interface{}) int {
	t.Helper()

	var keys int
	for _, v := range m {
		if sub, ok := v.(map[string]interface{}); ok {
			// Don't count the object, just its keys.
			keys += len(sub)
		}

		if _, ok := v.(string); ok {
			// v is just a string key.
			keys++
		}

		if _, ok := v.([]string); ok {
			// v is just a string key.
			keys++
		}
	}
	return keys
}

func TestSetNotiDefaults(t *testing.T) {
	v := viper.New()
	setNotiDefaults(v)

	haveKeys := countSettingsKeys(t, v.AllSettings())
	if haveKeys != len(baseDefaults) {
		t.Error("Unexpected base config length")
		t.Errorf("have=%d; want=%d", haveKeys, len(baseDefaults))
	}
}

func getNotiEnv(t *testing.T) map[string]string {
	t.Helper()

	notiEnv := make(map[string]string)
	for _, env := range keyEnvBindings {
		notiEnv[env] = os.Getenv(env)
	}
	return notiEnv
}

func clearNotiEnv(t *testing.T) {
	t.Helper()

	for _, env := range keyEnvBindings {
		if err := os.Unsetenv(env); err != nil {
			t.Fatalf("failed to clear noti env: %s", err)
		}
	}
}

func setNotiEnv(t *testing.T, m map[string]string) {
	t.Helper()

	for env, val := range m {
		if err := os.Setenv(env, val); err != nil {
			t.Fatalf("failed to set noti env: %s", err)
		}
	}
}

func TestBindNotiEnv(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)

	clearNotiEnv(t)

	v := viper.New()
	bindNotiEnv(v)

	haveKeys := countSettingsKeys(t, v.AllSettings())
	if haveKeys != 0 {
		t.Fatal("Environment should be cleared")
	}

	for _, env := range keyEnvBindings {
		if err := os.Setenv(env, "foo"); err != nil {
			t.Errorf("Setenv error: %s", err)
		}
	}

	haveKeys = countSettingsKeys(t, v.AllSettings())
	wantKeys := len(baseDefaults) - 2 // -1 for message, -1 for default.
	if haveKeys != wantKeys {
		t.Error("Unexpected base config length")
		t.Errorf("have=%d; want=%d", haveKeys, wantKeys)
	}
}

func TestSetupConfigFile(t *testing.T) {
	v := viper.New()
	// For tests, we prepend the testdata dir so that we check for a config
	// file there first.
	v.AddConfigPath("testdata")
	setupConfigFile(v)

	const want = 1
	have := countSettingsKeys(t, v.AllSettings())
	if have != want {
		t.Error("Unexpected number of keys")
		t.Errorf("have=%d; want=%d", have, want)
	}
}

func TestConfigureApp(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)
	clearNotiEnv(t)

	v := viper.New()
	// For tests, we prepend the testdata dir so that we check for a config
	// file there first.
	v.AddConfigPath("testdata")
	flags := pflag.NewFlagSet("testconfigureapp", pflag.ContinueOnError)
	defineFlags(flags)

	configureApp(v, flags)

	configDir := filepath.Base(filepath.Dir(v.ConfigFileUsed()))
	if configDir != "testdata" {
		t.Fatalf("Wrong config file used: %s", v.ConfigFileUsed())
	}

	t.Run("default and file", func(t *testing.T) {
		// File takes precedence.
		have := v.GetString("nsuser.soundName")
		want := "testdata"
		if have != want {
			t.Error("Unexpected config value")
			t.Errorf("have=%s; want=%s", have, want)
		}
	})

	t.Run("default, file, and env", func(t *testing.T) {
		// Env takes precedence.
		want := "foo"
		if err := os.Setenv("NOTI_SOUND", want); err != nil {
			t.Errorf("Failed to set env: %s", err)
		}
		defer setNotiEnv(t, orig)

		have := v.GetString("nsuser.soundName")
		if have != want {
			t.Error("Unexpected config value")
			t.Errorf("have=%s; want=%s", have, want)
		}
	})

	t.Run("default", func(t *testing.T) {
		// Default takes precedence.

		// Clear config file.
		v.ReadConfig(strings.NewReader(""))

		have := v.GetString("nsuser.soundName")
		want := baseDefaults["nsuser.soundName"]
		if have != want {
			t.Error("Unexpected config value")
			t.Errorf("have=%s; want=%s", have, want)
		}
	})
}

func TestEnabledServices(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)
	clearNotiEnv(t)

	t.Run("flag override", func(t *testing.T) {
		v := viper.New()
		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		defineFlags(flags)

		want := true
		flags.Set("slack", fmt.Sprint(want))
		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["slack"]
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})

	t.Run("env override", func(t *testing.T) {
		v := viper.New()
		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		defineFlags(flags)

		if err := os.Setenv("NOTI_DEFAULT", "slack"); err != nil {
			t.Fatal(err)
		}
		defer os.Unsetenv("NOTI_DEFAULT")

		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["slack"]
		want := true
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		v := viper.New()
		// For tests, we prepend the testdata dir so that we check for a config
		// file there first.
		v.AddConfigPath("testdata")

		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		defineFlags(flags)

		configureApp(v, flags)

		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["banner"]
		want := true
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})
}