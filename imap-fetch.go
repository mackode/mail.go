package main

import (
  "github.com/DusanKasan/parsemail"
  "github,com/emersion/go-imap/v2"
  "io/ioutil"
  "regexp"
  "strings"
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


