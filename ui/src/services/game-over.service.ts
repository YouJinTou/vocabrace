import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { ContextService, Status } from './context.service';

@Injectable({
  providedIn: 'root'
})
export class GameOverService {

  constructor(private httpClient: HttpClient, private contextService: ContextService) { }

  onGameOver(payload) {
    this.contextService.setStatus(new Status());
    // this.httpClient.post(environment.gameOverEndpoint, payload);
  }
}
