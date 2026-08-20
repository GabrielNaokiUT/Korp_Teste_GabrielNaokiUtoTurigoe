package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type NotaFiscal struct {
	ID        int    `json:"id"`
	Numeracao int    `json:"numeracao"`
	Status    string `json:"status"` // Aberta ou Fechada[cite: 1]
}

var db *sql.DB

func main() {
	var err error
	connStr := "user=postgres password=root dbname=faturamento sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/notas", CriarNota).Methods("POST")
	r.HandleFunc("/notas/{id}/imprimir", ImprimirNota).Methods("POST")

	log.Println("Serviço de Faturamento rodando na porta 8082...")
	log.Fatal(http.ListenAndServe(":8082", r))
}

// Cria a nota com status inicial Aberta[cite: 1]
func CriarNota(w http.ResponseWriter, r *http.Request) {
	var nota NotaFiscal
	json.NewDecoder(r.Body).Decode(&nota)
	nota.Status = "Aberta" // Status inicial obrigatório[cite: 1]

	err := db.QueryRow("INSERT INTO notas_fiscais (numeracao, status) VALUES ($1, $2) RETURNING id", nota.Numeracao, nota.Status).Scan(&nota.ID)
	if err != nil {
		http.Error(w, "Erro no banco", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nota)
}

// Imprime a nota e fala com o serviço de estoque[cite: 1]
func ImprimirNota(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var status string
	db.QueryRow("SELECT status FROM notas_fiscais WHERE id = $1", id).Scan(&status)
	if status != "Aberta" {
		http.Error(w, "Nota não está aberta", http.StatusBadRequest)
		return
	}

	itemRequest := map[string]int{"quantidade": 2} // Simulação de 2 itens para baixar[cite: 1]
	jsonReq, _ := json.Marshal(itemRequest)

	resp, err := http.Post("http://localhost:8081/produtos/1/baixar", "application/json", bytes.NewBuffer(jsonReq))

	// Tratamento de Falha: Se der erro, aborta a operação e mantém como Aberta
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Erro ao contatar serviço de estoque:", err)
		http.Error(w, "Falha na comunicação com o estoque. Nota fiscal não impressa.", http.StatusBadGateway)
		return
	}

	// Atualiza para Fechada se a comunicação for um sucesso[cite: 1]
	_, err = db.Exec("UPDATE notas_fiscais SET status = 'Fechada' WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Erro ao atualizar status", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Nota impressa e fechada com sucesso!")
}
