package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/controller"
	"backend/db"
	"backend/handlers"
	"backend/repository"
	"backend/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg.Database)
	defer database.Close()

	clinicRepo := repository.NewClinicRepository(database)
	slotRepo := repository.NewSlotRepository(database)
	resRepo := repository.NewReservationRepository(database)
	paymentRepo := repository.NewPaymentRepository(database)

	bookingCtrl := controller.NewBookingController(database, slotRepo, resRepo, paymentRepo)
	expiryCtrl := controller.NewExpiryController(database, slotRepo, resRepo, paymentRepo, time.Duration(cfg.Booking.HoldExpirySecs)*time.Second)

	clinicHandler := handlers.NewClinicHandler(clinicRepo)
	slotHandler := handlers.NewSlotHandler(slotRepo)
	bookingHandler := handlers.NewBookingHandler(bookingCtrl)
	h := handlers.NewHandlers(clinicHandler, slotHandler, bookingHandler)

	expWorker := worker.NewExpiryWorker(database, expiryCtrl, time.Duration(cfg.Booking.ExpiryPollIntervalSecs)*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go expWorker.Run(ctx)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/clinics", h.Clinic.ListClinics)
		v1.GET("/clinics/:clinicId", h.Clinic.GetClinic)
		v1.GET("/clinics/:clinicId/slots", h.Slot.ListSlots)

		v1.POST("/slots/:slotId/hold", h.Booking.HoldSlot)
		v1.POST("/slots/:slotId/book", h.Booking.BookSlot)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel() // stop worker

	ctxShutDown, cancelShutDown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutDown()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
