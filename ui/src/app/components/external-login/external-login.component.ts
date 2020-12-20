import { Component, OnInit } from '@angular/core';
import { ContextService, IsPlaying } from 'src/services/context.service';

@Component({
  selector: 'app-external-login',
  templateUrl: './external-login.component.html',
  styleUrls: ['./external-login.component.css']
})
export class ExternalLoginComponent implements OnInit {
  isPlaying: IsPlaying;

  constructor(private contextService: ContextService) { }

  ngOnInit(): void {
    this.contextService.isPlaying$.subscribe(i => {
      this.isPlaying = i;});
  }
}
