import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { ContextService } from './context.service';

@Injectable({
  providedIn: 'root'
})
export class GameOverService {

  constructor(private httpClient: HttpClient, private contextService: ContextService) { }

  onGameOver(payload) {
    this.contextService.setIsPlaying({ value: false, pid: '' });
    // this.httpClient.post(environment.gameOverEndpoint, payload);
  }
}
