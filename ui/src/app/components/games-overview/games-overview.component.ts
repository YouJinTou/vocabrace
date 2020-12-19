import { Component, OnInit } from '@angular/core';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'games-overview',
  templateUrl: './games-overview.component.html',
  styleUrls: ['./games-overview.component.css']
})
export class GamesOverviewComponent implements OnInit {
  constructor(private contextService: ContextService) { }

  ngOnInit(): void {
    this.contextService.setIsPlaying(false);
  }
}
