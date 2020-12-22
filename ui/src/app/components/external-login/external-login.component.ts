import { Component, OnInit } from '@angular/core';
import { ContextService, Status } from 'src/services/context.service';

@Component({
  selector: 'app-external-login',
  templateUrl: './external-login.component.html',
  styleUrls: ['./external-login.component.css']
})
export class ExternalLoginComponent implements OnInit {
  status: Status;

  constructor(private contextService: ContextService) { }

  ngOnInit(): void {
    this.contextService.status$.subscribe(i => {
      this.status = i;});
  }
}
