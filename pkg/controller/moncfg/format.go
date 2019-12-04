package moncfg

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/platform9/mon-operator/pkg/apis/monitoring/v1alpha1"
)

const (
	suffixLen = 8
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())
}

type format interface {
	formatAlert(recv *v1alpha1.Receivers, acfg *alertConfig) error
}

func getFormatter(ftype string) (format, error) {
	var f format
	switch ftype {
	case "slack":
		f = slackconfig{}
	case "email":
		f = emailconfig{}
	default:
		return nil, os.ErrInvalid
	}

	return f, nil
}

func formatReceiver(moncfg *v1alpha1.MonCfg, acfg *alertConfig) error {

	for _, recv := range moncfg.Spec.Alertmanager.Receivers {
		var f format
		log.Info("Parsing Receivers: ", "Type", recv.Type)
		f, err := getFormatter(recv.Type)
		if err != nil {
			return err
		}

		err = f.formatAlert(&recv, acfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f slackconfig) formatAlert(recv *v1alpha1.Receivers, acfg *alertConfig) error {
	var url, channel, severity string
	for _, param := range recv.Params {
		switch param.Name {
		case "url":
			url = param.Value
		case "channel":
			channel = param.Value
		case "severity":
			severity = param.Value
		}
	}

	if url == "" {
		return os.ErrInvalid
	}

	if channel == "" {
		return os.ErrInvalid
	}

	if severity == "" {
		return os.ErrInvalid
	}

	if acfg.Route.Routes == nil {
		acfg.Route.Routes = []routes{}
	}
	receiverName := fmt.Sprintf("%s-%s", "slack", RandString(suffixLen))

	acfg.Route.Routes = append(acfg.Route.Routes, routes{
		Receiver: receiverName,
		MatchRe: map[string]string{
			"severity": severity,
		},
	})

	acfg.Receivers = append(acfg.Receivers, receiver{
		Name: receiverName,
		SlackConfigs: []slackconfig{
			slackconfig{
				ApiURL:  url,
				Channel: channel,
			},
		},
	})

	return nil
}

func (f emailconfig) formatAlert(recv *v1alpha1.Receivers, acfg *alertConfig) error {
	var to, from, smarthost, severity string
	var auth_identity, auth_username, auth_password string
	for _, param := range recv.Params {
		switch param.Name {
		case "to":
			to = param.Value
		case "from":
			from = param.Value
		case "smarthost":
			smarthost = param.Value
		case "severity":
			severity = param.Value
		case "auth_identity":
			auth_identity = param.Value
		case "auth_username":
			auth_username = param.Value
		case "auth_password":
			auth_password = param.Value
		}
	}

	if to == "" {
		return os.ErrInvalid
	}

	if from == "" {
		return os.ErrInvalid
	}

	if smarthost == "" {
		return os.ErrInvalid
	}

	if severity == "" {
		return os.ErrInvalid
	}

	if auth_identity == "" {
		return os.ErrInvalid
	}

	if auth_username == "" {
		return os.ErrInvalid
	}

	if auth_password == "" {
		return os.ErrInvalid
	}

	if acfg.Route.Routes == nil {
		acfg.Route.Routes = []routes{}
	}

	receiverName := fmt.Sprintf("%s-%s", "email", RandString(suffixLen))
	acfg.Route.Routes = append(acfg.Route.Routes, routes{
		Receiver: receiverName,
		MatchRe: map[string]string{
			"severity": severity,
		},
	})

	acfg.Receivers = append(acfg.Receivers, receiver{
		Name: receiverName,
		EmailConfigs: []emailconfig{
			emailconfig{
				To:           to,
				From:         from,
				SmartHost:    smarthost,
				AuthUsername: auth_username,
				AuthIdentity: auth_identity,
				AuthPassword: auth_password,
			},
		},
	})

	return nil
}

// RandString returns a random string of size "n"
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
