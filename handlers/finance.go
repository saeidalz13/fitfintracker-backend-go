package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	sqlc "github.com/saeidalz13/LifeStyle2/lifeStyleBack/db/sqlc"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/models"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/utils"
)

type FinanceHandlersManager struct {
	Db           *sql.DB
	tokenManager token.TokenManager
}

func (f *FinanceHandlersManager) HandleGetAllBudgets(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	limitQry := ftx.Query("limit", "2")
	offsetQry := ftx.Query("offset", "0")
	log.Println("limit:" + limitQry + "offset" + offsetQry)
	convertedInts, err := utils.ConvertStringToInt64(limitQry, offsetQry)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch data"})
	}
	limit, offset := int32(convertedInts[0]), int32(convertedInts[1])

	budgets, err := q.SelectAllBudgets(ctx, sqlc.SelectAllBudgetsParams{UserID: user.ID, Offset: offset, Limit: limit})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Could not find the budgets!"})
	}

	numBudgets, err := q.CountBudgets(ctx, user.ID)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Could not find the number of budgets!"})
	}
	allBudgets := &models.OutgoingAllBudgets{
		Budgets:    budgets,
		NumBudgets: numBudgets,
	}

	log.Println("Budgets were found. Sending to front end...")
	return ftx.Status(fiber.StatusOK).JSON(allBudgets)

}

func (f *FinanceHandlersManager) HandleGetSingleBudget(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}
	idString := ftx.Params("id")
	budgetId, err := utils.FetchIntOfParamId(idString)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch the budget ID"})
	}

	singleBudget := sqlc.SelectSingleBudgetParams{
		BudgetID: int64(budgetId),
		UserID:   user.ID,
	}
	budget, err := q.SelectSingleBudget(ctx, singleBudget)
	if err != nil {
		if err.Error() == cn.SqlErrs.NoRows {
			return ftx.Status(fiber.StatusNotFound).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: err.Error()})
		}
		log.Println("Error fetching single budget", err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch budget from database"})
	}

	return ftx.Status(fiber.StatusOK).JSON(budget)
}

func (f *FinanceHandlersManager) HandleGetSingleBalance(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	idString := ftx.Params("id")
	budgetId, err := utils.FetchIntOfParamId(idString)
	if err != nil {
		log.Println("Conversion error:", err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch the budget ID"})
	}

	balance, err := q.SelectBalance(ctx, sqlc.SelectBalanceParams{UserID: user.ID, BudgetID: int64(budgetId)})
	if err != nil {
		if err.Error() == cn.SqlErrs.NoRows {
			return ftx.Status(fiber.StatusNotFound).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "No balance is available with this budget ID"})
		}
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch the balance!"})
	}

	log.Printf("Balance Found: %#v", balance)
	return ftx.Status(fiber.StatusOK).JSON(balance)
}

