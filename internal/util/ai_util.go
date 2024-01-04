package util

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ExecuteQuery(text string) string {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text("Generate 10 Multiple Choice Questions Based on the text below with Answers  "+"\n"+"\n"+text))
	if err != nil {
		log.Fatal(err)
	}

	questions := printResponse(resp)

	fmt.Println(questions)

	return questions

}

func printResponse(resp *genai.GenerateContentResponse) string {
	response := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response += fmt.Sprintf("%v ", part)
			}
		}
	}

	// fmt.Println(response)
	return response

}
