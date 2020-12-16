import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { Payload } from '../payload';

@Component({
  selector: 'scrabble-game-over',
  templateUrl: 'game-over.component.html',
  styleUrls: ['game-over.component.css']
})
export class GameOverDialog {
  winner: string;

  constructor(
    public dialogRef: MatDialogRef<GameOverDialog>, 
    @Inject(MAT_DIALOG_DATA) public payload: Payload,
    private router: Router
    ) {
      this.winner = payload.winnerName;
  }

  onHomeClicked() {
    this.dialogRef.close();
    this.router.navigate(['/']);
  }
}