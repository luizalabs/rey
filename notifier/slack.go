package notifier

import "github.com/drgarcia1986/slacker/slack"

type Slack struct {
	client   *slack.Client
	username string
	avatar   string
	channel  string
}

func (s *Slack) Notify(message string) error {
	return s.client.PostMessage(s.channel, s.username, s.avatar, message)
}

func NewSlack(token, username, avatar, channel string) *Slack {
	return &Slack{
		client:   slack.New(token),
		username: username,
		avatar:   avatar,
		channel:  channel,
	}
}
