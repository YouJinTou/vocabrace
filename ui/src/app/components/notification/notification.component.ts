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

  showError(message: string) {
    this.backgroundClass = "error";
    this.show(message);
  }

  showSuccess(message: string) {
    this.backgroundClass = "success";
    this.show(message);
  }

  showInfo(message: string) {
    this.backgroundClass = "info";
    this.show(message);
  }

  private show(message: string) {
    this.message = message;
    this.showClass = "in";
    setTimeout(() => {
      this.showClass = "out";
    }, 3000);
  }
}
