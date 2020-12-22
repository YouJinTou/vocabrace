import { Component, OnInit } from '@angular/core';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  isPlaying = false;

  constructor(private contextService: ContextService) {}
  
  ngOnInit() {
    this.contextService.status$.subscribe(s => this.isPlaying = s.value);
  }
}
