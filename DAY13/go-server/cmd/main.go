/*
 * Digimon Service API
 *
 * 提供孵化數碼蛋與培育等數碼寶貝養成服務
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"database/sql"
	"fmt"
	_dietRepo "go-server/diet/repository/postgresql"
	_dietUsecase "go-server/diet/usecase"
	_digimonHandlerGRPCDelivery "go-server/digimon/delivery/grpc"
	_digmonRepo "go-server/digimon/repository/postgresql"
	_digimonUsecase "go-server/digimon/usecase"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("dotenv")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Fatal error config file: %v\n", err)
	}
}

func main() {
	logrus.Info("HTTP server started")

	grpcPort := viper.GetString("GRPC_PORT")
	dbHost := viper.GetString("DB_HOST")
	dbDatabase := viper.GetString("DB_DATABASE")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbDatabase),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatal(err)
	}

	digimonRepo := _digmonRepo.NewpostgresqlDigimonRepository(db)
	dietRepo := _dietRepo.NewPostgresqlDietRepository(db)

	digimonUsecase := _digimonUsecase.NewDigimonUsecase(digimonRepo)
	dietUsecase := _dietUsecase.NewDietUsecase(dietRepo)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logrus.Fatal(err)
	}
	s := grpc.NewServer()

	_digimonHandlerGRPCDelivery.NewDigimonHandler(s, digimonUsecase, dietUsecase)

	logrus.Fatal(s.Serve(lis))
}