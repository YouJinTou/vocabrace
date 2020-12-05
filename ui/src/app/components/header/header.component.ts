import { Component, OnInit } from '@angular/core';
import { LoginStatusService } from 'src/services/login-status.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  loggedIn: boolean;

  constructor(private loginStatusService: LoginStatusService) { }

  ngOnInit(): void {
    this.loginStatusService.loggedIn$.subscribe(s => this.loggedIn = s);
  }

}
