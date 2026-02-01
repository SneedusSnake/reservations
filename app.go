package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/cache"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/system"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/mysql"
	"github.com/SneedusSnake/Reservations/internal/adapters/driving/telegram"
	"github.com/SneedusSnake/Reservations/internal/application"
	"github.com/SneedusSnake/Reservations/internal/ports"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/goose/v3"
	_ "time/tzdata"
)

const (
	CLOCK = "clock"

	STORE_SUBJECTS     = "subjects_store"
	STORE_USERS        = "users_store"
	STORE_TG_USERS     = "tg_users_store"
	STORE_RESERVATIONS = "reservations_store"
	STORE_READ_RESERVATIONS = "reservations_read_store"

	SERVICE_SUBJECT = "subject_service"
	SERVICE_USER = "user_service"
	SERVICE_TELEGRAM_USER = "telegram_user_service"
	SERVICE_RESERVATION = "reservation_service"

	TELERAM_BOT = "telegram_bot"
)

type App struct{
	Config Config
	Log *log.Logger
	container map[string]any
}

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
		Token string `envconfig:"TELEGRAM_API_TOKEN"`
	}
	Clock string `envconfig:"CLOCK_DRIVER"`
	CacheClockPath string `envconfig:"CACHE_CLOCK_PATH"`
	PersistenceDriver string `envconfig:"PERSISTENCE_DRIVER"`
	MysqlConnection struct {
		ConnectionString string `envconfig:"MYSQL_CONNECTION"`
		Host string `envconfig:"MYSQL_HOST"`
		Port string `envconfig:"MYSQL_PORT"`
		Database string `envconfig:"MYSQL_DATABASE"`
		User string `envconfig:"MYSQL_USER"`
		Password string `envconfig:"MYSQL_PASSWORD"`
	}
	TimeZone string `envconfig:"TZ"`
}

func (app *App) Resolve(dependency string) any {
	result, ok := app.container[dependency]
	if !ok {
		app.Log.Fatalf("Could not resolve %s", dependency)
	}

	return result
}

func Bootstrap() *App {
	app := &App{
		Log: log.Default(),
		container: make(map[string]any),
	}
	app.loadConfig()
	app.registerDependencies()

	return app
}

func (app *App) usersStore() users.UsersRepository {
	return app.Resolve(STORE_USERS).(users.UsersRepository)
}

func (app *App) tgUsersStore() telegram.TelegramUsersRepository {
	return app.Resolve(STORE_TG_USERS).(telegram.TelegramUsersRepository)
}

func (app *App) subjectsStore() reservations.SubjectsRepository {
	return app.Resolve(STORE_SUBJECTS).(reservations.SubjectsRepository)
}

func (app *App) reservationsStore() reservations.ReservationsRepository {
	return app.Resolve(STORE_RESERVATIONS).(reservations.ReservationsRepository)
}

func (app *App) reservationsReadStore() reservations.ReservationsReadRepository {
	return app.Resolve(STORE_READ_RESERVATIONS).(reservations.ReservationsReadRepository)
}

func (app *App) loadConfig() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		app.Log.Print(err)
		panic(err)
	}
	app.Config = cfg
}

func (app *App) registerDependencies() {
	var clock ports.Clock

	clock = system.SystemClock{}
	if app.Config.Clock == "cache" {
		clock = cache.NewClock(app.Config.CacheClockPath)
	} 
	app.container[CLOCK] = clock

	app.registerStores()
	app.registerServices()

	tgBot := app.telegramBot()

	app.container[TELERAM_BOT] = tgBot
	app.registerTelegramBotHandlers()
}

func (app *App) registerStores() {
	var subjectsStore reservations.SubjectsRepository
	var reservationsStore reservations.ReservationsRepository
	var reservationsReadStore reservations.ReservationsReadRepository
	var usersStore users.UsersRepository
	var tgUsersStore telegram.TelegramUsersRepository

	subjectsStore = inmemory.NewSubjectsStore()
	usersStore = inmemory.NewUsersStore()
	tgUsersStore = inmemory.NewTelegramUsersStore(usersStore)
	reservationsStore = inmemory.NewReservationStore()
	reservationsReadStore = inmemory.NewReservationReadStore(
		reservationsStore.(*inmemory.ReservationsStore),
		usersStore.(*inmemory.UsersStore), 
		subjectsStore.(*inmemory.SubjectsStore),
	)

	if app.Config.PersistenceDriver == "mysql" {
		db := app.ConnectDB()
		app.Migrate(db)

		subjectsStore = mysql.NewSubjectsRepository(db)
		usersStore = mysql.NewUsersRepository(db)
		tgUsersStore = mysql.NewTelegramUsersRepository(db)
		reservationsStore = mysql.NewReservationsRepository(db)
		reservationsReadStore = mysql.NewReservationsReadRepository(db)
	}

	app.container[STORE_SUBJECTS] = subjectsStore
	app.container[STORE_USERS] = usersStore
	app.container[STORE_TG_USERS] = tgUsersStore
	app.container[STORE_RESERVATIONS] = reservationsStore
	app.container[STORE_READ_RESERVATIONS] = reservationsReadStore
}

