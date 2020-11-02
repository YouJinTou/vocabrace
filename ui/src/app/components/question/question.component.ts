import { Component, OnInit } from '@angular/core';
import { Question } from './question';

@Component({
  selector: 'app-question',
  templateUrl: './question.component.html',
  styleUrls: ['./question.component.css']
})
export class QuestionComponent implements OnInit {
  question: Question;

  constructor() {
    this.question = {
      text: "red or black?", answers: [
        "Where is it?",
        "There it is"
      ], correctAnswerIdx: 0
    };
  }

  ngOnInit(): void {
  }

}
