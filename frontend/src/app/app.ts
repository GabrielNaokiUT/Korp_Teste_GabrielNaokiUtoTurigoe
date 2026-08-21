import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subscription } from 'rxjs';
import { NotaService } from './nota.service';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule, MatButtonModule, MatSnackBarModule],
  template: `
    <div style="padding: 40px; font-family: Roboto, sans-serif;">
      <h2>Sistema de Emissão de Notas Fiscais</h2>

      <div style="margin-bottom: 20px; padding: 15px; border: 1px solid #ccc; max-width: 400px;">
        <h3>Nota Fiscal #1</h3>
        <p><strong>Status:</strong> {{ status }}</p>

        <!-- Indicador de processamento do Angular Material[cite: 1, 2] -->
        <mat-spinner *ngIf="loading" diameter="30" style="margin-bottom: 15px;"></mat-spinner>

        <!-- Botão de impressão intuitivo, desabilita se não for "Aberta"[cite: 1] -->
        <button mat-raised-button color="primary"
                (click)="imprimir()"
                [disabled]="loading || status !== 'Aberta'">
          Imprimir Nota
        </button>
      </div>
    </div>
  `
})
export class AppComponent implements OnInit, OnDestroy {
  loading = false;
  status = 'Aberta'; // Status inicial obrigatório[cite: 1]
  private sub: Subscription = new Subscription();

  constructor(private notaService: NotaService, private snackBar: MatSnackBar) {}

  ngOnInit(): void {
    // Ciclo de vida: Disparado ao carregar o componente[cite: 2]
    console.log('Tela de nota fiscal carregada.');
  }

  imprimir(): void {
    this.loading = true; // Exibe o spinner de processamento[cite: 1]

    this.sub.add(
      this.notaService.imprimirNota(1).subscribe({
        next: () => {
          this.loading = false;
          this.status = 'Fechada'; // Atualiza status com sucesso[cite: 1]
          this.mostrarAviso('Nota impressa e finalizada com sucesso!');
        },
        error: (erro) => {
          this.loading = false;
          // Feedback apropriado sobre a falha (Resiliência)[cite: 1]
          this.mostrarAviso(erro.message, 'Entendi');
        }
      })
    );
  }

  mostrarAviso(mensagem: string, acao: string = 'Fechar') {
    this.snackBar.open(mensagem, acao, { duration: 5000 }); // Notificação visual[cite: 2]
  }

  ngOnDestroy(): void {
    // Ciclo de vida: Previne vazamento de memória cancelando inscrições[cite: 2]
    this.sub.unsubscribe();
  }
}
