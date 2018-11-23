package notifier

type Notifier interface {
	Notify(message string) error
}

type fake struct {
}

func (f *fake) Notify(message string) error {
	return nil
}

func New(token, username, avatar, channel string) Notifier {
	if token == "" {
		return &fake{}
	}
	return NewSlack(token, username, avatar, channel)
}
