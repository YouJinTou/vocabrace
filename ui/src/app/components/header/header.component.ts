import { Component, OnInit } from '@angular/core';
import { userStatusService } from 'src/services/user-status.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  greeting: string;

  constructor(private userStatusService: UserStatusService) { }

  ngOnInit(): void {
    this.userStatusService.user$.subscribe(
      u => this.greeting = u.loggedIn ? `Let's go, ${u.name}` : '');
  }
}
