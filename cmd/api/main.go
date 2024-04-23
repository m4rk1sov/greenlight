package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/joho/godotenv"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/jsonlog"
	"greenlight.m4rk1sov.github.com/internal/mailer"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// application version
const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
// config struct to hold all the configuration settings for our application.
// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	// Add a new limiter struct containing fields for the requests-per-second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	// altogether.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	// Update the config struct to hold the SMTP server settings.
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware
// Add a models field to hold our new Models struct.

// Change the logger field to have the type *jsonlog.Logger, instead of
// *log.Logger.

// Update the application struct to hold a new Mailer instance.

// Include a sync.WaitGroup in the application struct. The zero-value for a
// sync.WaitGroup type is a valid, useable, sync.WaitGroup with a 'counter' value of 0,
// so we don't need to do anything else to initialize it before we can use it.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// load .env file from given path
	// we keep it empty it will load .env from current directory
	envLoadErr := godotenv.Load(".env")

	if envLoadErr != nil {
		log.Fatalf("Error loading .env file")
	}
	// getting env variables
	dbHost := os.Getenv("ALMAS_DB_DSN")

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development"
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Use the value of the GREENLIGHT_DB_DSN environment variable as the default value
	// for our db-dsn command-line flag.

	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbHost, "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Read the SMTP server configuration settings into the config struct, using the
	// Mailtrap settings as the default values. IMPORTANT: If you're following along,
	// make sure to replace the default values for smtp-username and smtp-password
	// with your own Mailtrap credentials.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "eb444b7fc1bd66", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "397d341a1b10d0", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@almasmagzumov.mail.ru>", "SMTP sender")

	flag.Parse()

	////A new logger which writes messages to the standard out stream, current date and time.
	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		// Use the PrintFatal() method to write a log entry containing the error at the
		// FATAL level and exit. We have no additional properties to include in the log
		// entry, so we pass nil as the second parameter.
		logger.PrintFatal(err, nil)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	//// Also log a message to say that the connection pool has been successfully
	//// established.
	//logger.Printf("database connection pool established")

	// Likewise use the PrintInfo() method to write a message at the INFO level.
	logger.PrintInfo("database connection pool established", nil)

	// Declare an instance of the application struct, containing the config struct
	// Use the data.NewModels() function to initialize a Models struct, passing in the
	// connection pool as a parameter.
	// Initialize a new Mailer instance using the settings from the command line
	// flags, and add it to the application struct.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	// THIS GOT MOVED TO server.go

	//// Declare an HTTP server with some sensible timeout settings
	//srv := &http.Server{
	//	Addr:         fmt.Sprintf(":%d", cfg.port),
	//	Handler:      app.routes(),
	//	IdleTimeout:  time.Minute,
	//	ReadTimeout:  10 * time.Second,
	//	WriteTimeout: 30 * time.Second,
	//}
	////// Start the HTTP server.
	////logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	//
	//// Again, we use the PrintInfo() method to write a "starting server" message at the
	//// INFO level. But this time we pass a map containing additional properties (the
	//// operating environment and server address) as the final parameter.
	//logger.PrintInfo("starting server", map[string]string{
	//	"addr": srv.Addr,
	//	"env":  cfg.env,
	//})
	//
	//// Because the err variable is now already declared in the code above, we need
	//// to use the = operator here, instead of the := operator.
	//err = srv.ListenAndServe()
	//// Use the PrintFatal() method to log the error and exit.
	//logger.PrintFatal(err, nil)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
