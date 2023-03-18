package main

import (
	"fmt"
	"strings"
)

func CreateMessageText(status, name, summary, id string) string {
	var result string

	if len(id) > 0 {
		id = fmt.Sprintf(" (%s)", id)
	}

	result = fmt.Sprintf("%s %s: %s%s", strings.ToUpper(status), name, summary, id)

	return result
}

func FormatAlert(alert *Alert, labels bool) string {
	var result string

	status := "alert"
	if alert.Status == "resolved" {
		status = "resolved"
	} else if alert.Status == "suppressed" {
		status = "silenced"
	} else if sev, ok := alert.Labels["severity"]; ok {
		status = sev
	}

	var summary string
	if v, ok := alert.Annotations["summary"]; ok {
		summary = v
	}

	var alertName string
	if v, ok := alert.Labels["alertname"]; ok {
		alertName = v
	}

	// Format main message
	result = CreateMessageText(status, alertName, summary, alert.Fingerprint)

	// Add labels
	if labels {
		var lbls []string

		for n, v := range alert.Labels {
			lbls = append(lbls, fmt.Sprintf(`%s="%s"`, n, v))
		}

		result += ", labels:" + strings.Join(lbls, " ")
	}

	return result
}
