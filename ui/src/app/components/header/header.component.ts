import { Component, OnInit } from '@angular/core';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  greeting: string;

  constructor(private contextService: ContextService) { }

  ngOnInit(): void {
    this.contextService.user$.subscribe(u => this.greeting = u.loggedIn ? u.name : '');
  }
}
