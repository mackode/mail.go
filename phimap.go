package main

func main() {
  c := NewIMAP()
  err := c.Open()
  if err != nil {
    c.Log.Fatalw("conn", err)
  }
  defer c.Close()

  ids, err := c.UnreadEmails()
  if err != nil {
    c.Log.Fatalw("List", err)
  }

  emails, err := c.FetchEmails(ids)
  if err != nil {
    c.Log.Fatalw("Fetch", err)
  }

  for _, email := range emails {
    err := c.ProcessEmail(email)
    if err != nil {
      c.Log.Fatalw("Parse", err)
    }
  }
}
