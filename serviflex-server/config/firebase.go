package config

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
	ctx := context.Background()
	credsPath := "C:/Users/mathe/Downloads/Arquivo/config/pdsi-23063-firebase-adminsdk-fbsvc-8df5cdf4f5.json"
	if credsPath == "" {
		log.Fatal("Arquivo de credenciais do Firebase não encontrado")
	}
	fmt.Println("Usando credenciais:", credsPath)
	opt := option.WithCredentialsFile(credsPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: "pdsi-23063",
	}, opt)
	if err != nil {
		log.Fatalf("Erro ao inicializar o Firebase: %v", err)
	}
	App = app
}
