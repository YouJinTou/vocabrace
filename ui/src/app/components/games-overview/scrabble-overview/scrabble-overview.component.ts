import { Component, OnInit } from '@angular/core';

import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { UserStatusService } from 'src/services/user-status.service';
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
    private userStatusService: UserStatusService, 
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
    this.wsService.connect(environment.wsEndpoint, 'scrabble', {
      'players': parseInt(this.selectedPlayers),
      'language': this.selectedLanguage,
      'userId': this.userStatusService.current.id,
    }).subscribe({
        next: m => this.router.navigate(['scrabble', m['pid']]),
        error: e => console.log(e)
      });
  }
}
