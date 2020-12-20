import { Component, OnDestroy, OnInit } from '@angular/core';

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
  selectedPlayers: string;
  selectedLanguage: string;

  constructor(
    private wsService: WebsocketService,
    private contextService: ContextService,
    private router: Router) { }

  ngOnInit(): void {
  }

  onSelectChanged() {
    this.connect();
  }

  private connect() {
    if (!(this.selectedLanguage && this.selectedPlayers)) {
      return;
    }
    
    this.wsService.connect(environment.wsEndpoint, {
      'game': 'wordlines',
      'players': parseInt(this.selectedPlayers),
      'language': this.selectedLanguage,
      'userID': this.contextService.user.id,
    }).subscribe({
      next: m => {
        if ('pid' in m) {
          this.contextService.setWordlines({
            players: parseInt(this.selectedPlayers),
            userId: this.contextService.user.id,
            language: this.selectedLanguage,
            poolId: m['pid']
          });
          this.contextService.setIsPlaying(true);
          this.router.navigate(['wordlines', m['pid']]);
        }
      },
      error: e => console.log(e)
    });
  }
}
