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
    exchangeTiles: Tile[]
    placedCells: Cell[]
    players: Player[]
    tiles: Tile[]

    constructor(m: any) {
        console.log(m);
        this.isError = this.returnedError(m);
        if (this.isError) {
            return;
        }
        this.yourMove = m['y'];
        this.isStart = !('d' in m);
        if (this.isStart) {
            this.tiles = this.getTiles(m['t']);
            this.players = this.getPlayers(m['p']);
        } else {
            this.lastAction = m['l'];
            this.lastMovedId = m['i'];
            this.wasExchange = this.lastAction === 'Exchange';
            this.wasPlace = this.lastAction === 'Place';
            this.exchangeTiles = this.getExchangeTiles(m);
            this.placedCells = this.getPlacedCells(m);
        }
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

    private getPlayers(players: []): Player[] {
        let result = [];
        for (var p of players) {
            result.push(new Player(p['n'], p['p']));
        }
        return players;
    }

    private getTile(s: string): Tile {
        let tokens = s.split("|");
        let tile = new Tile(tokens[0], tokens[1], parseInt(tokens[2]));
        return tile;
    }
}