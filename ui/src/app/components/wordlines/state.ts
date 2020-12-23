import { Cell, getCellClass } from './cell';
import { Payload } from './payload';
import { Player } from './player';
import { Tile } from './tile';

export class State {
    cells: Cell[];
    players: Player[];
    tiles: Tile[];
    blanks: Tile[];
    isGameOver: boolean;
    blankClicked: boolean;
    selectedBlank: Tile;
    yourMove: boolean;
    isError: boolean;
    poolId: string;
    winnerName: string;
    toMoveId: string;
    tilesRemaining: number;
    currentPlacedCells: Cell[];
    lastAction: string;
    displayMessage: string;
    clientId: string;
    clientLastMoved: boolean;

    constructor() {
        this.cells = [];
        this.players = [];
        this.tiles = [];
        this.blanks = [];
        this.currentPlacedCells = [];
        this.tilesRemaining = 0;
    }

    copy(cells?: Cell[], players?: Player[], tiles?: Tile[]): State {
        let s = new State();
        s.cells = cells ?? this.cells;
        s.players = players ?? this.players;
        s.tiles = tiles ?? this.tiles;
        s.currentPlacedCells = this.currentPlacedCells;
        s.blanks = this.blanks;
        s.blankClicked = this.blankClicked;
        s.selectedBlank = this.selectedBlank;
        s.isGameOver = this.isGameOver;
        s.yourMove = this.yourMove;
        s.isError = this.isError;
        s.winnerName = this.winnerName;
        s.toMoveId = this.toMoveId;
        s.poolId = this.poolId;
        s.tilesRemaining = this.tilesRemaining;
        s.lastAction = this.lastAction;
        s.displayMessage = this.displayMessage;
        s.clientId = this.clientId;
        s.clientLastMoved = this.clientLastMoved;
        return s;
    }

    public apply(p: Payload): State {
        let copy = this.copy();
        copy.isError = p.isError;
        copy.clientId = copy.clientId ? copy.clientId :
            copy.yourMove ? copy.toMoveId : '';
        copy.clientLastMoved = copy.clientId == p.lastMovedId;
        copy = copy.setLastAction(p);
        copy = copy.setDisplayMessage(p);
        if (p.isError) {
            return copy.cancel();
        }
        copy.currentPlacedCells = [];
        copy = copy.getCells();
        copy = copy.copy(null, p.players, p.tiles);
        copy = copy.handleExchange(p);
        copy = copy.handleSomeoneElsePlaced(p);
        copy = copy.handlePlayerPlaced(p);
        copy.isGameOver = p.isGameOver;
        copy.yourMove = p.yourMove;
        copy.blanks = p.blanks;
        copy.poolId = copy.poolId || p.poolId;
        copy.winnerName = p.winnerName;
        copy.toMoveId = p.toMoveId;
        copy.tilesRemaining = p.tilesRemaining;
        return copy;
    }

    public clickPlayerTile(t: Tile): State {
        let copy = this.copy();
        if (!copy.yourMove || copy.isGameOver) {
            return copy;
        }

        let tCopy = t.copy();
        tCopy.selected = !tCopy.selected;
        copy.blankClicked = tCopy.selected && tCopy.isBlank();
        copy.selectedBlank = copy.blankClicked ? tCopy : null;

        copy.tiles.find(ti => ti.id == tCopy.id).selected = tCopy.selected;

        return copy;
    }

    public setBlank(t: Tile): State {
        let copy = this.copy();
        let playerBlank = copy.tiles.find(pt => pt.id == copy.selectedBlank.id);
        playerBlank.letter = t.letter;
        return copy;
    }

    public clickCellTile(c: Cell): State {
        let copy = this.copy();
        if (!copy.yourMove || copy.isGameOver) {
            return copy;
        }

        let result = copy.removeCellTile(c);

        if (result[1]) {
            return result[0];
        }

        if (c.isEmpty() && copy.singleTileSelected()) {
            copy = copy.setCellTile(c, copy.current());
        }

        return copy;
    }

    public cancel(): State {
        let copy = this.copy();

        for (var pc of copy.currentPlacedCells) {
            if (!copy.tiles.some(t => t.id == pc.tile.id)) {
                copy.tiles.push(pc.tile.copy());
            }

            for (var c of copy.cells) {
                if (pc.id == c.id) {
                    c.tile = null;
                }
            }
        }

        copy.tiles = copy.tiles.map(t => { t.selected = false; return t; });
        copy.currentPlacedCells = [];

        return copy;
    }

