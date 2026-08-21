import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-cadastro-produto',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  templateUrl: './cadastro-produto.html'
})
export class CadastroProdutoComponent {
  produto = { codigo: '', descricao: '', saldo: 0 };
  mensagem = '';

  constructor(private http: HttpClient) {}

  cadastrar() {
    this.http.post('http://localhost:8081/produtos', this.produto).subscribe({
      next: () => {
        this.mensagem = 'Produto cadastrado com sucesso!';
        this.produto = { codigo: '', descricao: '', saldo: 0 };
      },
      error: () => this.mensagem = 'Erro ao cadastrar produto.'
    });
  }
}
