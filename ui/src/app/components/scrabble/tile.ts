export class Tile {
    id: string
    letter: string
    value: number
  
    constructor(id: string, letter: string, value: number) {
      this.id = id;
      this.letter = letter;
      this.value = value;
    }
  
    copy(): Tile {
      return new Tile(this.id, this.letter, this.value);
    }
  }