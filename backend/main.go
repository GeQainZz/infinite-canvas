package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"infinite-canvas-server/config"
	"infinite-canvas-server/handler"
	"infinite-canvas-server/model"
	"infinite-canvas-server/repository"
	"infinite-canvas-server/router"
	"infinite-canvas-server/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.DBDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.CreditAccount{},
		&model.CreditTransaction{},
		&model.CreditPricing{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	tenantRepo := repository.NewTenantRepo(db)
	creditRepo := repository.NewCreditRepo(db)

	authService := service.NewAuthService(cfg, userRepo, tenantRepo)
	userService := service.NewUserService(userRepo)
	creditService := service.NewCreditService(creditRepo)

	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)
	creditHandler := handler.NewCreditHandler(creditService, creditRepo)

	r := gin.Default()
	router.Setup(r, authService, authHandler, userHandler, creditHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
