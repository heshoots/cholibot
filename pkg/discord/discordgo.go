package discord

func (d *discordgo.session) open() {
	d.Open()
}

func (d *discordgo.session) addHandler(h interface{}) {
	d.AddHandler(h)
}

func (d *discordgo.session) addChatHandler(h interface{}) {
	d.AddHandler(h)
}

func (d *discordgo.session) channelMessageSend(channel, out string) {
	d.ChannelMessageSend(channel, out)
}
