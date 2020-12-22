import { UsernameService } from 'src/services/username.service'
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
    returnedTiles: Tile[]
    placedCells: Cell[]
    players = [];
    tiles: Tile[]
    blanks: Tile[]
    tilesRemaining: number;
    toMoveId: string;
    winnerId: string;
    winnerName: string;
    isGameOver: boolean;

    constructor(m: any, private usernameService: UsernameService) {
        this.isError = this.returnedError(m);
        this.yourMove = true;

        if (this.isError) {
            return;
        }

        this.yourMove = m['y'];
        this.isStart = ('s' in m) && m['s'];
        this.language = m['z'];
        this.tilesRemaining = m['r'];
        this.toMoveId = m['m'];
        this.winnerId = m['w'];
        this.isGameOver = this.winnerId != undefined;

        if (this.isStart) {
            this.tiles = this.getTiles(m['t']);
            this.players = this.getPlayersOnStart(m['p']);
            this.poolId = m['pid'];
        } else {
            this.lastAction = m['l'];
            this.lastMovedId = m['i'];
            this.wasExchange = this.lastAction === 'Exchange';
            this.wasPlace = this.lastAction === 'Place';
            this.exchangeTiles = this.getTiles(m['v']);
            this.returnedTiles = this.getTiles(m['e']);
            this.placedCells = this.getPlacedCells(m);
            this.players = this.getUpdatedPlayers(m);
        }
        this.blanks = this.getBlanks();
        this.winnerName = this.getWinnerName();
    }

    private returnedError(m: any): boolean {
        let isServerError = 'message' in m && m['message'].indexOf('Internal server error') > -1
        let isBadMove = 'Type' in m && m['Type'] == 'ERROR';
        return isServerError || isBadMove;
    }

    private getPlacedCells(m: any): Cell[] {
        if (!this.wasPlace) {
            return null;
        }
        let cells = [];
        for (var c of m['n']['c']) {
            cells.push(new Cell(c['i'], this.getTile(c['t']), getCellClass(c['i'])))
        }
        return cells;
    }

    private getTiles(m: any): Tile[] {
        let tiles = [];
        if (!m) {
            return tiles;
        }
        for (var s of m) {
            tiles.push(this.getTile(s));
        }
        return tiles;
    }

    private getPlayersOnStart(players: []): Player[] {
        let result = [];
        for (var p of players) {
            let name = this.usernameService.get(p['n']);
            result.push(new Player(p['i'], name, p['p']));
        }
        return result;
    }

    private getUpdatedPlayers(m: any): Player[] {
        if (!('p' in m && m['p'])) {
            return this.players;
        }

        let result = [];
        for (let id in m['p']) {
            let name = this.usernameService.get(id);
            let points = m['p'][id];
            result.push(new Player(id, name, points));
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
            let tile = new Tile(l, l, 0);
            result.push(tile);
        }
        return result;
    }
    
    private getWinnerName(): string {
        for (var p of this.players) {
            if (p.id == this.winnerId) {
                return p.name;
            }
        }
        return 'No winner';
    }
}