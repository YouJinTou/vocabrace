import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CookieService } from 'ngx-cookie-service';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'games-overview',
  templateUrl: './games-overview.component.html',
  styleUrls: ['./games-overview.component.css']
})
export class GamesOverviewComponent implements OnInit {
  ongoingGameExists = false;

  constructor(
    private contextService: ContextService,
    private cookieService: CookieService,
    private router: Router) { }

  ngOnInit(): void {
    this.ongoingGameExists = this.contextService.isPlaying.value || (
      this.contextService.user.loggedIn && this.cookieService.check('pid'));
  }

  backToGame() {
    const pid = this.contextService.isPlaying.pid || this.cookieService.get('pid');
    this.router.navigate(['wordlines', pid]);
  }
}
