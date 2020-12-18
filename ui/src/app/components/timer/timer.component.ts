import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-timer',
  templateUrl: './timer.component.html',
  styleUrls: ['./timer.component.css']
})
export class TimerComponent implements OnInit {
  private intervalId;
  remaining: number;
  @Input() timeout: number;
  @Output() onExpire = new EventEmitter<any>();

  constructor() { }

  ngOnInit(): void {
  }

  start() {
    this.countdown();
  }

  reset() {
    clearTimeout(this.intervalId);
  }

  private countdown() {
    this.remaining = this.timeout + 1;
    this.intervalId = setInterval(() => {
      this.remaining -= 1;
    }, 1000);
  }
}
