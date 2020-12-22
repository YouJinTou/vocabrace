import { Component, OnInit } from '@angular/core';

import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { ContextService } from 'src/services/context.service';
import { WebsocketService } from 'src/services/websocket.service';

@Component({
  selector: 'app-wordlines-overview',
  templateUrl: './wordlines-overview.component.html',
  styleUrls: ['./wordlines-overview.component.css']
})
export class WordlinesOverviewComponent implements OnInit {
  gameExists = true;
  selectedPlayers: string;
  selectedLanguage: string;

  constructor(
    private wsService: WebsocketService,
    private contextService: ContextService,
    private router: Router) { }

  ngOnInit(): void {
    this.contextService.isPlaying$.subscribe(r => this.gameExists = r.value);
  }

  onSelectChanged() {
    this.connect();
  }
  
  backToGame() {
    this.router.navigate(['wordlines', this.contextService.isPlaying.pid]);
  }

  private connect() {
    if (!(this.selectedLanguage && this.selectedPlayers)) {
      return;
    }

    if (!this.contextService.isPlaying.value) {
      this.wsService.close();
    }

    this.wsService.connect(environment.wsEndpoint, {
      'game': 'wordlines',
      'players': parseInt(this.selectedPlayers),
      'language': this.selectedLanguage,
      'userID': this.contextService.user.id,
    }).subscribe({
      next: m => {
        if ('pid' in m) {
          this.contextService.setIsPlaying({
            value: true,
            pid: m['pid'],
            language: this.selectedLanguage,
            players: parseInt(this.selectedPlayers),
          });
          this.router.navigate(['wordlines', m['pid']]);
        }
      },
      error: e => console.log(e)
    });
  }
}