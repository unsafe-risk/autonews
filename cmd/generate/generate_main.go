package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/lemon-mint/godotenv"
	"github.com/lemon-mint/vermittlungsstelle/llm"
	"github.com/lemon-mint/vermittlungsstelle/llm/generativelanguage"
	"github.com/unsafe-risk/autonews"
)

type Section struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Post struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`

	Urls     []string  `json:"urls"`
	Sections []Section `json:"sections"`
}

func mkpath(title string) string {
	runes := []rune(title)
	for i, r := range runes {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			runes[i] = r
		} else {
			runes[i] = '-'
		}
	}
	return string(runes)
}

func timeID(t time.Time) string {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(t.Unix()))
	return hex.EncodeToString(b[:])
}

func main() {
	godotenv.Load()
	client, err := generativelanguage.NewClient(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "creating client:", err)
		os.Exit(1)
	}
	defer client.Close()

	var post Post
	post.Date = time.Now()

	post_json, err := os.ReadFile("post.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(post_json, &post)
	if err != nil {
		panic(err)
	}

	an := autonews.NewAutoNews(generativelanguage.NewModel(client, "gemini-1.5-flash-latest", &llm.Config{
		SafetyFilterThreshold: llm.BlockOnlyHigh,
	}))

	for _, url := range post.Urls {
		title, text, err := an.GeneratePost(context.Background(), strings.TrimSpace(url))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err, "URL:", url)
			os.Exit(1)
		}
		fmt.Printf("%s:\n\n%s\n\n", title, text)
		post.Sections = append(post.Sections, Section{Title: title, Text: text, URL: url})
	}

	var post_body strings.Builder
	for _, section := range post.Sections {
		post_body.WriteString("## [")
		post_body.WriteString(section.Title)
		post_body.WriteString("]")
		post_body.WriteString("(")
		post_body.WriteString(section.URL)
		post_body.WriteString(")\n\n")
		post_body.WriteString(section.Text)
		post_body.WriteString("\n\n")
	}

	title, err := an.NewsTitle(context.Background(), post_body.String())
	if err != nil {
		panic(err)
	}
	title = strings.TrimSpace(title)
	title = mkpath(title)
	title = strings.ReplaceAll(title, "--", "-")
	title = strings.ReplaceAll(title, "--", "-")
	fmt.Println("Title:", title)
	post.Title = title

	basepath := fmt.Sprintf("articles/%s-%s", timeID(post.Date), post.Title)
	err = os.MkdirAll(basepath, 0755)
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/index.json", basepath), data, 0644)
	if err != nil {
		panic(err)
	}

	var postmd strings.Builder

	type metadata struct {
		Title string    `yaml:"title"`
		Date  time.Time `yaml:"date"`
	}

	md := metadata{Title: post.Title, Date: post.Date}
	data, err = yaml.Marshal(md)
	if err != nil {
		panic(err)
	}

	postmd.WriteString("---\n")
	postmd.WriteString(string(data))
	postmd.WriteString("---\n\n")
	postmd.WriteString(post_body.String())
	postmd.WriteString("\n")

	post_en := strings.TrimSpace(postmd.String())

	err = os.WriteFile(fmt.Sprintf("%s/post.md", basepath), []byte(post_en), 0644)
	if err != nil {
		panic(err)
	}

	post_ko, err := an.Translate(context.Background(), post_en)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/post.ko.md", basepath), []byte(post_ko), 0644)
	if err != nil {
		panic(err)
	}
}
