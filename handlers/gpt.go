package handlers

// import (
// 	"context"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// 	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
// 	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/utils"
// 	"github.com/sashabaranov/go-openai"
// )

// func HandlePostGptApi(ftx *fiber.Ctx) error {
// 	if err := utils.ValidateContentType(ftx); err != nil {
// 		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
// 	}
// 	_, err := extractEmailFromClaim(ftx)
// 	if err != nil {
// 		log.Println(err)
// 		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
// 	}

// 	var body map[string]string
// 	if err := ftx.BodyParser(&body); err != nil {
// 		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
// 	}

// 	client := openai.NewClient(cn.EnvVars.GptApiKey)
// 	response, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT3Dot5Turbo0613,
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: body["prompt"],
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("ChatCompletion error: %v\n", err)
// 		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Faild to get response from GPT API"})
// 	}
// 	log.Println(response.Choices[0].Message.Content)
// 	return ftx.Status(fiber.StatusOK).JSON(map[string]string{"GPTResp": response.Choices[0].Message.Content})
// }
