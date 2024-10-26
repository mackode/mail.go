package main

import (
  "github.com/emersion/go-imap/v2"
	"io"
	"regexp"
	"strings"

	"github.com/DusanKasan/parsemail"
)

func (c *conn) UnreadEmails() (*imap.NumSet, error) {
  ids := new(imap.NumSet)

  // read/write!
  mbox, err := c.Cli.Select("INBOX", &imap.SelectOptions{ReadOnly: false}).Wait()
  if err != nil {
    return ids, err
  }

  c.Log.Debug("Select ok")
  c.Log.Debugw("Inbox", "messages", mbox.NumMessages)
  if mbox.NumMessages == 0 {
    c.Log.Debug("No message in mailbox")
    return ids, nil
  }

  searchCriteria := &imap.SearchCriteria{Not: []imap.SearchCriteria{{
    Flag: []imap.Flag{imap.FlagSeen},
  }}}
  
  data, err := c.Cli.UIDSearch(searchCriteria, nil).Wait()
  if err != nil {
    return ids, err
  }

  c.Log.Debugw("Unread", "msgs", data.AllUIDs())
  return &data.All, nil
}

func (c *conn) FetchEmails(ids *imap.NumSet) ([]string, error) {
  msgs := []string{}
  if ids == nil {
    c.Log.Debug("No emails")
    return msgs, nil
  }

  fetchOptions := &imap.FetchOptions{
    UID:          true,
    Envelope:     true,
    BodySection:  []*imap.FetchItemBodySection{{}},
  }

  c.Log.Debugw("Fetching", "uids", *ids)
  messages, err := c.Cli.Fetch(*ids, fetchOptions).Collect()
  if err != nil {
    c.Log.Error("Fetch failed")
    return msgs, err
  }

  c.Log.Debugw("Fetched", "msgs", len(messages))
  for _, msg := range messages {
    rawEmail := ""
    for _, buf := range msg.BodySection {
      rawEmail += string(buf)
    }
    msgs = append(msgs, rawEmail)
  }

  return msgs, nil
}

func (c *conn) ProcessEmail(rawEmail string) error {
  email, err := parsemail.Parse(strings.NewReader(rawEmail))
  if err != nil {
    return err
  }

  c.Log.Debugw("Fetched email",
      "subject", email.Subject,
      "size", len(email.HTMLBody),
      "attms", len(email.Attachments),
    )

  for _, a := range email.Attachments {
    data, err := io.ReadAll(a.Data)
    if err != nil {
      return err
    }
    c.Log.Debugw("Attachments",
      "file", a.Filename,
      "size", len(data),
      "type", a.ContentType)

    err = c.toStore(a.Filename, data)
    if err != nil {
      return err
    }
  }

  for _, e := range email.EmbeddedFiles {
    data, err := io.ReadAll(e.Data)
    if err != nil {
      return err
    }

    c.Log.Debugw("Embedded",
      "size", len(data),
      "type", e.ContentType)
  
    namerx := regexp.MustCompile(`name="(.*)"`)
    matches := namerx.FindStringSubmatch(e.ContentType)

    name := "unknown"
    if len(matches) >= 2 {
      name = matches[1]
    }
    err = c.toStore(name, data)
    if err != nil {
      return err
    }
  }

  return nil
}
