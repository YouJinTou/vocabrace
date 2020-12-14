import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
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
  private poolID: string;
  payload: Payload;
  players: Player[] = [];
  tiles: Tile[] = [];
  cells: Cell[] = [];
  blanks: Tile[] = [];
  blankClicked: boolean;

  constructor(
    public blanksDialog: MatDialog,
    private wsService: WebsocketService,
    private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.loadCells();

    if (this.wsService.last) {
      this.pipeline(this.wsService.last);
    }

    let connection$ = this.wsService.connection();

    if (connection$) {
      connection$.pipe(
        takeUntil(this.destroyed$)
        )
        .subscribe({
        next: m => this.pipeline(m),
        error: e => console.log(e)
      });
    } else {
      console.log('not connected.');
    }
  }

  ngOnDestroy(): void {
    this.destroyed$.next();
  }

  onPlayerTileClicked(t: Tile) {
    if (!this.payload.yourMove) {
      return;
    }

    t.selected = !t.selected;
    this.blankClicked = t.selected && t.isBlank();

    if (this.blankClicked) {
      this.openBlanks(t);
    }
  }

  onCellTileClicked(c: Cell) {
    if (!this.payload.yourMove) {
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
    this.payload = new Payload(m);
    console.log('entering...')
    if (this.payload.isError) {
    console.log('error...')
    this.cancel();
      return;
    }

    this.placedTiles = [];
    console.log('players...')
    this.renderPlayers();
    console.log('tiles...')
    this.renderPlayerTiles();
    console.log('exch...')
    this.handleExchange();
    console.log('place...')
    this.handlePlace();
    this.originalTiles = this.tiles;
    this.blanks = this.payload.blanks;
    this.blankClicked = false;
    this.poolID = this.poolID ? this.poolID : this.payload.poolId;
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
      console.log('placed cells:');
      console.log(this.payload.placedCells);
      for (var c of this.payload.placedCells) {
      console.log('setting...');
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
}

@Component({
  selector: 'scrabble-blanks',
  templateUrl: 'blanks.component.html',
  styleUrls: ['./scrabble.component.css']
})
export class BlanksDialog {
  blanks: Tile[];

  constructor(
    public dialogRef: MatDialogRef<BlanksDialog>,
    @Inject(MAT_DIALOG_DATA) public data: { blanks: Tile[] }) {
    this.blanks = data['blanks'];
  }

  onSelected(t: Tile) {
    this.dialogRef.close(t);
  }
}