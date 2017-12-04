// Pacakge tgwriter wraps telegram-bot-api as a writer foor logging purpose
//
//     w, err := tgwriter.ToChannel("@mychan", "my-secret-token")
//     if err != nil {
//         // error handling
//     }
//     l := log.New(w, "[MYAPP]", 0)
//     log.Printf("some error happened: %s", errors.new("my error"))
//
// WARNING: WRITING ERRORS ARE SILENTLY IGNORED BY log.New
package tgwriter

import (
	"io"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type writer struct {
	api   *tgbotapi.BotAPI
	maker func(string) tgbotapi.MessageConfig
}

func (w *writer) Write(data []byte) (n int, err error) {
	_, err = w.api.Send(w.maker(string(data[:len(data)-1])))
	if err == nil {
		n = len(data)
	}

	return
}

// ToUser creates a writer send messages to a user or private channel
func ToUser(uid int64, token string) (io.Writer, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &writer{
		api: api,
		maker: func(s string) tgbotapi.MessageConfig {
			return tgbotapi.NewMessage(uid, s)
		},
	}, nil
}

// ToChannel creates a writer send messages to a public channel or a named user
func ToChannel(username, token string) (io.Writer, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &writer{
		api: api,
		maker: func(s string) tgbotapi.MessageConfig {
			return tgbotapi.NewMessageToChannel(username, s)
		},
	}, nil
}
