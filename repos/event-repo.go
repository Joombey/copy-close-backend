package repos

type EventSession struct {
	OrderCreationEventChan <-chan string
}