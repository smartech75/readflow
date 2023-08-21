package mail

import (
	"context"
	"io"
	"mime"
	"net/mail"
	"strings"

	"github.com/emersion/go-smtp"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrActionAborted = &smtp.SMTPError{
		Code:         451,
		EnhancedCode: smtp.EnhancedCode{4, 5, 1},
		Message:      "Action aborted: Local error in processing",
	}
	ErrMailboxUnavailable = &smtp.SMTPError{
		Code:         550,
		EnhancedCode: smtp.EnhancedCode{5, 5, 0},
		Message:      "Requested action not taken: mailbox unavailable",
	}
	ErrBadSequence = &smtp.SMTPError{
		Code:         503,
		EnhancedCode: smtp.EnhancedCode{5, 0, 3},
		Message:      "Bad sequence of commands",
	}
)

type Session struct {
	from   string
	ctx    context.Context
	logger zerolog.Logger
}

func (s *Session) AuthPlain(username, password string) error {
	return smtp.ErrAuthUnsupported
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// TODO limit body size
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// decode recipient
	dest := strings.Split(to, "@")[0]
	parts := strings.Split(dest, "-")
	hash := parts[len(parts)-1]
	alias := strings.TrimSuffix(dest, "-"+hash)

	// retrieve the corresponding user
	user, err := service.Lookup().GetUserByHashID(s.ctx, hash)
	if err != nil {
		return ErrMailboxUnavailable
	}
	s.ctx = context.WithValue(s.ctx, constant.ContextUser, *user)
	s.ctx = context.WithValue(s.ctx, constant.ContextUserID, *user.ID)
	s.ctx = context.WithValue(s.ctx, constant.ContextIsAdmin, false)
	s.logger = s.logger.With().Uint("uid", *user.ID).Logger()

	// retrieve the corresponding incoming webhook
	// TODO escape alias? so find another way to retrieve the webhook (full scan ?)
	incomingWebhook, err := service.Lookup().GetIncomingWebhookByAlias(s.ctx, alias)
	if err != nil || incomingWebhook == nil {
		return ErrMailboxUnavailable
	}
	s.ctx = context.WithValue(s.ctx, constant.ContextIncomingWebhook, incomingWebhook)
	s.logger = s.logger.With().Str("webhook", incomingWebhook.Alias).Logger()

	return nil
}

func (s *Session) Data(r io.Reader) error {
	if s.ctx.Value(constant.ContextIncomingWebhook) == nil {
		return ErrBadSequence
	}
	s.logger.Debug().Msg("receiving mail using incoming webhook...")
	msg, err := mail.ReadMessage(r)
	if err != nil {
		s.logger.Error().Err(err).Msg("unable to read message")
		return ErrActionAborted
	}
	// extract Media Type
	contentType := msg.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		s.logger.Error().Err(err).Msg("invalid Content-Type")
		return ErrActionAborted
	}
	// TODO support multipart
	if !(strings.HasPrefix(mediaType, "text/plain") || strings.HasPrefix(mediaType, "text/html")) {
		s.logger.Error().Str("type", mediaType).Msg("invalid Media-Type")
		return ErrActionAborted
	}
	// read body
	// TODO limit body size
	b, err := io.ReadAll(msg.Body)
	if err != nil {
		s.logger.Error().Err(err).Msg("unable to read body")
		return ErrActionAborted
	}

	// build article
	builder := model.NewArticleCreateFormBuilder()
	builder.Origin(msg.Header.Get("From"))
	builder.Title(msg.Header.Get("Subject"))
	if strings.HasPrefix(mediaType, "text/plain") {
		builder.Text(string(b))
	} else {
		builder.HTML(string(b))
	}

	// create article
	_, err = service.Lookup().CreateArticle(s.ctx, *builder.Build(), service.ArticleCreationOptions{IgnoreHydrateError: true})
	if err != nil {
		s.logger.Error().Err(err).Msg("unable to create article")
		return ErrActionAborted
	}

	return nil
}

func (s *Session) Reset() {
	s.ctx = context.Background()
}

func (s *Session) Logout() error {
	s.ctx = nil
	return nil
}

func NewSession() *Session {
	logger := log.With().Str("component", "mail").Logger()
	return &Session{
		from:   "",
		ctx:    context.Background(),
		logger: logger,
	}
}