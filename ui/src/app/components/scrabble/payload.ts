import { load } from './alphabet'
import { Cell, getCellClass } from './cell'
import { Player } from './player'
import { Tile } from './tile'

export class Payload {
    isError: boolean
    isStart: boolean
    wasExchange: boolean
    wasPlace: boolean
    yourMove: boolean
    lastAction: string
    lastMovedId: string
    language: string
    poolId: string
    exchangeTiles: Tile[]
    placedCells: Cell[]
    players: Player[]
    tiles: Tile[]
    blanks: Tile[]

    constructor(m: any) {
        this.isError = this.returnedError(m);
        this.yourMove = true;

        if (this.isError) {
            return;
        }

        this.yourMove = m['y'];
        this.isStart = !('d' in m);
        this.language = "bulgarian";
        
        if (this.isStart) {
            this.tiles = this.getTiles(m['t']);
            this.players = this.getPlayersOnStart(m['p']);
            this.poolId = m['pid'];
        } else {
            this.lastAction = m['l'];
            this.lastMovedId = m['i'];
            this.wasExchange = this.lastAction === 'Exchange';
            this.wasPlace = this.lastAction === 'Place';
            this.exchangeTiles = this.getExchangeTiles(m);
            this.placedCells = this.getPlacedCells(m);
            this.players = this.getUpdatedPlayers(m);
        }
        this.blanks = this.getBlanks();
    }

    private returnedError(m: any): boolean {
        let isServerError = 'message' in m && m['message'].indexOf('Internal server error') > -1
        let isBadMove = 'Type' in m && m['Type'] == 'ERROR';
        return isServerError || isBadMove;
    }

    private getExchangeTiles(m: any): Tile[] {
        if (!((this.wasExchange || this.wasPlace) && Array.isArray(m['d']))) {
            return null;
        }
        return this.getTiles(m['d']);
    }

    private getPlacedCells(m: any): Cell[] {
        if (!(this.wasPlace && this.yourMove)) {
            return null;
        }
        let cells = [];
        for (var c of m['o']['Cells']) {
            cells.push(new Cell(c['i'], this.getTile(c['t']), getCellClass(c['i'])))
        }
        return cells;
    }

    private getTiles(m: any): Tile[] {
        let tiles = [];
        for (var s of m) {
            tiles.push(this.getTile(s));
        }
        return tiles;
    }

    private getPlayersOnStart(players: []): Player[] {
        let result = [];
        for (var p of players) {
            result.push(new Player(p['n'], p['p']));
        }
        return result;
    }

    private getUpdatedPlayers(m: any): Player[] {
        if (!('p' in m && m['p'])) {
            return this.players;
        }

        let result = [];
        for (let key in m['p']) {
            let value = m['p'][key];
            result.push(new Player(key, value));
        }
        return result;
    }

    private getTile(s: string): Tile {
        let tokens = s.split("|");
        let tile = new Tile(tokens[0], tokens[1], parseInt(tokens[2]));
        return tile;
    }

    private getBlanks(): Tile[] {
        let letters = load(this.language)
        let result = [];
        for (var l of letters) {
            let tile = new Tile("to_be_replaced", l, 0);
            result.push(tile);
        }
        return result;
    }
}