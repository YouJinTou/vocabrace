import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Tile } from '../tile';

@Component({
  selector: 'wordlines-blanks',
  templateUrl: 'blanks.component.html',
  styleUrls: ['../wordlines.component.css']
})
export class BlanksDialog {
  blanks: Tile[];

  constructor(
    public dialogRef: MatDialogRef<BlanksDialog>,
    @Inject(MAT_DIALOG_DATA) public data: { blanks: Tile[] }) {
    this.blanks = data['blanks'];
  }

  onSelected(t: Tile) {
    this.dialogRef.close(t);
  }
}
