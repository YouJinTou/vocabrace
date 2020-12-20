export class Tile {
    id: string
    letter: string
    value: number
    selected: boolean
  
    constructor(id: string, letter: string, value: number) {
      this.id = id;
      this.letter = letter;
      this.value = value;
      this.selected = false;
    }

    isBlank(): boolean {
      return this.value == 0;
    }
  
    copy(): Tile {
      let t = new Tile(this.id, this.letter, this.value);
      t.selected = this.selected;
      return t;
    }
  }