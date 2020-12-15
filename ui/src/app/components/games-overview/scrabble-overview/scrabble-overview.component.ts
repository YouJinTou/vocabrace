import { Component, OnDestroy, OnInit } from '@angular/core';

import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { ContextService } from 'src/services/context.service';
import { WebsocketService } from 'src/services/websocket.service';

@Component({
  selector: 'app-scrabble-overview',
  templateUrl: './scrabble-overview.component.html',
  styleUrls: ['./scrabble-overview.component.css']
})
export class ScrabbleOverviewComponent implements OnInit {
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
      'game': 'scrabble',
      'players': parseInt(this.selectedPlayers),
      'language': this.selectedLanguage,
      'userID': this.contextService.user.id,
      'isAnonymous': this.contextService.user.id ? true : false
    }).subscribe({
      next: m => {
        if ('pid' in m) {
          this.contextService.setIsPlaying(true);
          this.router.navigate(['scrabble', m['pid']]);
        }
      },
      error: e => console.log(e)
    });
  }
}
