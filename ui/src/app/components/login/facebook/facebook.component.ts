import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';

declare var FB: any;

@Component({
  selector: 'app-facebook',
  templateUrl: './facebook.component.html',
  styleUrls: ['./facebook.component.css']
})
export class FacebookComponent implements OnInit {
  loggedIn: boolean;
  constructor(private httpClient: HttpClient) { }

  ngOnInit(): void {
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
          this.loggedIn = true;
        } else {
          console.log('not connected');
          this.loggedIn = false;
        }
      });
    };
  }

  login() {
    FB.login((r) => {
      if (r.authResponse) {
        FB.api('/me?fields=name,email', (r) => {
          let url = `${environment.iamEndpoint}/provider-auth`;
          this.httpClient.post(url, r).subscribe(r => this.loggedIn = true);
        });
      } else {
        console.log('User cancelled login or did not fully authorize.');
      }
    }, {
      scope: "email"
    });
  }
}
