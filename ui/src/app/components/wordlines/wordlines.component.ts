import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { environment } from 'src/environments/environment';
import { ContextService } from 'src/services/context.service';
import { GameOverService } from 'src/services/game-over.service';
import { UsernameService } from 'src/services/username.service';
import { WebsocketService } from 'src/services/websocket.service';
import { TimerComponent } from '../timer/timer.component';
import { BlanksDialog } from './blanks/blanks.component';
import { Cell, getCellClass } from './cell';
import { GameOverDialog } from './game-over/game-over.component';
import { Payload } from './payload';
import { Player } from './player';
import { Tile } from './tile';

const GAME = 'wordlines';

@Component({
  selector: 'wordlines',
  templateUrl: './wordlines.component.html',
  styleUrls: ['./wordlines.component.css']
})
export class WordlinesComponent implements OnInit {
  private placedTiles: Cell[] = [];
  private originalTiles: Tile[] = [];
  private poolID: string;
  @ViewChild(TimerComponent)
  private timer: TimerComponent;
  private blanks: Tile[] = [];
  private blankClicked: boolean;
  timeout = 900;
  payload: Payload;
  players: Player[] = [];
  tiles: Tile[] = [];
  cells: Cell[] = [];

  constructor(
    public blanksDialog: MatDialog,
    public gameOverDialog: MatDialog,
    private wsService: WebsocketService,
    private contextService: ContextService,
    private gameOverService: GameOverService,
    private usernameService: UsernameService) { }

  ngOnInit(): void {
    this.loadCells();
    this.connect();
  }

  onTimeout() {
    if (this.payload.yourMove) {
      this.onPassClicked();
      this.cancel();
    }
  }

  onPlayerTileClicked(t: Tile) {
    if (!this.payload.yourMove || this.payload.isGameOver) {
      return;
    }

    t.selected = !t.selected;
    this.blankClicked = t.selected && t.isBlank();

    if (this.blankClicked) {
      this.openBlanks(t);
    }
  }

  onCellTileClicked(c: Cell) {
    if (!this.payload.yourMove || this.payload.isGameOver) {
      return;
    }
    if (this.removeCellTile(c)) {
      return;
    }
    if (c.isEmpty() && this.singleTileSelected()) {
      this.setCellTile(c);
    }
  }

  onPlaceClicked() {
    if (!this.placedTiles.length) {
      return
    }
    let payload = {
      g: GAME,
      p: true,
      w: [],
      pid: this.poolID
    };
    for (var c of this.placedTiles) {
      payload.w.push({
        c: c.id,
        t: c.tile.id,
        b: c.tile.isBlank() ? c.tile.letter : null,
      })
    }
    this.wsService.send(payload);
  }

  onExchangeClicked() {
    if (!this.selected()) {
      return;
    }
    let payload = {
      g: GAME,
      e: true,
      t: this.selected().map(t => t.id),
      pid: this.poolID
    };
    this.wsService.send(payload);
  }

  onPassClicked() {
    let payload = {
      g: GAME,
      q: true,
      pid: this.poolID
    };
    this.wsService.send(payload);
  }

  onCancelClicked() {
    this.cancel();
  }

  private pipeline(m: any) {
    console.log(m);
    this.payload = new Payload(m, this.usernameService);

    if (this.payload.isError) {
      this.cancel();
      return;
    }

    this.placedTiles = [];
    this.renderPlayers();
    this.renderPlayerTiles();
    this.handleExchange();
    this.handlePlace();
    this.onGameOver();
    this.originalTiles = this.tiles;
    this.blanks = this.payload.blanks;
    this.blankClicked = false;
    this.poolID = this.poolID ? this.poolID : this.payload.poolId;
    this.startTimer();
  }

  private startTimer() {
    if (this.payload.isGameOver) {
      return;
    }

    setTimeout(() => {
      if (this.timer) {
        this.timer.restart();
      }
    }, 0.5);
  }

  private cancel() {
    this.tiles = this.originalTiles.map(ot => { ot.selected = false; return ot; });
    for (var pc of this.placedTiles)
      for (var c of this.cells) {
        if (pc.id == c.id) {
          c.tile = null;
        }
      }
    this.placedTiles = [];
  }

  private loadCells() {
    this.cells = [];
    let i = 0;
    for (let r = 0; r < 15; r++) {
      for (let c = 0; c < 15; c++) {
        let cell = new Cell(i, null, getCellClass(i));
        this.cells.push(cell);
        i++;
      }
    }
  }

  private renderPlayers() {
    if (!this.payload.players) {
      return;
    }

    this.players = this.payload.players;
  }

  private renderPlayerTiles() {
    if (!this.payload.tiles) {
      return
    }

    this.tiles = [];
    for (var t of this.payload.tiles) {
      this.tiles.push(t);
      this.originalTiles.push(t.copy());
    }
  }

  private handleExchange() {
    if (!(this.payload.wasExchange && this.payload.exchangeTiles)) {
      return;
    }
    this.tiles = this.tiles.filter(t => !t.selected);
    this.tiles.push(...this.payload.exchangeTiles);
  }

  private handlePlace() {
    if (this.payload.yourMove && this.payload.wasPlace) {
      for (var c of this.payload.placedCells) {
        this.cells[c.id] = c;
      }
    } else if (this.payload.wasPlace && this.payload.exchangeTiles) {
      this.tiles = this.tiles.filter(t => !t.selected);
      this.tiles.push(...this.payload.exchangeTiles);
    }
  }

  private selected(): Tile[] {
    let selectedTiles = this.tiles.filter(t => t.selected);
    return selectedTiles.length == 0 ? null : selectedTiles;
  }

  private current(): Tile {
    if (!this.selected()) {
      return null;
    }
    if (this.selected().length == 1) {
      return this.selected()[0];
    }
    return null;
  }

  private singleTileSelected(): boolean {
    return this.current() != null;
  }

  private setCellTile(c: Cell) {
    c.tile = this.current().copy()
    this.tiles = this.tiles.filter(t => t.id != this.current().id);
    this.placedTiles.push(c.copy());
  }

  private removeCellTile(c: Cell): boolean {
    let shouldReturnTile = this.placedTiles.filter(pc => pc.id == c.id).length > 0;

    if (shouldReturnTile && !this.selected()) {
      this.placedTiles = this.placedTiles.filter(t => t.id != c.id);
      this.tiles.push(c.tile.copy());
      c.tile = null;
      return true;
    }

    return false;
  }

  private openBlanks(selectedBlank: Tile) {
    const blanksRef = this.blanksDialog.open(BlanksDialog, {
      data: { blanks: this.blanks }
    });

    blanksRef.afterClosed().subscribe(result => {
      if (result) {
        selectedBlank.letter = result.letter;
      }
    });
  }

  private onGameOver() {
    if (!this.payload.isGameOver) {
      return;
    }

    this.gameOverDialog.open(GameOverDialog, { data: this.payload });
    this.gameOverService.onGameOver(this.payload);
  }

  private connect() {
    console.log('history', this.wsService.history);
    this.wsService.history.map(h => this.pipeline(h));

    const params = {
      'game': 'wordlines',
      'players': this.contextService.wordlines.players,
      'language': this.contextService.wordlines.language,
      'userID': this.contextService.wordlines.userId,
      'pid': this.contextService.wordlines.poolId
    };
    this.wsService.connect(environment.wsEndpoint, params).subscribe({
      next: m => this.pipeline(m),
      error: e => console.log(e)
    });
  }
}