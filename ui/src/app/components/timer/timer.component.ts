import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-timer',
  templateUrl: './timer.component.html',
  styleUrls: ['./timer.component.css']
})
export class TimerComponent implements OnInit {
  private intervalId;
  showWarning = false;
  remaining: number;
  @Input() timeout: number;
  @Input() warning = 10;
  @Input() warningClass = 'warning'
  @Input() warningEnabled = true;
  @Output() timedOut = new EventEmitter();

  constructor() { }

  ngOnInit(): void {
  }

  start() {
    this.countdown();
  }

  reset() {
    clearTimeout(this.intervalId);
  }

  restart() {
    this.reset();
    this.start();
  }

  private countdown() {
    this.remaining = this.timeout;
    this.intervalId = setInterval(() => {
      this.remaining -= 1;
      if (this.remaining <= 0) {
        this.timedOut.emit();
        this.reset();
      } else if (this.warningEnabled && this.remaining <= this.warning) {
        this.showWarning = true;
      }
    }, 1000);
  }
}
