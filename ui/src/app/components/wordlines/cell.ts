import { Tile } from './tile';

export const DOUBLE_LETTER_INDICES = [3, 11, 36, 38, 45, 52, 59, 92, 96, 98, 102, 108, 116, 122, 126, 128, 132, 165, 172, 179, 186, 188, 213, 221];
export const DOUBLE_WORD_INDICES = [16, 32, 48, 64, 112, 160, 176, 192, 208, 28, 42, 56, 70, 154, 168, 182, 196];
export const TRIPLE_LETTER_INDICES = [20, 24, 76, 80, 84, 88, 136, 140, 144, 148, 200, 204];
export const TRIPLE_WORD_INDICES = [0, 7, 14, 105, 119, 210, 217, 224];

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

export function getCellClass(i: number): string {
  let cls = DOUBLE_LETTER_INDICES.includes(i) ? 'double-letter' :
    TRIPLE_LETTER_INDICES.includes(i) ? 'triple-letter' :
      DOUBLE_WORD_INDICES.includes(i) ? 'double-word' :
        TRIPLE_WORD_INDICES.includes(i) ? 'triple-word' :
          'tile'
  return cls;
}