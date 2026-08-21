import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-cadastro-nota',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  templateUrl: './cadastro-nota.html'
})
export class CadastroNotaComponent implements OnInit {
  nota = {
    numero: 1,
    status: 'Aberta',
    itens: [{ produtoCodigo: '', quantidade: 1 }]
  };
  mensagem = '';
  notasList: any[] = [];
  carregandoId: number | null = null;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.carregarNotas();
  }

  carregarNotas() {
    this.http.get<any[]>('http://localhost:8082/notas').subscribe({
      next: (res) => this.notasList = res || [],
      error: () => {}
    });
  }

  adicionarItem() {
    this.nota.itens.push({ produtoCodigo: '', quantidade: 1 });
  }

  cadastrar() {
    this.http.post('http://localhost:8082/notas', this.nota).subscribe({
      next: () => {
        this.mensagem = 'Nota Fiscal criada com sucesso!';
        this.nota = { numero: this.nota.numero + 1, status: 'Aberta', itens: [{ produtoCodigo: '', quantidade: 1 }] };
        this.carregarNotas();
      },
      error: () => this.mensagem = 'Erro ao criar nota fiscal.'
    });
  }

  imprimir(id: number) {
    this.carregandoId = id;
    this.mensagem = '';
    this.http.post(`http://localhost:8082/notas/${id}/imprimir`, {}).subscribe({
      next: (res: any) => {
        this.carregandoId = null;
        this.mensagem = res.mensagem || 'Nota impressa e fechada!';
        this.carregarNotas();
      },
      error: (err) => {
        this.carregandoId = null;
        this.mensagem = err.error?.error || 'Erro ao processar impressão.';
        this.carregarNotas();
      }
    });
  }
}
