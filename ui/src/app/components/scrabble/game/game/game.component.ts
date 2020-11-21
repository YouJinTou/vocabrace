import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { WebsocketService } from 'src/services/websocket.service';

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.css']
})
export class GameComponent implements OnInit, OnDestroy {
  private destroyed$ = new Subject();

  constructor(private wsService: WebsocketService) { }

  ngOnInit(): void {
    this.wsService.connect(environment.wsEndpoint).pipe(
      takeUntil(this.destroyed$)
    ).subscribe(m => console.log(m));
  }

  ngOnDestroy(): void {
    this.destroyed$.next();
  }
}
