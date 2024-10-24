package main

import (
  "github.com/emersion/go-imap/v2/imapclient"
  "go.uber.org/zap"
)

type conn struct {
  HostPort  string
  User      string
  Pass      string
  Cli       *imapclient.Client
  Log       *zap.SugaredLogger
}

func NewIMAP() *conn {
  c = conn{
    HostPort: "imap.foo.com:993",
    User:     "me@foo.com",
    Pass:     "PASSWORD",
  }

  c.Log = zap.NewExample().Sugar()
  return &c
}

func (c *conn) Close() {
  c.Log.Debug("Closing connection")
  c.Log.Sync()
  c.Cli.Close()
}

func (c *conn) Open() error {
  c.Log.Debugw("Connecting", "host", c.HostPort)
  cli, err := imapclient.DialTLS(c.HostPort, nil)
  c.Cli = cli
  if err != nil {
    return 0, err
  }

  c.Log.Debug("Connect OK")
  c.Log.Debugw("Login", "user", c.User)
  if err := c.Cli.Login(c.User, c.Pass).Wait(); err != nil {
    return 0, err
  }

  c.Log.Debug("Login OK")
  return nil
}

