import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { WebsocketService } from 'src/services/websocket.service';
import { Cell, getCellClass } from './cell';
import { Payload } from './payload';
import { Player } from './player';
import { Tile } from './tile';

const GAME = 'scrabble';

@Component({
  selector: 'scrabble',
  templateUrl: './scrabble.component.html',
  styleUrls: ['./scrabble.component.css']
})
export class ScrabbleComponent implements OnInit, OnDestroy {
  private destroyed$ = new Subject();
  private placedTiles: Cell[] = [];
  private originalTiles: Tile[] = [];
  private payload: Payload;
  players: Player[] = [];
  tiles: Tile[] = [];
  cells: Cell[] = [];

  constructor(private wsService: WebsocketService) { }

  ngOnInit(): void {
    this.loadCells();
    this.wsService.connect(environment.wsEndpoint, GAME).pipe(
      takeUntil(this.destroyed$)
    ).subscribe({
      next: m => this.pipeline(m),
      error: e => console.log(e)
    });
  }

  ngOnDestroy(): void {
    this.destroyed$.next();
  }

  onPlayerTileClicked(t: Tile) {
    t.selected = !t.selected;
  }

  onCellTileClicked(c: Cell) {
    if (this.removeCellTile(c)) {
      return;
    }
    if (c.isEmpty() && this.singleTileSelected()) {
      this.setCellTile(c);
    }
  }

  onPlaceClicked() {
    let payload = {
      g: GAME,
      p: true,
      w: []
    };
    for (var c of this.placedTiles) {
      payload.w.push({
        c: c.id,
        t: c.tile.id,
        b: c.tile.isBlank() ? c.tile.letter : null,
      })
    }
    this.wsService.send(payload);
    this.placedTiles = [];
  }

  onExchangeClicked() {
    if (!this.selected()) {
      return;
    }
    let payload = {
      g: GAME,
      e: true,
      t: this.selected().map(t => t.id)
    };
    this.wsService.send(payload);
  }

  onPassClicked() {
    let payload = {
      g: GAME,
      q: true
    };
    this.wsService.send(payload);
  }

  onCancelClicked() {
    this.tiles = this.originalTiles.map(ot => { ot.selected = false; return ot; });
    for (var pc of this.placedTiles)
      for (var c of this.cells) {
        if (pc.id == c.id) {
          c.tile = null;
        }
      }
    this.placedTiles = [];
  }

  private pipeline(m: any) {
    this.payload = new Payload(m);
    if (this.payload.isError) {
      return;
    }

    this.renderPlayers();
    this.renderPlayerTiles();
    this.handleExchange();
    this.handlePlace();
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
}
