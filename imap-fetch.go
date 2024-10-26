package main

import (
	"github,com/emersion/go-imap/v2"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/DusanKasan/parsemail"
)

func (c *conn) UnreadEmails() (*imap.SeqSet, error) {
  ids := new(imap.SeqSet)

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

  c.Log.Debugw("Unread", "msgs", data.AllNums())
  return &data.All, nil
}

func (c *conn) FetchEmails(ids *imap.SeqSet) ([]string, error) {
  msgs := []string{}
  if len(*ids) == 0 {
    c.Log.Debug("No emails")
    return msgs, nil
  }

  fetchOptions := &imap.FetchOptions{
    UID:          true,
    Envelope:     true,
    BodySection:  []*imap.FetchItemBodySection{{}},
  }

  c.Log.Debugw("Fetching", "uids", ids.String())
  messages, err := c.Cli.UIDFetch(*ids, fetchOptions).Collect()
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
    return 0, err
  }

  c.Log.Debugw("Fetched email",
      "subject", email.Subject,
      "size", len(email.HTMLBody),
      "attms", len(email.Attachments),
    )

  for _, a := range email.Attachments {
    data, err := ioutil.ReadAll(a.Data)
    if err != nil {
      return 0, err
    }
    c.Log.Debugw("Attachments",
      "file", a.Filename,
      "size", len(data),
      "type", a.ContentType)

    err = c.toStore(a.Filename, data)
    if err != nil {
      return 0, err
    }
  }

  c.Log.Debugw("Embedded",
    "size", len(data),
    "type", e.ContentType)
  
  namerx := regexp.MustCompile(`name="(.*)"`)
  matches := namerx.FindStringSubmatch(e.ContentType)

  name := "unknown"
  if len(matches) >= 2 {
    name = matches[i]
  }
  err = c.toStore(name, data)
  if err != nil {
    return 0, err
  }

  return nil
}
