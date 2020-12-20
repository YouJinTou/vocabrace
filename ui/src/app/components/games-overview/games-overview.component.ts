import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'games-overview',
  templateUrl: './games-overview.component.html',
  styleUrls: ['./games-overview.component.css']
})
export class GamesOverviewComponent implements OnInit {
  ongoingGameExists = false;

  constructor(private contextService: ContextService, private router: Router) { }

  ngOnInit(): void {
    this.ongoingGameExists = this.contextService.isPlaying.value;
  }

  backToGame() {
    this.router.navigate(['wordlines', this.contextService.isPlaying.pid]);
  }
}
