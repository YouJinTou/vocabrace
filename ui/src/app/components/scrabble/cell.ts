import { Tile } from './tile';

export class Cell {
    id: number
    tile: Tile
    cssClass: string
  
    constructor(id: number, tile: Tile, cssClass: string) {
      this.id = id;
      this.tile = tile;
      this.cssClass = cssClass;
    }
  
    isEmpty(): boolean {
      return this.tile == null;
    }

    copy(): Cell {
      return new Cell(this.id, this.tile, this.cssClass);
    }
  }