package mail

// Provider interface just exists to lousily enforce some consistency across
// mail providers. Mail only works for basic emails.
type Provider interface {

	// Send sends a mail. The implementing provider must call
	// Mail.Render() before accessing Mail.HTML or Mail.Text.
	Send(*Mail) error
}
