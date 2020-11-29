import { Component, OnInit } from '@angular/core';

import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { WebsocketService } from 'src/services/websocket.service';

@Component({
  selector: 'app-scrabble-overview',
  templateUrl: './scrabble-overview.component.html',
  styleUrls: ['./scrabble-overview.component.css']
})
export class ScrabbleOverviewComponent implements OnInit {
  selectedPlayers: string;
  selectedLanguage: string;

  constructor(private wsService: WebsocketService, private router: Router) { }

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
    }).subscribe({
      next: m => this.router.navigate(['scrabble', m['pid']]),
      error: e => console.log(e)
    });
  }
}
