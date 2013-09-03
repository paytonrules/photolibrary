package photolibrary

type Events interface {
  Find(eventName string) (Event, error)
}
