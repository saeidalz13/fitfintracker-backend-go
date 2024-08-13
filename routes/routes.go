package routes

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	h "github.com/saeidalz13/LifeStyle2/lifeStyleBack/handlers"
	m "github.com/saeidalz13/LifeStyle2/lifeStyleBack/middlewares"
	"github.com/sashabaranov/go-openai"
)

func AuthSetup(app *fiber.App, hc *h.HandlersManager) {
	app.Get(cn.URLS.OAuthSignIn, h.HandleGetGoogleSignIn)
	app.Get(cn.URLS.OAuthCallback, hc.AuthHandler.HandleGetGoogleCallback)
	app.Get(cn.URLS.Home, hc.AuthHandler.HandleGetHome)
	app.Get(cn.URLS.Profile, hc.AuthHandler.HandleGetProfile)
	app.Get(cn.URLS.SignOut, h.HandleGetSignOut)

	app.Post(cn.URLS.SignUp, hc.AuthHandler.HandlePostSignUp)
	app.Post(cn.URLS.Login, hc.AuthHandler.HandlePostLogin)

	app.Delete(cn.URLS.DeleteProfile, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.AuthHandler.HandleDeleteUser)
}

func FitnessSetup(app *fiber.App, hc *h.HandlersManager) {
	app.Get(cn.URLS.FetchSinglePlan, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetSinglePlan)
	app.Get(cn.URLS.AllPlans, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetAllFitnessPlans)
	app.Get(cn.URLS.AllDayPlans, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetAllFitnessDayPlans)
	app.Get(cn.URLS.AllDayPlanMoves, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetAllFitnessDayPlanMoves)
	app.Get(cn.URLS.FetchDayPlanMovesWorkout, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetAllFitnessDayPlanMovesWorkout)
	app.Get(cn.URLS.FetchPlanRecords, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetPlanRecords)
	app.Get(cn.URLS.FetchWeekPlanRecords, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetWeekPlanRecords)
	app.Get(cn.URLS.FetchNumAvailableWeeksPlanRecords, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetNumAvailableWeeksPlanRecords)
	app.Get(cn.URLS.FetchCurrentWeekCompletedExercises, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetCurrentWeekCompletedExercises)
	app.Get(cn.URLS.FetchRecordedTime, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleGetRecordedTime)

	app.Post(cn.URLS.AddPlan, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandlePostAddPlan)
	app.Post(cn.URLS.EditPlan, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandlePostEditPlan)
	app.Post(cn.URLS.AddDayPlanMoves, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandlePostAddDayPlanMoves)
	app.Post(cn.URLS.AddPlanRecord, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandlePostAddPlanRecord)
	app.Post(cn.URLS.AddPlanRecordedTime, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandlePostRecordedTime)

	app.Delete(cn.URLS.DeletePlan, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleDeletePlan)
	app.Delete(cn.URLS.DeleteDayPlan, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleDeleteDayPlan)
	app.Delete(cn.URLS.DeleteDayPlanMove, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleDeleteDayPlanMove)
	app.Delete(cn.URLS.DeleteWeekPlanRecords, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.HandleDeleteWeekFromPlanRecords)
	app.Delete(cn.URLS.DeletePlanRecord, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.DeleteSetFromPlanRecord)

	app.Patch(cn.URLS.UpdatePlanRecord, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FitnessHandler.PatchPlanRecord)
}

func FinanceSetup(app *fiber.App, hc *h.HandlersManager) {
	app.Get(cn.URLS.ShowBudgets, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetAllBudgets)
	app.Get(cn.URLS.EachBalance, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetSingleBalance)
	app.Get(cn.URLS.EachBudget, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetSingleBudget)

	app.Post(cn.URLS.PostNewBudget, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandlePostNewBudget)
	app.Post(cn.URLS.EachExpense, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandlePostExpenses)
	app.Get(cn.URLS.CapitalExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetCapitalExpenses)
	app.Get(cn.URLS.EatoutExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetEatoutExpenses)
	app.Get(cn.URLS.EntertainmentExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleGetEntertainmentExpenses)

	app.Delete(cn.URLS.EachBudget, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.HandleDeleteBudget)
	app.Delete(cn.URLS.DeleteCapitalExpense, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.DeleteCapitalExpense)
	app.Delete(cn.URLS.DeleteEatoutExpense, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.DeleteEatoutExpense)
	app.Delete(cn.URLS.DeleteEntertainmentExpense, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.DeleteEntertainmentExpense)

	app.Patch(cn.URLS.UpdateBudget, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.PatchBudget)
	app.Patch(cn.URLS.UpdateCapitalExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.PatchCapitalExpenses)
	app.Patch(cn.URLS.UpdateEatoutExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.PatchEatoutExpenses)
	app.Patch(cn.URLS.UpdateEntertainmentExpenses, m.IsLoggedIn(hc.AuthHandler.TokenManager), hc.FinanceHandler.PatchEntertainmentExpenses)
}

func Setup(app *fiber.App, hc *h.HandlersManager) {
	AuthSetup(app, hc)
	FinanceSetup(app, hc)
	FitnessSetup(app, hc)

	// GPT
	// app.Post(cn.URLS.GptApi, h.HandlePostGptApi)

	// Websockets
	app.Get("/ws", websocket.New(func(wsc *websocket.Conn) {
		defer wsc.Close()
		client := openai.NewClient(cn.EnvVars.GptApiKey)

		for {
			// Handle WebSocket messages here
			messageType, msg, err := wsc.ReadMessage()
			if err != nil {
				// Handle the error
				break
			}
			prompt := string(msg)

			// Prepare req for GTP API
			req := openai.ChatCompletionRequest{
				Model:     openai.GPT3Dot5Turbo,
				MaxTokens: 500,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt + "In maximum 100 words.",
					},
				},
				Stream: true,
			}

			// Returns stream
			stream, err := client.CreateChatCompletionStream(context.Background(), req)
			if err != nil {
				log.Printf("ChatCompletionStream error: %v\n", err)
				return
			}
			defer stream.Close()

			// Send the stream to frontend
			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					log.Println("Stream finished")
					break
				}
				if err != nil {
					log.Printf("\nStream error: %v\n", err)
					break
				}

				// Send response to WebSocket client
				if messageType == websocket.TextMessage {
					if err := wsc.WriteMessage(websocket.TextMessage, []byte(response.Choices[0].Delta.Content)); err != nil {
						log.Printf("WebSocket send error: %v\n", err)
						break
					}
				}
			}
		}
	}))
}
