import { Component, OnInit } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';
import { ContextService } from 'src/services/context.service';

@Component({
  selector: 'app-cookies',
  templateUrl: './cookies.component.html',
  styleUrls: ['./cookies.component.css']
})
export class CookiesComponent implements OnInit {
  accepted = false;

  constructor(private service: CookieService, private contextService: ContextService) { }

  ngOnInit(): void {
    const acceptedCookie = this.service.get('accepted');
    this.accepted = acceptedCookie && acceptedCookie === 'true';
    this.contextService.setCookies({ accepted: this.accepted });
  }

  acceptAndClose() {
    this.accepted = true;
    const expires = new Date(new Date().getTime() + (1000 * 5));
    this.service.set('accepted', 'true', { expires: expires, sameSite: 'Strict' });
    this.contextService.setCookies({ accepted: this.accepted });
  }
}
