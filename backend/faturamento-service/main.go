package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type ItemNota struct {
	ProdutoCodigo string `json:"produtoCodigo"`
	Quantidade    int    `json:"quantidade"`
}

type NotaFiscal struct {
	ID     int        `json:"id"`
	Numero int        `json:"numero"`
	Status string     `json:"status"` // "Aberta" ou "Fechada"
	Itens  []ItemNota `json:"itens"`
}

var (
	db    *sql.DB
	notas = make(map[int]*NotaFiscal)
	mu    sync.Mutex
)

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=korp sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("PostgreSQL indisponível. Usando modo em memória.")
		return
	}
	if err = db.Ping(); err != nil {
		log.Println("PostgreSQL indisponível. Usando modo em memória.")
		db = nil
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS notas (
		id SERIAL PRIMARY KEY,
		numero INT NOT NULL,
		status VARCHAR(20) NOT NULL,
		itens JSONB NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Println("Erro ao verificar tabela notas:", err)
	} else {
		log.Println("Conectado ao PostgreSQL com sucesso no Faturamento!")
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cadastrarNota(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var n NotaFiscal
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, `{"error":"JSON inválido"}`, http.StatusBadRequest)
		return
	}

	n.Status = "Aberta"

	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		itensJSON, _ := json.Marshal(n.Itens)
		err := db.QueryRow("INSERT INTO notas (numero, status, itens) VALUES ($1, $2, $3) RETURNING id",
			n.Numero, n.Status, string(itensJSON)).Scan(&n.ID)
		if err != nil {
			log.Println("Erro DB:", err)
			http.Error(w, `{"error":"Erro ao salvar nota no banco"}`, http.StatusInternalServerError)
			return
		}
	} else {
		n.ID = n.Numero
		notas[n.ID] = &n
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func listarNotas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		rows, err := db.Query("SELECT id, numero, status, itens FROM notas ORDER BY id DESC")
		if err == nil {
			defer rows.Close()
			var lista []NotaFiscal
			for rows.Next() {
				var n NotaFiscal
				var itensStr string
				rows.Scan(&n.ID, &n.Numero, &n.Status, &itensStr)
				json.Unmarshal([]byte(itensStr), &n.Itens)
				lista = append(lista, n)
			}
			json.NewEncoder(w).Encode(lista)
			return
		}
	}

	var lista []NotaFiscal
	for _, n := range notas {
		lista = append(lista, *n)
	}
	json.NewEncoder(w).Encode(lista)
}

func imprimirNota(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"ID de nota inválido"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	var nota *NotaFiscal

	if db != nil {
		var n NotaFiscal
		var itensStr string
		err := db.QueryRow("SELECT id, numero, status, itens FROM notas WHERE id = $1", id).Scan(&n.ID, &n.Numero, &n.Status, &itensStr)
		if err != nil {
			mu.Unlock()
			http.Error(w, `{"error":"Nota fiscal não encontrada"}`, http.StatusNotFound)
			return
		}
		json.Unmarshal([]byte(itensStr), &n.Itens)
		nota = &n
	} else {
		var exists bool
		nota, exists = notas[id]
		if !exists {
			nota = &NotaFiscal{
				ID:     id,
				Numero: id,
				Status: "Aberta",
				Itens:  []ItemNota{{ProdutoCodigo: "PROD01", Quantidade: 1}},
			}
			notas[id] = nota
		}
	}

	if nota.Status != "Aberta" {
		mu.Unlock()
		http.Error(w, `{"error":"Apenas notas com status 'Aberta' podem ser impressas"}`, http.StatusBadRequest)
		return
	}
	mu.Unlock()

	// Simula tempo de processamento de impressão
	time.Sleep(1 * time.Second)

	// Comunicação HTTP REST com o Serviço de Estoque
	client := http.Client{Timeout: 3 * time.Second}

	for _, item := range nota.Itens {
		prodCod := item.ProdutoCodigo
		if prodCod == "" {
			prodCod = "PROD01"
		}
		url := fmt.Sprintf("http://localhost:8081/produtos/%s/baixar", prodCod)

		payload, _ := json.Marshal(map[string]interface{}{
			"produtoCodigo": prodCod,
			"quantidade":    item.Quantidade,
		})

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))

		// Trata cenário de falha (Resiliência)
		if err != nil {
			log.Printf("FALHA: Serviço de Estoque indisponível: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Falha na impressão: Serviço de Estoque indisponível. A nota continua ABERTA.",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("FALHA no Estoque: %s\n", string(bodyBytes))
			w.WriteHeader(resp.StatusCode)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Falha ao dar baixa no estoque do produto '%s'. A nota permanece ABERTA.", prodCod),
			})
			return
		}
	}

	// Atualiza status da nota para Fechada se a baixa no estoque funcionou
	mu.Lock()
	nota.Status = "Fechada"

	if db != nil {
		_, _ = db.Exec("UPDATE notas SET status = 'Fechada' WHERE id = $1", id)
	}
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensagem": "Nota Fiscal impressa e fechada com sucesso!",
		"nota":     nota,
	})
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.Use(enableCORS)

	r.HandleFunc("/notas", cadastrarNota).Methods("POST", "OPTIONS")
	r.HandleFunc("/notas", listarNotas).Methods("GET", "OPTIONS")
	r.HandleFunc("/notas/{id}/imprimir", imprimirNota).Methods("POST", "OPTIONS")

	fmt.Println("Serviço de Faturamento rodando na porta :8082...")
	log.Fatal(http.ListenAndServe(":8082", r))
}
