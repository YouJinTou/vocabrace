import { PlatformRef } from '@angular/core'
import { Player } from './player'
import { Tile } from './tile'

export class Payload {
    isError: boolean
    isStart: boolean
    wasExchange: boolean
    yourMove: boolean
    lastAction: string
    lastMovedId: string
    exchangeTiles: Tile[]
    players: Player[]
    tiles: Tile[]

    constructor(m: any) {
        console.log(m);
        this.isError = 'message' in m && m['message'].indexOf('Internal server error') > -1;
        this.yourMove = m['y'];
        this.isStart = !('d' in m);
        if (this.isStart) {
            this.tiles = this.getTiles(m['t']);
            this.players = this.getPlayers(m['p']);
        } else {
            this.lastAction = m['l'];
            this.lastMovedId = m['i'];
            this.wasExchange = this.lastAction === 'Exchange';
            this.exchangeTiles = this.getExchangeTiles(m);
        }
    }

    private getExchangeTiles(m: any): Tile[] {
        if (!(this.wasExchange && Array.isArray(m['d']))) {
            return null;
        }
        return this.getTiles(m['d']);
    }

    private getTiles(m: any): Tile[] {
        let tiles = [];
        for (var s of m) {
            let tokens = s.split("|");
            tiles.push(new Tile(tokens[0], tokens[1], tokens[2]));
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
}