func (app *App) registerServices() {
	subjectsStore := app.subjectsStore()
	reservationsStore := app.reservationsStore()
	reservationsReadStore := app.reservationsReadStore()
	usersStore := app.usersStore()
	tgUsersStore := app.tgUsersStore()

	reservationService := application.NewReservationService(
		subjectsStore,
		reservationsStore,
		reservationsReadStore,
		usersStore,
		app.Resolve(CLOCK).(ports.Clock),
	)
	subjectService := application.NewSubjectService(subjectsStore)
	userService := application.NewUserService(usersStore)
	tgUserService := telegram.NewTelegramUserService(tgUsersStore, *userService)

	app.container[SERVICE_RESERVATION] = reservationService
	app.container[SERVICE_SUBJECT] = subjectService
	app.container[SERVICE_USER] = userService
	app.container[SERVICE_TELEGRAM_USER] = tgUserService
}

func (app *App) telegramBot() *bot.Bot {
	var opts []bot.Option
	url := app.Config.TelegramApi.Host
	if url != "" {
		opts = append(opts, bot.WithServerURL(app.Config.TelegramApi.Host))
	}

	b, err := bot.New(app.Config.TelegramApi.Token, opts...)

	if err != nil {
		app.Error(err)
	}

	return b
}

func (app *App) registerTelegramBotHandlers() {
	b := app.Resolve(TELERAM_BOT).(*bot.Bot)
	adapter := telegram.NewAdapter(
		app.Resolve(SERVICE_SUBJECT).(*application.SubjectService),
		app.Resolve(SERVICE_RESERVATION).(*application.ReservationService),
		app.Resolve(SERVICE_USER).(*application.UserService),
		app.Resolve(SERVICE_TELEGRAM_USER).(*telegram.TelegramUserService),
		app.Resolve(CLOCK).(ports.Clock),
		app.Log,
	)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_tags", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, botHandlerFunc(adapter.ListSubjectsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tags", bot.MatchTypePrefix, botHandlerFunc(adapter.ListSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reserved", bot.MatchTypePrefix, botHandlerFunc(adapter.ActiveReservationsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reserve", bot.MatchTypePrefix, botHandlerFunc(adapter.CreateReservationHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remove", bot.MatchTypePrefix, botHandlerFunc(adapter.RemoveReservationHandler))
}

type UpdateHandler func(ctx context.Context, b *bot.Bot, update *models.Update) (string, error)

func botHandlerFunc(h UpdateHandler) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		text, err := h(ctx, b, update)

		if err != nil {
			log.Print(err)
			text = "An error occured"
		}

		if text != "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				MessageThreadID: update.Message.MessageThreadID,
				Text:   text,
			})
		}
	}
}

func (app *App) Error(err error) {
	app.Log.Fatal(err)
}

func (app *App) ConnectDB() *sql.DB{
	connectionString := app.Config.MysqlConnection.ConnectionString

	if  connectionString == "" {
		cfg := mysqlDriver.NewConfig()
		cfg.User = app.Config.MysqlConnection.User
		cfg.Passwd = app.Config.MysqlConnection.Password
		cfg.Net = "tcp"
		cfg.Addr = fmt.Sprintf("%s:%s", app.Config.MysqlConnection.Host, app.Config.MysqlConnection.Port)
		cfg.DBName = app.Config.MysqlConnection.Database
		cfg.Loc = time.Local
		cfg.ParseTime = true

		connectionString = cfg.FormatDSN()
	}

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		app.Error(err)
	}
	err = db.Ping()
	if err != nil {
		app.Error(err)
	}

	return db
}

func (app *App) Migrate(db *sql.DB) {
	goose.SetDialect("mysql")
	wd, err := os.Getwd()
	if err != nil {
		app.Error(err)
	}

	err = goose.Up(db, wd + "/migrations")
	if err != nil {
		app.Error(err)
	}
}