func (f *FinanceHandlersManager) HandleGetCapitalExpenses(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	userId, budgetID, limit, offset, searchString, err := prepareInputArgsForGetExpenses(ctx, ftx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	// redisKey := fmt.Sprint(userId) + ":" + fmt.Sprint(budgetID) + ":" + fmt.Sprint(limit) + ":" + fmt.Sprint(offset) + ":" + searchString
	// var cachedExpenses []sqlc.FetchAllCapitalExpensesRow
	// value, err := rUtils.GetRedisStruct(f.Rdb, redisKey, cachedExpenses)
	// if err != nil {
	// 	if err.Error() == "EOF" {
	// 		log.Printf("no cached capital expenses: %s\n", redisKey)
	// 	}
	// 	log.Printf("error get this key from redis %s: %v\n%", redisKey, err)
	// } else {
	// 	log.Printf("%+v", value)
	// }

	capitalExpenses, err := q.FetchAllCapitalExpenses(ctx, sqlc.FetchAllCapitalExpensesParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Limit:       limit,
		Offset:      offset,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to fetch capital expenses"})
	}

	capitalRowCountTotal, err := q.FetchTotalRowCountCapital(ctx, sqlc.FetchTotalRowCountCapitalParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to row counts and total of capital"})
	}

	log.Printf("%#v\n", f.Db.Stats())
	return ftx.Status(fiber.StatusOK).JSON(models.CapitalExpensesResponse{
		ExpenseType:          "capital",
		CapitalExpenses:      capitalExpenses,
		CapitalTotalRowCount: capitalRowCountTotal,
	})
}

func (f *FinanceHandlersManager) HandleGetEatoutExpenses(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	userId, budgetID, limit, offset, searchString, err := prepareInputArgsForGetExpenses(ctx, ftx, q, f.tokenManager)

	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	eatoutExpenses, err := q.FetchAllEatoutExpenses(ctx, sqlc.FetchAllEatoutExpensesParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Limit:       limit,
		Offset:      offset,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
	}
	eatoutRowCountTotal, err := q.FetchTotalRowCountEatout(ctx, sqlc.FetchTotalRowCountEatoutParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
	}
	return ftx.Status(fiber.StatusOK).JSON(models.EatoutExpensesResponse{
		ExpenseType:         "eatout",
		EatoutExpenses:      eatoutExpenses,
		EatoutTotalRowCount: eatoutRowCountTotal,
	})
}

func (f *FinanceHandlersManager) HandleGetEntertainmentExpenses(ftx *fiber.Ctx) error {
	q := sqlc.New(f.Db)
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()

	userId, budgetID, limit, offset, searchString, err := prepareInputArgsForGetExpenses(ctx, ftx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	entertainmentExpenses, err := q.FetchAllEntertainmentExpenses(ctx, sqlc.FetchAllEntertainmentExpensesParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Limit:       limit,
		Offset:      offset,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
	}

	entertainmentRowCountTotal, err := q.FetchTotalRowCountEntertainment(ctx, sqlc.FetchTotalRowCountEntertainmentParams{
		UserID:      userId,
		BudgetID:    budgetID,
		Description: searchString,
	})
	if err != nil {
		log.Println(err)
	}
	return ftx.Status(fiber.StatusOK).JSON(models.EntertainmentExpensesResponse{
		ExpenseType:                "entertainment",
		EntertainmentExpenses:      entertainmentExpenses,
		EntertainmentTotalRowCount: entertainmentRowCountTotal,
	})
}

func (f *FinanceHandlersManager) HandlePostNewBudget(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if !isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	var newBudget sqlc.CreateBudgetParams
	if err := ftx.BodyParser(&newBudget); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	operationBudget := sqlc.CreateBudgetParams(newBudget)
	operationBudget.UserID = user.ID
	op := sqlc.NewQWithTx(f.Db)
	result, err := op.CreateBudgetBalance(ctx, operationBudget)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "unique_combination_constraint") {
			return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Budget NAME already exists. Choose another one please!"})
		}
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to add the budget!"})
	}
	log.Printf("Budget Creation Success: %v", result)
	return ftx.Status(fiber.StatusCreated).JSON(&cn.ApiRes{ResType: cn.ResTypes.Success, Msg: "Budget Created Successfully!"})
}

func (f *FinanceHandlersManager) HandlePostExpenses(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	// JSON Parsing Stage
	var newExpense models.ExpenseReq
	if err := ftx.BodyParser(&newExpense); err != nil {
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	// Trim and lower case description
	utils.NormalizeInput(&newExpense.ExpenseDesc)

	q2 := sqlc.NewQWithTx(f.Db)
	updatedBalance, err := q2.AddExpenseUpdateBalance(ctx, sqlc.AddExpenseUpdateBalanceTx{
		BudgetID:    newExpense.BudgetID,
		UserID:      user.ID,
		Expenses:    newExpense.ExpenseAmount,
		ExpenseType: newExpense.ExpenseType,
		Description: newExpense.ExpenseDesc,
	})

	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to complete the expense adding transaction"})
	}

	log.Println(newExpense.ExpenseType + " expense was added for " + user.Email)
	return ftx.Status(fiber.StatusOK).JSON(updatedBalance)
}

// DELETE
func (f *FinanceHandlersManager) HandleDeleteBudget(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	// Extract Budget ID
	idString := ftx.Params("id")
	budgetId, err := utils.FetchIntOfParamId(idString)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ExtractUrlParam})
	}

	if err = q.DeleteBudget(ctx, sqlc.DeleteBudgetParams{
		BudgetID: int64(budgetId),
		UserID:   user.ID,
	}); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to delete the budget!"})
	}

	log.Printf("Budget ID: %v -> Deleted", budgetId)
	return ftx.SendStatus(fiber.StatusNoContent)
	// return ftx.Status(fiber.StatusAccepted).JSON(&cn.ApiRes{ResType: cn.ResTypes.Success, Msg: "Budget was deleted successfully!"})
}

func (f *FinanceHandlersManager) DeleteCapitalExpense(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	qwtx := sqlc.NewQWithTx(f.Db)
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	var deleteArgTx sqlc.DeleteSingleCapitalExpenseParams
	if err := ftx.BodyParser(&deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}
	deleteArgTx.UserID = user.ID

	if err = qwtx.DeleteCapitalExpenseBalance(ctx, &deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to delete the capital expense"})
	}
	return ftx.SendStatus(fiber.StatusNoContent)
}

func (f *FinanceHandlersManager) DeleteEatoutExpense(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	qwtx := sqlc.NewQWithTx(f.Db)
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}
	var deleteArgTx sqlc.DeleteSingleEatoutExpenseParams
	if err := ftx.BodyParser(&deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}
	deleteArgTx.UserID = user.ID

	if err = qwtx.DeleteEatoutExpenseBalance(ctx, &deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to delete the eatout expense"})
	}
	return ftx.SendStatus(fiber.StatusNoContent)
}

