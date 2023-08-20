package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
)

type GopherInfo struct {
	ID, X, Y string
}

func main() {
	e := echo.New()
	m := melody.New()

	e.Any("/", func(c echo.Context) error {
		m.HandleRequest(c.Response(), c.Request())
		return nil
	})
	m.HandleConnect(func(s *melody.Session) {
		ss, _ := m.Sessions()
		for _, o := range ss {
			value, exists := o.Get("info")

			if !exists {
				continue
			}
			info := value.(*GopherInfo)
			s.Write([]byte("set " + info.ID + " " + info.X + " " + info.Y))
		}

		id := uuid.NewString()
		// fmt.Printf("id: %v\n", id)
		s.Write([]byte("iam " + id))
	})
	m.HandleDisconnect(func(s *melody.Session) {
		value, exists := s.Get("info")

		if !exists {
			return
		}

		info := value.(*GopherInfo)

		m.BroadcastOthers([]byte("dis "+info.ID), s)
	})
	data := []byte("dasd")
	e.GET("/test", func(c echo.Context) error {
		m.BroadcastFilter(data, func(q *melody.Session) bool {
			fmt.Println("ion")
			return true
		})
		return nil
	})
	e.Logger.Fatal(e.Start(":8000"))
}
