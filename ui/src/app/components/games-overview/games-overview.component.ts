import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { WebsocketService } from 'src/services/websocket.service';

@Component({
  selector: 'games-overview',
  templateUrl: './games-overview.component.html',
  styleUrls: ['./games-overview.component.css']
})
export class GamesOverviewComponent implements OnInit {
  constructor(private wsService: WebsocketService, private router: Router) { }

  ngOnInit(): void {
  }

  onScrabblePlayersChanged(e) {
    this.wsService.connect(environment.wsEndpoint, 'scrabble').subscribe({
      next: m => this.redirect(m),
      error: e => console.log(e)
    });
  }

  private redirect(m: any) {
    this.router.navigate(['scrabble', 'test']);
  }
}
