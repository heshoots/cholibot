package discord

type session interface {
	open()
	addChatHandler(interface{})
	addHandler(interface{})
	ChannelMessageSend(channel, out string)
}

type HandlerFunc func(s *session, m interface{})

func (f HandlerFunc) Call(s *session, m interface{}) {
	f(w, r)
}
