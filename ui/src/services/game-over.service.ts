import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { ContextService, IsPlaying } from './context.service';

@Injectable({
  providedIn: 'root'
})
export class GameOverService {

  constructor(private httpClient: HttpClient, private contextService: ContextService) { }

  onGameOver(payload) {
    this.contextService.setIsPlaying(new IsPlaying());
    // this.httpClient.post(environment.gameOverEndpoint, payload);
  }
}
