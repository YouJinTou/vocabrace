import { Component, OnInit } from '@angular/core';
import { Cell } from './cell';

@Component({
  selector: 'app-board',
  templateUrl: './board.component.html',
  styleUrls: ['./board.component.css']
})
export class BoardComponent implements OnInit {
  cells: Cell[];

  constructor() { 
  }

  ngOnInit(): void {
    this.cells = this.getInitialBoard();
  }

  private getInitialBoard(): Cell[] {
    let board: Cell[] = [];

    for (let r = 0; r < 15; r++) {
      for (let c = 0; c < 15; c++) {
        let cell = new Cell()
        cell.value = (r * 15) + c;
        cell.letter = "a";

        board.push(cell);
      }
    }

    return board;
  }
}