func (f *FinanceHandlersManager) DeleteEntertainmentExpense(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	qwtx := sqlc.NewQWithTx(f.Db)
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}
	var deleteArgTx sqlc.DeleteSingleEntertainmentExpenseParams
	if err := ftx.BodyParser(&deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}
	deleteArgTx.UserID = user.ID

	if err = qwtx.DeleteEntertainmentExpenseBalance(ctx, &deleteArgTx); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to delete the entertainment expense"})
	}
	return ftx.SendStatus(fiber.StatusNoContent)
}

/*
PATCH Section
*/
func (f *FinanceHandlersManager) PatchBudget(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	// Prepare the input
	var updateBudgetBalance sqlc.UpdateBudgetParams
	if err := ftx.BodyParser(&updateBudgetBalance); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}
	updateBudgetBalance.UserID = user.ID
	log.Printf("Incoming: %#v", updateBudgetBalance)

	// Do the transaction
	q2 := sqlc.NewQWithTx(f.Db)
	updatedBudget, updatedBalance, err := q2.UpdateBudgetBalance(ctx, updateBudgetBalance)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to update the budget"})
	}

	log.Println("Budget and balance updated successfully!")
	return ftx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"updated_budget":  updatedBudget,
		"updated_balance": updatedBalance,
	})
}

func (f *FinanceHandlersManager) PatchCapitalExpenses(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	var updateCapitalExpensesInfo sqlc.IncomingUpdateCapitalExpenses
	if err := ftx.BodyParser(&updateCapitalExpensesInfo); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	qwtx := sqlc.NewQWithTx(f.Db)
	err = qwtx.UpdateExpensesBalanceCapital(ctx, &sqlc.ExpenseBalanceCapital{
		IncomingUpdateCapitalExpenses: updateCapitalExpensesInfo,
		UserId:                        user.ID,
	})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to update capital expenses"})
	}

	return ftx.SendStatus(fiber.StatusOK)
}

func (f *FinanceHandlersManager) PatchEatoutExpenses(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	var updateEatoutExpensesInfo sqlc.IncomingUpdateEatoutExpenses
	if err := ftx.BodyParser(&updateEatoutExpensesInfo); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	qwtx := sqlc.NewQWithTx(f.Db)
	err = qwtx.UpdateExpensesBalanceEatout(ctx, &sqlc.ExpenseBalanceEatout{
		IncomingUpdateEatoutExpenses: updateEatoutExpensesInfo,
		UserId:                       user.ID,
	})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to update eatout expenses"})
	}

	return ftx.SendStatus(fiber.StatusOK)
}

func (f *FinanceHandlersManager) PatchEntertainmentExpenses(ftx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), cn.CONTEXT_TIMEOUT)
	defer cancel()
	q := sqlc.New(f.Db)

	if isContentTypeJson(ftx) {
		return ftx.Status(fiber.StatusBadRequest).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ContentType})
	}

	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, f.tokenManager)
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusUnauthorized).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.UserValidation})
	}

	var updateEntertainmentExpensesInfo sqlc.IncomingUpdateEntertainmentExpenses
	if err := ftx.BodyParser(&updateEntertainmentExpensesInfo); err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: cn.ErrsFitFin.ParseJSON})
	}

	qwtx := sqlc.NewQWithTx(f.Db)
	err = qwtx.UpdateExpensesBalanceEntertainment(ctx, &sqlc.ExpenseBalanceEntertainment{
		IncomingUpdateEntertainmentExpenses: updateEntertainmentExpensesInfo,
		UserId:                              user.ID,
	})
	if err != nil {
		log.Println(err)
		return ftx.Status(fiber.StatusInternalServerError).JSON(&cn.ApiRes{ResType: cn.ResTypes.Err, Msg: "Failed to update entertainment expenses"})
	}

	return ftx.SendStatus(fiber.StatusOK)
}

// HELPERS
func prepareInputArgsForGetExpenses(ctx context.Context, ftx *fiber.Ctx, q *sqlc.Queries, tokenManager token.TokenManager) (int64, int64, int32, int32, string, error) {
	user, err := fetchUserFromToken(ftx.Cookies(cn.PasetoCookieName), ctx, q, tokenManager)
	if err != nil {
		return -1, -1, -1, -1, "", errors.New("failed to validate user")
	}

	idString := ftx.Params("id")
	budgetId, err := utils.FetchIntOfParamId(idString)
	if err != nil {
		return -1, -1, -1, -1, "", errors.New("failed to fetch budget ID from url params")
	}

	limitQry := ftx.Query("limit", "10")
	offsetQry := ftx.Query("offset", "0")
	convertedInts, err := utils.ConvertStringToInt64(limitQry, offsetQry)

	searchString := ftx.Query("search", "")
	searchString = utils.PrepareSearchString(searchString)

	if err != nil {
		return -1, -1, -1, -1, "", errors.New("failed to fetch budget ID from req json body")
	}
	limit, offset := int32(convertedInts[0]), int32(convertedInts[1])
	return user.ID, int64(budgetId), limit, offset, searchString, nil
}
