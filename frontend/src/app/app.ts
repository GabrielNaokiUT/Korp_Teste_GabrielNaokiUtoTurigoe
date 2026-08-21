import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { CadastroProdutoComponent } from './pages/cadastro-produto/cadastro-produto';
import { CadastroNotaComponent } from './pages/cadastro-nota/cadastro-nota';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    CadastroProdutoComponent,
    CadastroNotaComponent
  ],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  title = 'frontend';
}
