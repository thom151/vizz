package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func genThread(client *openai.Client, book string) (openai.Thread, error) {
	ctx := context.Background()

	threadRequest := openai.ThreadRequest{
		Messages: []openai.ThreadMessage{
			{
				Role:    "user",
				Content: "The book is " + book,
			},
		},
	}

	thread, err := client.CreateThread(ctx, threadRequest)
	if err != nil {
		return openai.Thread{}, err
	}
	return thread, nil

}

func genMessage(c *openai.Client, thread_id, chapter string) ([]string, error) {
	request := openai.MessageRequest{
		Role:    "user",
		Content: "This is the chapter: " + chapter,
	}
	message, err := c.CreateMessage(context.Background(), thread_id, request)
	if err != nil {
		return nil, err
	}
	fmt.Println("Received message content: ", message.Content)
	return nil, nil
}

func getRun(c *openai.Client, thread_id, assistant_id string) (openai.Run, error) {
	run, err := c.CreateRun(context.Background(), thread_id, openai.RunRequest{
		AssistantID: assistant_id,
	})

	if err != nil {
		return openai.Run{}, err
	}
	return run, nil
}

func getResponse(c *openai.Client, thread_id, run_id string) ([]string, error) {
	for {
		run, err := c.RetrieveRun(context.Background(), thread_id, run_id)
		if err != nil {
			return nil, err
		}

		if run.Status == "completed" {
			messagesList, err := c.ListMessage(context.Background(), thread_id, nil, nil, nil, nil, nil)
			if err != nil {
				return nil, err
			}

			message := messagesList.Messages[0].Content[0].Text.Value
			prompts := strings.Split(message, "\n\n")
			return prompts, nil

		}
	}
}

func genImageBase64(c *openai.Client, text string) (string, error) {

	resp, err := c.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Model:          openai.CreateImageModelDallE3,
			Quality:        "standard",
			Prompt:         "A cartoon version of " + text + "There should be no text",
			ResponseFormat: openai.CreateImageResponseFormatURL,
			N:              1,
		},
	)

	if err != nil {
		return "", err
	}

	fmt.Println("IMG URLS: ", resp.Data[0].URL)
	return resp.Data[0].URL, nil

}
