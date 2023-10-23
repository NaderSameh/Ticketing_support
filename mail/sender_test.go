package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	sender := NewGmailSender("Cypod", "cypodsolutionsticket@gmail.com", "qfvimvtxllahvkwy")

	subject := "new open ticket"
	content := `
	<h1>Customer issue</h1>
	<p>Please find a new ticket <a href="http://cypod-test.web.app">Cypodsolutions</a></p>
	`
	to := []string{"naders@cypodsolutions.com"}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
