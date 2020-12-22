import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Observable, of } from 'rxjs';
import { environment } from 'src/environments/environment';
import { ContextService } from 'src/services/context.service';
import { GameOverService } from 'src/services/game-over.service';
import { UsernameService } from 'src/services/username.service';
import { WebsocketService } from 'src/services/websocket.service';
import { TimerComponent } from '../timer/timer.component';
import { BlanksDialog } from './blanks/blanks.component';
import { Cell } from './cell';
import { GameOverDialog } from './game-over/game-over.component';
import { Payload } from './payload';
import { State } from './state';
import { Tile } from './tile';

const GAME = 'wordlines';

@Component({
  selector: 'wordlines',
  templateUrl: './wordlines.component.html',
  styleUrls: ['./wordlines.component.css']
})
export class WordlinesComponent implements OnInit {
  @ViewChild(TimerComponent)
  private timer: TimerComponent;
  timeout = 60;
  state = new State();
  tilesRemaining = [];

  constructor(
    public blanksDialog: MatDialog,
    public gameOverDialog: MatDialog,
    private wsService: WebsocketService,
    private contextService: ContextService,
    private gameOverService: GameOverService,
    private usernameService: UsernameService) { }

  ngOnInit(): void {
    this.connect();
  }

  canDeactivate(): Observable<boolean> {
    return of(confirm('You will not be able to return. Continue?'));
  }

  onTimeout() {
    if (this.state.yourMove) {
      this.onPassClicked();
      this.state = this.state.cancel();
    }
  }

  onPlayerTileClicked(t: Tile) {
    this.state = this.state.clickPlayerTile(t);
    if (this.state.blankClicked) {
      this.openBlanks();
    }
  }

  onCellTileClicked(c: Cell) {
    this.state = this.state.clickCellTile(c);
  }

  onPlaceClicked() {
    if (!this.state.currentPlacedCells.length) {
      return
    }

    const payload = {
      g: GAME,
      p: true,
      w: this.state.currentPlacedCells.map(p => ({
        c: p.id,
        t: p.tile.id,
        b: p.tile.isBlank() ? p.tile.letter : null,
      })),
      pid: this.state.poolId
    };
    this.wsService.send(payload);
  }

  onExchangeClicked() {
    if (!this.state.selected()) {
      return;
    }
    const payload = {
      g: GAME,
      e: true,
      t: this.state.selected().map(t => t.id),
      pid: this.state.poolId
    };
    this.wsService.send(payload);
  }

  onPassClicked() {
    const payload = {
      g: GAME,
      q: true,
      pid: this.state.poolId
    };
    this.wsService.send(payload);
  }

  onCancelClicked() {
    this.state = this.state.cancel();
  }

  private process(m: any) {
    const p = new Payload(m, this.usernameService);
    this.state = this.state.apply(p);
    this.tilesRemaining = Array(this.state.tilesRemaining).fill(1);

    if (this.state.isError) {
      return;
    }

    this.onGameOver();
    this.startTimer();
  }

  private startTimer() {
    if (this.state.isGameOver) {
      return;
    }

    setTimeout(() => {
      if (this.timer) {
        this.timer.restart();
      }
    }, 0.5);
  }

  private openBlanks() {
    const blanksRef = this.blanksDialog.open(BlanksDialog, {
      data: { blanks: this.state.blanks }
    });

    blanksRef.afterClosed().subscribe(result => {
      if (result) {
        this.state = this.state.setBlank(result);
      }
    });
  }

  private onGameOver() {
    if (!this.state.isGameOver) {
      return;
    }

    this.gameOverDialog.open(GameOverDialog, { data: this.state });
    this.gameOverService.onGameOver(this.state);
  }

  private connect() {
    this.wsService.history.map(h => this.process(h));
    const params = {
      'game': 'wordlines',
      'players': this.contextService.status.players,
      'language': this.contextService.status.language,
      'userID': this.contextService.user.id,
      'pid': this.contextService.status.pid
    };
    this.wsService.connect(environment.wsEndpoint, params).subscribe({
      next: m => {
        if (Array.isArray(m)) {
          m.forEach(x => {
            this.process(JSON.parse(x));
          });
        } else {
          this.process(m);
        }
      },
      error: e => console.log(e)
    });
  }
}