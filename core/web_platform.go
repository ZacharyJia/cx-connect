package core

import "context"

type webAdminPlatform struct{}

func newWebAdminPlatform() Platform {
	return &webAdminPlatform{}
}

func (p *webAdminPlatform) Name() string {
	return "web"
}

func (p *webAdminPlatform) Start(handler MessageHandler) error {
	return nil
}

func (p *webAdminPlatform) Reply(ctx context.Context, replyCtx any, content string) error {
	return nil
}

func (p *webAdminPlatform) Send(ctx context.Context, replyCtx any, content string) error {
	return nil
}

func (p *webAdminPlatform) ReplyWithButtons(ctx context.Context, replyCtx any, content string, buttons []Button) error {
	return nil
}

func (p *webAdminPlatform) Stop() error {
	return nil
}

func (p *webAdminPlatform) ReconstructReplyCtx(sessionKey string) (any, error) {
	return sessionKey, nil
}
