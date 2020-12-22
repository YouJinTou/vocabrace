import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
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
    private httpClient: HttpClient,
    private router: Router) { }

  ngOnInit(): void {
    if (this.contextService.user.loggedIn) {
      this.getUserPool(this.contextService.user.id).subscribe({
        next: r => {
          this.ongoingGameExists = r && r.PoolID != '' && r.PoolID != undefined && r.PoolID != null;
          
          if (this.ongoingGameExists) {
            console.log(r);
            this.contextService.setIsPlaying({
              value: true, pid: r.PoolID, language: r.Language, players: r.Players
            });
          }
        },
        error: e => {
          console.log(e);
        }
      });
    }
  }

  backToGame() {
    this.router.navigate(['wordlines', this.contextService.isPlaying.pid]);
  }

  getUserPool(userId: string): Observable<any> {
    const url = `${environment.poolingEndpoint}/userpools/${userId}`
    return this.httpClient.get<any>(url);
  }
}
