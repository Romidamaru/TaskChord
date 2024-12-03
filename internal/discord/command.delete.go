package discord

func (b *Bot) DeleteCommands() error {
	cmds, err := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", cmd.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
