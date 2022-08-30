package main

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "log"
  "encoding/base64"
  "os"
  "path/filepath"

  "github.com/gocolly/colly"
)

type Emoji struct {
  Id     string `json:"id"`
  Code   string `json:"code"`
  Sample string `json:"sample"`
  Name   string `json:"name"`
}

func main() {
  // { "Smileys & Emotion": { "face-smiling": { "1": { "id": 1", "code": "U+1f600", "sample": "😀", "name": "grinning face" } } } }
  // [ { "Smileys & Emotion": [ { "face-smiling": [ { "id": 1", "code": "U+1f600", "sample": "😀", "name": "grinning face" } ] } ] } ]
  data := make(map[string]map[string]map[string]Emoji)
  c := colly.NewCollector()

  // Find and print all links
  c.OnHTML("table", func(e *colly.HTMLElement) {
    bighead := ""
    mediumhead := ""
    emojiId := ""

    e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
      if el.ChildAttr("th:nth-child(1)", "class") == "bighead" {
        bighead = el.ChildText("th:nth-child(1)")
        data[bighead] = make(map[string]map[string]Emoji)
      } else if el.ChildAttr("th:nth-child(1)", "class") == "mediumhead" {
        mediumhead = el.ChildText("th:nth-child(1)")
        data[bighead][mediumhead] = make(map[string]Emoji)
      } else if el.ChildAttr("th:nth-child(1)", "class") == "rchars" {
        // do nothing
      } else {
        emojiId = el.ChildText("td:nth-child(1)")
        data[bighead][mediumhead][emojiId] = Emoji{
          Id: emojiId,
          Code: el.ChildText("td:nth-child(2)"),
          Sample: el.ChildAttr("td:nth-child(3) img", "alt"),
          Name: el.ChildText("td:nth-child(4)"),
        }
        decodeBase64String(el.ChildAttr("td:nth-child(3) img", "src"), emojiId)
      }
    })
    content, err := json.Marshal(data)
    if err != nil {
      fmt.Println(err)
    }

    err = ioutil.WriteFile("data.json", content, 0644)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Println("Scraping Complete")
  })
  c.Visit("https://unicode.org/emoji/charts-14.0/emoji-list.html")
}

func decodeBase64String(s string, id string) {
  // remove the data:image/png;base64, part
  s = s[22:]

  // decode base64 string
  dec, err := base64.StdEncoding.DecodeString(s)
  if err != nil {
      panic(err)
  }

  // make directory
  newPath := filepath.Join(".", "images")
  if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
    panic(err)
  }

  // write this byte array to disk
  f, err := os.Create(`./images/` + id + `.png`)
  if err != nil {
      panic(err)
  }
  defer f.Close()

  // write bytes to file
  if _, err := f.Write(dec); err != nil {
      panic(err)
  }
  // save file
  if err := f.Sync(); err != nil {
      panic(err)
  }
}