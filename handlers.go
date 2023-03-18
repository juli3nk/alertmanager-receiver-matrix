package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"maunium.net/go/mautrix"
)

func handleEvent(mcli *mautrix.Client, roomID string) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp, err := mcli.JoinRoom(roomID, "", nil)
		if err != nil {
			return err
		}

		msg := new(Message)
		if err := c.Bind(msg); err != nil {
			return err
		}

		for _, a := range msg.Alerts {
			msg := FormatAlert(a, true)

			if _, err := mcli.SendText(resp.RoomID, msg); err != nil {
				return err
			}
		}

		return c.NoContent(http.StatusCreated)
	}
}
