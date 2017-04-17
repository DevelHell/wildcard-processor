package wildcard_processor

import (
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/flashmob/go-guerrilla/response"
	"path/filepath"
	"strings"
)

type WildcardConfig struct {
	WildcardHosts string `json:"wildcard_hosts"`
}

type wildcard struct {
	config       *WildcardConfig
	allowedHosts map[string]bool
}

func newWildcardProcessor(c *WildcardConfig) *wildcard {
	w := wildcard{config: c}
	w.allowedHosts = make(map[string]bool)
	for _, h := range strings.Split(c.WildcardHosts, ",") {
		w.allowedHosts[h] = true
	}

	return &w
}

func (w *wildcard) allowsRcpt(rcpt string) bool {
	if _, ok := w.allowedHosts["*"]; ok {
		return true
	}
	if _, ok := w.allowedHosts[rcpt]; ok {
		return true
	}
	for pattern := range w.allowedHosts {
		matched, err := filepath.Match(pattern, rcpt)
		if err != nil {
			return false
		}
		if matched {
			return true
		}
	}
	return false
}

var WildcardProcessor = func() backends.Decorator {

	var (
		w *wildcard
	)

	initializer := backends.InitializeWith(func(backendConfig backends.BackendConfig) error {
		configType := backends.BaseConfig(&WildcardConfig{})
		bcfg, err := backends.Svc.ExtractConfig(backendConfig, configType)
		if err != nil {
			return err
		}
		c := bcfg.(*WildcardConfig)
		w = newWildcardProcessor(c)

		return nil
	})

	backends.Svc.AddInitializer(initializer)

	return func(c backends.Processor) backends.Processor {
		return backends.ProcessWith(func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
			if task == backends.TaskValidateRcpt {
				// all recipients must pass
				for _, rcpt := range e.RcptTo {
					if !w.allowsRcpt(rcpt.Host) {
						backends.Log().Debugf("Recipients host %s did not pass", rcpt.Host)
						return backends.NewResult(
								response.Canned.FailNoSenderDataCmd),
							backends.NoSuchUser
					}
				}
			}
			return c.Process(e, task)
		})
	}
}
