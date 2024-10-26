package main

import (
  "errors"
  "fmt"
  "os"
  "path/filepath"
  "strings"
)

func (c *conn) toStore(fpath string, data []byte) error {
  var err error
  home, err := os.UserHomeDir()
  if err != nil {
    return err
  }

  photoDir := filepath.Join(home, "photos")
  os.Mkdir(photoDir, 0755)
  base := filepath.Base(fpath)
  npath := filepath.Join(photoDir, base)

  var f *os.File
  if _, err := os.Stat(npath); errors.Is(err, os.ErrNotExist) {
    file, e := os.Create(npath)
    f = file
    err = e
  } else {
    suffix := filepath.Ext(base)
    prefix := strings.TrimSuffix(base, suffix)
    f, err = os.CreateTemp(photoDir, fmt.Sprintf("%s-*%s", prefix, suffix))
  }

  if err != nil {
    return err
  }

  c.Log.Debugw("Write", "name", f.Name(), "size", len(data))
  _, err = f.Write(data)
  if err != nil {
    return err
  }

  err = f.Close()
  if err != nil {
    return err
  }

  return nil
}