    public selected(): Tile[] {
        let selectedTiles = this.copy().tiles.filter(t => t.selected).map(t => t.copy());
        return selectedTiles.length == 0 ? null : selectedTiles;
    }

    private setCellTile(c: Cell, t: Tile): State {
        let copy = this.copy();
        let cell = copy.cells.find(cell => cell.id == c.id);
        cell.tile = t.copy();
        copy.tiles = copy.tiles.filter(ti => ti.id != t.id);
        copy.currentPlacedCells.push(cell.copy());
        return copy;
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

    private removeCellTile(c: Cell): [State, boolean] {
        let copy = this.copy();
        let shouldReturnTile = copy.currentPlacedCells.filter(pc => pc.id == c.id).length > 0;

        if (shouldReturnTile && !copy.selected()) {
            copy.currentPlacedCells = copy.currentPlacedCells.filter(t => t.id != c.id);
            copy.tiles.push(c.tile.copy());
            c.tile = null;
            return [copy, true];
        }

        return [copy, false];
    }

    private getCells(): State {
        let copy = this.copy();
        if (copy.cells.length == 0) {
            let i = 0;
            for (let r = 0; r < 15; r++) {
                for (let c = 0; c < 15; c++) {
                    let cell = new Cell(i, null, getCellClass(i));
                    copy.cells.push(cell);
                    i++;
                }
            }
        }
        return copy;
    }

    private handleExchange(p: Payload): State {
        if (!(p.wasExchange && p.exchangeTiles.length > 0)) {
            return this;
        }

        let copy = this.copy();
        let result = copy.tiles;
        for (var rt of p.returnedTiles) {
            result = result.filter(t => rt.id != t.id);
        }
        result.push(...p.exchangeTiles);
        copy.tiles = result;
        return copy;
    }

    private handleSomeoneElsePlaced(p: Payload): State {
        if (!(p.yourMove && p.wasPlace)) {
            return this;
        }

        let copy = this.copy();
        for (var c of p.placedCells) {
            copy.cells[c.id] = c;
        }
        return copy;
    }

    private handlePlayerPlaced(p: Payload): State {
        if (!(p.wasPlace && p.exchangeTiles.length > 0)) {
            return this;
        }

        let copy = this.copy();
        for (var c of p.placedCells) {
            copy.tiles = copy.tiles.filter(t => t.id != c.tile.id);
            copy.cells[c.id] = c;
        }
        copy.tiles.push(...p.exchangeTiles);
        return copy;
    }

    private setLastAction(p: Payload): State {
        let copy = this.copy();
        if (p.wasPlace) {
            copy.lastAction = 'place';
        } else if (p.wasExchange) {
            copy.lastAction = 'exchange';
        } else if (p.isStart) {
            copy.lastAction = '';
        } else {
            copy.lastAction = 'pass'
        }
        return copy;
    }

    private setDisplayMessage(p: Payload): State {
        let copy = this.copy();
        if (p.isError) {
            copy.displayMessage = this.formatErrorMessage(p.error);
            return copy;
        }

        if (p.isStart) {
            copy.displayMessage = 'Game started.';
            return copy;
        }

        if (copy.clientLastMoved) {
            copy.displayMessage = 'Done.';
            return copy;
        }

        const lastMoved = copy.lastMovedName(p);
        switch (copy.lastAction) {
            case 'place':
                copy.displayMessage = `${lastMoved} played a word.`
                break;
            case 'exchange':
                copy.displayMessage = `${lastMoved} exchanged some tiles.`
                break;
            case 'pass':
                copy.displayMessage = `${lastMoved} passed.`
                break;
            default:
                copy.displayMessage = '';
        }
        return copy;
    }

    private lastMovedName(payload: Payload): string {
        for (var p of this.players) {
            if (p.id == payload.lastMovedId) {
                return p.name;
            }
        }
        return '';
    }

    private formatErrorMessage(error: string): string {
        console.log(error);
        if (error.indexOf('cross the origin') > -1) {
            return 'First move should cross the origin.';
        }

        if (error.indexOf('minimum field size of 1') > -1) {
            return 'Place at least two tiles.';
        }
        
        if (error.indexOf('adjacent tiles not found') > -1) {
            return 'Word must touch at least one tile.';
        }

        if (error.indexOf('invalid words: ["') > -1) {
            let result = error.match(/invalid words: \[(.*)\]/)
            return result[0].charAt(0).toUpperCase() + result[0].slice(1);
        }

        return error;
    }
}