package handlers_test

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	ha "github.com/saeidalz13/LifeStyle2/lifeStyleBack/handlers"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
)

const (
	randomString32Chars = "c808cd4bc8639e0808216f5180277189"
	dbTestUrl           = "postgresql://root:testpassword@0.0.0.0:5432/lfdb?sslmode=disable"
	migrationDir        = "file:../db/migration"

	validEmail       = "test@gmail.com"
	validPassword    = "SomePassword13"
	validContentType = "application/json"

	invalidToken = ""

	testRequestTimeout = 5000
	anotherValidEmail  = "anotheremail@gmail.com"
	nonExistentEmail   = "emaildoesnotexist@gmail.com"
	testDbDriver       = "postgres"
)

var (
	validPasetoToken string
	// validBudgetId    string
	testMock         sqlmock.Sqlmock

	testHandlersManager *ha.HandlersManager
	app                 = fiber.New()
)

type TestCase[T any] struct {
	expectedStatus int
	name           string
	url            string
	contentType    string
	token          string
	reqPayload     T
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Entry point of all the tests in handlers package
func TestMain(m *testing.M) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	testMock = mock

	// Env vars
	config.EnvVars = &config.DotEnvVars{
		DbTestUrl: dbTestUrl,
		PasetoKey: randomString32Chars,
		DevStage:  "dev",
	}

	// config.OAuthConfigFitFin = googleOAuthConfig

	// Paseto token
	tokenManager, err := token.NewPasetoMaker(config.EnvVars.PasetoKey)
	if err != nil {
		log.Fatalln("Failed to extract Paseto Key!")
	}
	t, err := tokenManager.CreateToken(validEmail, config.PasetoTokenDuration)
	checkErr(err)
	validPasetoToken = t

	hm := ha.NewHandlersConfig(
		&ha.AuthHandlersManager{Db: db, TokenManager: tokenManager},
		&ha.FinanceHandlersManager{Db: db},
		&ha.FitnessHandlersManager{Db: db},
	)
	testHandlersManager = hm

	os.Exit(m.Run())
}

// googleOAuthConfig := &oauth2.Config{
// 	RedirectURL:  envVarsTest.GRedirUrl,
// 	ClientID:     envVarsTest.GClientId,
// 	ClientSecret: envVarsTest.GClientSec,
// 	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 	Endpoint:     google.Endpoint,
// }

// func mustInsertFitnessExercises(db *sql.DB) error {
// 	// Add move types
// 	ctx, cancel := context.WithTimeout(context.Background(), config.CONTEXT_TIMEOUT)
// 	defer cancel()
// 	q := sqlc.New(db)
// 	for _, moveType := range config.MOVE_TYPES_SLICE {
// 		if err := q.AddMoveType(ctx, moveType); err != nil {
// 			return err
// 		}
// 	}

// 	for i, moves := range config.Exercises {
// 		mType, err := q.FetchMoveTypeId(ctx, config.MOVE_TYPES_SLICE[i])
// 		if err != nil {
// 			return fmt.Errorf(config.ErrsFitFin.InvalidMoveType)
// 		}
// 		for _, move := range moves {
// 			if err := q.AddMoves(ctx, sqlc.AddMovesParams{
// 				MoveName:   move,
// 				MoveTypeID: mType.MoveTypeID,
// 			}); err != nil {
// 				return fmt.Errorf(config.ErrsFitFin.MoveInsertion)
// 			}
// 		}
// 	}
// 	return nil
// }

// func mustMigrate(db *sql.DB, migrationDir string) {
// 	driver, err := postgres.WithInstance(db, &postgres.Config{})
// 	checkErr(err)

// 	m, err := migrate.NewWithDatabaseInstance(
// 		migrationDir,
// 		"lfdb", // databaseName (random string for logging)
// 		driver, // Driver
// 	)

// 	checkErr(err)

// 	if err = m.Up(); err != nil {
// 		if err.Error() == config.ErrsFitFin.NoChangeMigration {
// 			// If no new migration, just start the server
// 			return
// 		}
// 		log.Fatalln(err)
// 	}
// }
