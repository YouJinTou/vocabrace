import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.css']
})
export class NotificationComponent implements OnInit {
  message: string;
  showClass = "out";
  backgroundClass = "info";

  constructor() { }

  ngOnInit(): void {
  }

  showError(message: string, millis?: number) {
    this.backgroundClass = "error";
    this.show(message, millis);
  }

  showSuccess(message: string, millis?: number) {
    this.backgroundClass = "success";
    this.show(message, millis);
  }

  showInfo(message: string, millis?: number) {
    this.backgroundClass = "info";
    this.show(message, millis);
  }

  private show(message: string, millis?: number) {
    this.message = message;
    this.showClass = "in";
    setTimeout(() => {
      this.showClass = "out";
    }, millis ?? 4500);
  }
}
