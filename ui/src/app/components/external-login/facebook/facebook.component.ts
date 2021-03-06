import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { ContextService } from 'src/services/context.service';

declare var FB: any;

@Component({
  selector: 'app-facebook',
  templateUrl: './facebook.component.html',
  styleUrls: ['./facebook.component.css']
})
export class FacebookComponent implements OnInit {
  loggedIn: boolean;
  constructor(private httpClient: HttpClient, private contextService: ContextService) { }

  ngOnInit(): void {
    this.contextService.user$.subscribe(u => this.loggedIn = u.loggedIn);
    this.init();
  }

  private init() {
    (window as any).fbAsyncInit = () => {
      FB.init({
        appId: '222363819258996',
        cookie: true,
        xfbml: true,
        version: 'v3.1'
      });

      FB.AppEvents.logPageView();

      FB.getLoginStatus((r) => {
        if (r.status === 'connected') {
          console.log('connected');
          this.contextService.setLoginStatus(true);
        } else {
          console.log('not connected');
          this.contextService.setLoginStatus(false);
        }
      });
    };
  }

  login() {
    FB.login((r) => {
      if (r.authResponse) {
        FB.api('/me?fields=name,email', (res) => {
          let url = `${environment.iamEndpoint}/provider-auth`;
          this.httpClient.post(url, res).subscribe({
            next: (u: any) => {this.contextService.setUser({
              loggedIn: true,
              username: u.Username,
              id: u.ID,
              name: u.Name
            }); console.log(u);},
            error: e => console.log(e)
          });
        });
      } else {
        console.log('User cancelled login or did not fully authorize.');
      }
    }, {
      scope: "email"
    });
  }
}
