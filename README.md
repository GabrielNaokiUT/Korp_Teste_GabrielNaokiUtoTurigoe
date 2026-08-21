# Sistema de Emissão de Notas Fiscais - Desafio Técnico KORP

Este projeto é a solução desenvolvida para o desafio técnico da KORP, contemplando um sistema de emissão de notas fiscais construído com arquitetura de microsserviços[cite: 2]. 

## 🏗️ Arquitetura da Solução
O sistema foi dividido em três camadas principais:
- **Frontend:** Desenvolvido em Angular[cite: 1].
- **Backend (Microsserviços):** Desenvolvido em Golang, dividido em Serviço de Estoque e Serviço de Faturamento[cite: 1, 2].
- **Banco de Dados:** PostgreSQL compartilhado pelos dois serviços, com tabelas independentes[cite: 2].

## ⚙️ Funcionalidades Implementadas
1. **Cadastro de Produtos:** Permite cadastrar produtos com Código, Descrição e Saldo[cite: 1].
2. **Cadastro de Notas Fiscais:** Criação de notas com numeração sequencial, status inicial "Aberta" e inclusão de múltiplos produtos[cite: 1].
3. **Impressão de Notas Fiscais:** Botão de impressão que processa a nota, comunica-se com o estoque para baixar o saldo e atualiza o status para "Fechada"[cite: 1, 2].
4. **Tratamento de Falhas:** Caso o Serviço de Estoque falhe durante a impressão, a operação é abortada, a nota permanece "Aberta" e o erro é informado ao usuário de forma clara[cite: 1, 2].

## 🛠️ Detalhamento Técnico

### Frontend (Angular)
- **Ciclos de vida utilizados:** 
  - `ngOnInit`: Para carregar listagens iniciais e formulários[cite: 2].
  - `ngOnDestroy`: Para cancelar inscrições (subscriptions) e evitar memory leaks[cite: 2].
  - `ngOnChanges`: Para reagir a alterações de propriedades (Inputs)[cite: 2].
  - `ngAfterViewInit`: Para configurações de paginação/ordenação na tabela[cite: 2].
- **Uso do RxJS:** Utilizado para gerenciar requisições e estados com os operadores `switchMap` (cancelar requisições em andamento), `tap` (efeitos colaterais/logs), `catchError` (tratamento de exceções HTTP), `forkJoin` (requisições paralelas) e `BehaviorSubject` (compartilhamento de estado)[cite: 2].
- **Bibliotecas Visuais:** Angular Material (tabelas, formulários, botões, spinners e snackbars)[cite: 2].
- **Outras bibliotecas:** `@angular/forms` (formulários reativos), `ng-toastr` (notificações) e `date-fns` (formatação de datas)[cite: 2].

### Backend (Golang)
- **Gerenciamento de dependências:** Feito nativamente via Go Modules (`go.mod` e `go.sum`)[cite: 2].
- **Frameworks:** Utilizado o `Gorilla Mux` para roteamento de endpoints HTTP e o pacote nativo `database/sql` em conjunto com o driver `github.com/lib/pq` para conexão com o PostgreSQL[cite: 2]. *(Nota: Não foi utilizado C# ou LINQ no projeto[cite: 2]).*
- **Tratamento de Erros:** 
  - Validação de entrada com retorno de HTTP 400[cite: 2].
  - Falhas de banco de dados registradas em log e retornando HTTP 500[cite: 2].
  - Falha de comunicação entre microsserviços (Estoque e Faturamento) resulta na interrupção da operação, log do erro e manutenção do status da nota fiscal, garantindo a consistência[cite: 2].

## 🚀 Como Executar o Projeto
**Pré-requisitos:** Node.js, Go e PostgreSQL instalados[cite: 2].

1. **Serviço de Estoque:** `cd backend/estoque-service` -> `go run main.go`[cite: 2].
2. **Serviço de Faturamento:** `cd backend/faturamento-service` -> `go run main.go`[cite: 2].
3. **Frontend:** `cd frontend` -> `npm install` -> `ng serve`[cite: 2].

