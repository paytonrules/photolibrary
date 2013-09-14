package library

type Events interface {
  Find(eventName string) (Event, error)
}
