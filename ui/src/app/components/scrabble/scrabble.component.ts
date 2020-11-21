import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { WebsocketService } from 'src/services/websocket.service';
import { Cell } from './cell';
import { Tile } from './tile';

const GAME = 'scrabble';
const DOUBLE_LETTER_INDICES = [3, 11, 36, 38, 45, 52, 59, 92, 96, 98, 102, 108, 116, 122, 126, 128, 132, 165, 172, 179, 186, 188, 213, 221];
const DOUBLE_WORD_INDICES = [16, 32, 48, 64, 112, 160, 176, 192, 208, 28, 42, 56, 70, 154, 168, 182, 196];
const TRIPLE_LETTER_INDICES = [20, 24, 76, 80, 84, 88, 136, 140, 144, 148, 200, 204];
const TRIPLE_WORD_INDICES = [0, 7, 14, 105, 119, 210, 217, 224];

@Component({
  selector: 'scrabble',
  templateUrl: './scrabble.component.html',
  styleUrls: ['./scrabble.component.css']
})
export class ScrabbleComponent implements OnInit, OnDestroy {
  private destroyed$ = new Subject();
  private selected: Tile;
  private placedTiles: Cell[] = [];
  private originalTiles: Tile[] = [];
  tiles: Tile[] = [];
  cells: Cell[] = [];

  constructor(private wsService: WebsocketService) { }

  ngOnInit(): void {
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
    let shouldDeselect = t == this.selected;
    if (shouldDeselect) {
      this.selected = null;
      return;
    }
    this.selected = t;
  }

  onCellTileClicked(c: Cell) {
    if (this.selected == null) {
      return;
    }
    if (c.tile != null) {
      return;
    }
    c.tile = this.selected.copy()
    this.tiles = this.tiles.filter(t => t.id != this.selected.id);
    this.placedTiles.push(c.copy());
    this.selected = null;
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
        b: c.tile.letter,
      })
    }
    this.wsService.send(payload);
    this.placedTiles = [];
  }

  onExchangeClicked() {
    let payload = {
      g: GAME,
      e: true,
      t: [] //TODO
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
    this.tiles = this.originalTiles;
  }

  private pipeline(m: any) {
    console.log(m);
    if (this.isError(m)) {
      return;
    }

    this.loadCells();
    this.renderPlayerTiles(m);
  }

  private isError(m: any): boolean {
    return 'message' in m && m['message'].indexOf('Internal server error') > -1;
  }

  private loadCells() {
    this.cells = [];
    let i = 0;
    for (let r = 0; r < 15; r++) {
      for (let c = 0; c < 15; c++) {
        let cell = new Cell(i, null,
          DOUBLE_LETTER_INDICES.includes(i) ? 'double-letter' :
            TRIPLE_LETTER_INDICES.includes(i) ? 'triple-letter' :
              DOUBLE_WORD_INDICES.includes(i) ? 'double-word' :
                TRIPLE_WORD_INDICES.includes(i) ? 'triple-word' :
                  'tile'
        );
        this.cells.push(cell);
        i++;
      }
    }
  }

  private renderPlayerTiles(data: any) {
    if ('t' in data) {
      this.tiles = [];
      for (var t of data['t']) {
        let tokens = t.split('|');
        let tile = new Tile(tokens[0], tokens[1], tokens[2]);
        this.tiles.push(tile);
        this.originalTiles.push(tile);
      }
    }
  }
}
