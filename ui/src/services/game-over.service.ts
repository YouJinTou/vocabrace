import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class GameOverService {

  constructor(private httpClient: HttpClient) { }

  onGameOver(payload) {
    // this.httpClient.post(environment.gameOverEndpoint, payload);
  }
}
