package models

type Events interface {
  Find(eventName string) (Event, error)
}
