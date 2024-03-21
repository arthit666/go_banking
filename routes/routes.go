package routes

import (
	"os"

	"github.com/arthit666/make_app/account"

	"github.com/arthit666/make_app/middleware"
	"github.com/arthit666/make_app/pocket"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"gorm.io/gorm"

	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/arthit666/make_app/docs"
)

func RegRoute(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	a := account.New(db)
	app.Post("/login/", a.Login)
	app.Get("/refresh/", a.RefreshAccessToken)
	app.Post("/accounts/", a.CreateAccount)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.ErrUnauthorized
		},
	}))

	app.Use(middleware.ExtractUserFromJWT)

	app.Get("/accounts/", a.GetAllAccounts)
	app.Get("/accounts/:id", a.GetAccountById)
	app.Post("/accounts/transfer", a.Transfer)

	p := pocket.New(db)
	app.Post("/pockets/", p.CreatePocket)
	app.Get("/pockets/", p.GetAllPockets)
	app.Get("/pockets/:id", p.GetPocketById)
	app.Put("/pockets/:id", p.UpdatePocket)
	app.Delete("/pockets/:id", p.DeletePocket)
	app.Post("/pockets/transfer", p.Transfer)

	return app
}
