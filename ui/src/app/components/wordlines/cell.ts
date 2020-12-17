import { Tile } from './tile';
import { Spiral } from './layout';

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
  let cls = Spiral.DOUBLE_LETTER_INDICES.includes(i) ? 'double-letter' :
    Spiral.TRIPLE_LETTER_INDICES.includes(i) ? 'triple-letter' :
      Spiral.DOUBLE_WORD_INDICES.includes(i) ? 'double-word' :
        Spiral.TRIPLE_WORD_INDICES.includes(i) ? 'triple-word' :
          'tile'
  return cls;
}