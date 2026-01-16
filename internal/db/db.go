package db // Package db provides database connection helpers.

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sales-management-api/internal/config"
	"sales-management-api/internal/models"
)

func Connect(cfg config.Config) (*gorm.DB, error) { // Connect opens the DB and prepares it.
	dsn := fmt.Sprintf( // Build the Postgres DSN from config.
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta", // DSN template with timezone.
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, // DSN values from config.
	) // Finish DSN formatting.

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{ // Open GORM with Postgres.
		Logger: logger.Default.LogMode(logger.Warn), // Set logging level to warn.
	}) // Close GORM config.
	if err != nil { // Check for connection errors.
		return nil, err // Return error if open fails.
	} // End error check.

	sqlDB, err := db.DB() // Get the underlying *sql.DB.
	if err != nil {       // Handle errors from db.DB().
		return nil, err // Return error to caller.
	} // End error check.
	if err := sqlDB.Ping(); err != nil { // Verify the connection is alive.
		return nil, err // Return error if ping fails.
	} // End ping check.

	// auto migrate tables
	if err := db.AutoMigrate( // Create or update schema for models.
		&models.User{},     // Users table.
		&models.Product{},  // Products table.
		&models.Sale{},     // Sales table.
		&models.SaleItem{}, // Sale items table.
	); err != nil { // Check for migration errors.
		return nil, err // Return error if migrations fail.
	} // End migration block.

	if err := seedAdmin(db, cfg.SeedAdminUser, cfg.SeedAdminPass); err != nil { // Seed admin if configured.
		return nil, err // Return error if seeding fails.
	} // End seed block.

	return db, nil // Return the ready DB handle.
}

func seedAdmin(db *gorm.DB, username, password string) error { // Create admin user if missing.
	if username == "" || password == "" { // Skip if credentials are empty.
		return nil // No-op when not configured.
	} // End empty check.

	var count int64                                                        // Holds number of matching users.
	db.Model(&models.User{}).Where("username = ?", username).Count(&count) // Count existing admin user.
	if count > 0 {                                                         // If user exists, skip seeding.
		return nil // Nothing to do.
	} // End exists check.

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Hash the password.
	if err != nil {                                                                // Handle hash error.
		return err // Return hash error.
	} // End hash check.

	u := models.User{ // Build the admin user record.
		Username:     username,         // Set username.
		PasswordHash: string(hash),     // Store hashed password.
		Role:         models.RoleAdmin, // Assign admin role.
	} // End user struct literal.
	return db.Create(&u).Error // Insert admin user.
}
