import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-cell',
  templateUrl: './cell.component.html',
  styleUrls: ['./cell.component.css']
})
export class CellComponent implements OnInit {
  letter: string = null;
  value: number = null;
  isDoubleLetterScore: boolean = false;
  isTripleLetterScore: boolean = false;
  isDoubleWordScore: boolean = false;
  isTripleWordScore: boolean = false;

  constructor() { }

  ngOnInit(): void {
  }

}
