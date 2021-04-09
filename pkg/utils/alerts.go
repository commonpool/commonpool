package utils

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

type Alert struct {
	Class   string
	Message string
}

func (a Alert) HTML() template.HTML {
	return template.HTML(a.Message)
}

type AlertManager struct {
	store sessions.Store
}

func NewAlertManager(store sessions.Store) *AlertManager {
	return &AlertManager{
		store: store,
	}
}

func (a *AlertManager) AddAlert(r *http.Request, w http.ResponseWriter, alert Alert) error {
	session, err := a.store.Get(r, "alerts")
	if err != nil {
		return err
	}
	alerts, err := a.GetAlerts(r)
	if err != nil {
		return err
	}

	alerts = append(alerts, alert)
	session.Values["alerts"] = alerts
	return session.Save(r, w)
}

func (a *AlertManager) GetAlerts(r *http.Request) ([]Alert, error) {
	session, err := a.store.Get(r, "alerts")
	if err != nil {
		return nil, err
	}
	alertsIntf, ok := session.Values["alerts"]
	if !ok {
		return []Alert{}, nil
	}
	alerts := alertsIntf.([]Alert)
	return alerts, nil
}

func (a *AlertManager) ClearAlerts(r *http.Request, w http.ResponseWriter) error {
	session, err := a.store.Get(r, "alerts")
	if err != nil {
		return err
	}
	session.Values["alerts"] = []Alert{}
	return session.Save(r, w)
}
