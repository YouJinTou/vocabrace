import { Component, OnInit } from '@angular/core';
import { User, UserStatusService } from 'src/services/user-status.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  greeting: string;

  constructor(private UserStatusService: UserStatusService) { }

  ngOnInit(): void {
    this.UserStatusService.user$.subscribe(u => {
      this.greeting = u.loggedIn ? `Let's go, ${u.name}` : '';
    });
  }
}
