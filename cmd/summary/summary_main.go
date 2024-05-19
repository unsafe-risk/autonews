package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lemon-mint/godotenv"
	"github.com/lemon-mint/vermittlungsstelle/llm/generativelanguage"
	"github.com/unsafe-risk/autonews"
)

func main() {
	godotenv.Load()
	client, err := generativelanguage.NewClient(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "creating client:", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Print("URL >>")
	url := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Printf("Crawling %s...\n", url)

	an := autonews.NewAutoNews(generativelanguage.NewModel(client, "gemini-1.5-pro-latest", nil))

	summary, err := an.GetSummary(context.Background(), strings.TrimSpace(url))
	if err != nil {
		fmt.Fprintln(os.Stderr, "getting summary:", err)
		os.Exit(1)
	}

	fmt.Println("SUMMARY:")
	fmt.Println(summary)

	fmt.Println("Translating...")
	translated, err := an.Translate(context.Background(), summary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "translating:", err)
		os.Exit(1)
		return
	}
	fmt.Println("TRANSLATED:")
	fmt.Println(translated)
}
