package main

import (
	"time"
)

type Server struct {
	Http   Http
	Matrix Matrix
}

type Http struct {
	BindPort         int      `split_words:"true" default:"8080"`
	CorsAllowOrigins []string `split_words:"true"`
	AuthTokenFile    string   `split_words:"true" default:"/var/run/secrets/alertmanager/receiver-matrix-token"`
}

type Matrix struct {
	HomeserverUrlFile string `split_words:"true" default:"/run/secrets/matrix/homeserver-url"`
	UserNameFile      string `split_words:"true" default:"/run/secrets/matrix/user-name"`
	UserPasswordFile  string `split_words:"true" default:"/run/secrets/matrix/user-password"`
	RoomIdFile        string `split_words:"true" default:"/run/secrets/matrix/room-id"`
}

type Message struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []*Alert          `json:"alerts"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}
