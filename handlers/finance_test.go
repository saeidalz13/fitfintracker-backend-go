// /*
// This file will test all the handlers together since
// it is necessary for fitness and finance module to
// have a valid user in the test db.
// */

package handlers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
// 	sqlc "github.com/saeidalz13/LifeStyle2/lifeStyleBack/db/sqlc"
// 	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/models"
// 	// "github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
// 	// "github.com/saeidalz13/LifeStyle2/lifeStyleBack/utils"
// )


// func TestHandlePostNewBudget(t *testing.T) {
// 	app.Post(cn.URLS.PostNewBudget, testHandlersManager.FinanceHandler.HandlePostNewBudget)

// 	tests := []TestCase[sqlc.CreateBudgetParams]{
// 		{
// 			name:           "create budget valid token valid params",
// 			expectedStatus: http.StatusCreated,
// 			url:            cn.URLS.PostNewBudget,
// 			contentType:    validContentType,
// 			token:          validPasetoToken,
// 			reqPayload: sqlc.CreateBudgetParams{
// 				BudgetName:    "random",
// 				StartDate:     time.Now(),
// 				EndDate:       time.Now().Add(time.Hour),
// 				Savings:       "100000",
// 				Capital:       "1000",
// 				Eatout:        "1000",
// 				Entertainment: "2000",
// 			},
// 		},
// 		{
// 			name:           "not create budget same budget name",
// 			expectedStatus: http.StatusInternalServerError,
// 			url:            cn.URLS.PostNewBudget,
// 			contentType:    validContentType,
// 			token:          validPasetoToken,
// 			reqPayload: sqlc.CreateBudgetParams{
// 				BudgetName:    "random",
// 				StartDate:     time.Now(),
// 				EndDate:       time.Now().Add(time.Hour),
// 				Savings:       "100000",
// 				Capital:       "1000",
// 				Eatout:        "1000",
// 				Entertainment: "2000",
// 			},
// 		},
// 		{
// 			name:           "not create budget invalid money params",
// 			expectedStatus: http.StatusInternalServerError,
// 			url:            cn.URLS.PostNewBudget,
// 			contentType:    validContentType,
// 			token:          validPasetoToken,
// 			reqPayload: sqlc.CreateBudgetParams{
// 				BudgetName:    "new_random",
// 				StartDate:     time.Now(),
// 				EndDate:       time.Now().Add(time.Hour),
// 				Savings:       "gs",
// 				Capital:       "g",
// 				Eatout:        "100grd50",
// 				Entertainment: "dg",
// 			},
// 		},
// 		{
// 			name:           "not create budget invalid token",
// 			expectedStatus: http.StatusUnauthorized,
// 			url:            cn.URLS.PostNewBudget,
// 			contentType:    validContentType,
// 			token:          invalidToken,
// 			reqPayload: sqlc.CreateBudgetParams{
// 				BudgetName:    "new random",
// 				StartDate:     time.Now(),
// 				EndDate:       time.Now().Add(time.Hour),
// 				Savings:       "100000",
// 				Capital:       "1000",
// 				Eatout:        "1000",
// 				Entertainment: "2000",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			body, err := json.Marshal(test.reqPayload)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req := httptest.NewRequest(http.MethodGet, test.url, bytes.NewBuffer(body))
// 			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})
// 			req.Header.Set("Content-Type", test.contentType)

// 			resp, err := app.Test(req, 10)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if resp.StatusCode != test.expectedStatus {
// 				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
// 			}
// 		})
// 	}
// }

// func TestHandleGetAllBudgets(t *testing.T) {
// 	tests := []TestCase[any]{
// 		{
// 			expectedStatus: http.StatusOK,
// 			name:           "get all budgets valid token",
// 			url:            cn.URLS.ShowBudgets,
// 			token:          validPasetoToken,
// 		},
// 		{
// 			expectedStatus: http.StatusUnauthorized,
// 			name:           "no budgets invalid token",
// 			url:            cn.URLS.ShowBudgets,
// 			token:          invalidToken,
// 		},
// 	}

// 	for i, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodGet, test.url, nil)
// 			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

// 			resp, err := app.Test(req, 10)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if resp.StatusCode != test.expectedStatus {
// 				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
// 			}

// 			if i == 0 {
// 				var jsonResp models.OutgoingAllBudgets
// 				bodyBytes, err := io.ReadAll(resp.Body)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				if err = json.Unmarshal(bodyBytes, &jsonResp); err != nil {
// 					t.Fatal(err)
// 				}
// 				validBudgetId = fmt.Sprint(jsonResp.Budgets[0].BudgetID)
// 			}
// 		})
// 	}
// }

// func TestHandleGetSingleBudget(t *testing.T) {
// 	tests := []TestCase[string]{
// 		{
// 			expectedStatus: http.StatusOK,
// 			name:           "get single budget valid token",
// 			url:            cn.URLS.EachBudget,
// 			token:          validPasetoToken,
// 			reqPayload:     validBudgetId,
// 		},
// 		{
// 			expectedStatus: http.StatusNotFound,
// 			name:           "no budget invalid id",
// 			url:            cn.URLS.EachBudget,
// 			token:          validPasetoToken,
// 			reqPayload:     "-1",
// 		},
// 		{
// 			expectedStatus: http.StatusUnauthorized,
// 			name:           "no budget invalid token",
// 			url:            cn.URLS.EachBudget,
// 			token:          invalidToken,
// 			reqPayload:     validBudgetId,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodGet, strings.Replace(test.url, ":id", test.reqPayload, 1), nil)
// 			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

// 			resp, err := app.Test(req, 10)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if resp.StatusCode != test.expectedStatus {
// 				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
// 			}
// 		})
// 	}
// }

// func TestHandleGetSingleBalance(t *testing.T) {
// 	tests := []TestCase[string]{
// 		{
// 			name:           "get balance valid budget id",
// 			expectedStatus: http.StatusOK,
// 			url:            cn.URLS.EachBalance,
// 			token:          validPasetoToken,
// 			reqPayload:     validBudgetId,
// 		},
// 		{
// 			name:           "no balance invalid budget id",
// 			expectedStatus: http.StatusNotFound,
// 			url:            cn.URLS.EachBalance,
// 			token:          validPasetoToken,
// 			reqPayload:     "-1",
// 		},
// 		{
// 			name:           "no balance invalid token",
// 			expectedStatus: http.StatusUnauthorized,
// 			url:            cn.URLS.EachBalance,
// 			token:          invalidToken,
// 			reqPayload:     validBudgetId,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodGet, strings.Replace(test.url, ":id", test.reqPayload, 1), nil)
// 			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

// 			resp, err := app.Test(req, 10)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			if resp.StatusCode != test.expectedStatus {
// 				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
// 			}
// 		})
// 	}
// }
