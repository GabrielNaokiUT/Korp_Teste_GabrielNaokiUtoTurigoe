import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class NotaService {
  // Aponta para o nosso backend de Faturamento em Golang
  private apiUrl = 'http://localhost:8082/notas';

  constructor(private http: HttpClient) {}

  // Utiliza RxJS para capturar erros e repassar de forma amigável
  imprimirNota(id: number): Observable<any> {
    return this.http.post(`${this.apiUrl}/${id}/imprimir`, {}, { responseType: 'text' }).pipe(
      catchError(erro => {
        console.error('Falha no microsserviço:', erro);
        return throwError(() => new Error('Falha na comunicação com o estoque. A nota não foi impressa.'));
      })
    );
  }
}
