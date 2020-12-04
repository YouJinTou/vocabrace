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

  constructor(private httpClient: HttpClient) { }

  ngOnInit(): void {
    this.init();
  }

  private init() {
    (window as any).fbAsyncInit = function () {
      FB.init({
        appId: '222363819258996',
        cookie: true,
        xfbml: true,
        version: 'v3.1'
      });
      FB.AppEvents.logPageView();
    };
  }

  login() {
    let $httpClient = this.httpClient
    FB.login(function (response) {
      if (response.authResponse) {
        console.log('Welcome!  Fetching your information.... ');
        console.log(response);
        FB.api('/me?fields=name,email', function (response) {
          let url = `${environment.iamEndpoint}/provider-auth`;
          console.log(url);
          $httpClient.post(url, response).subscribe(r => console.log(r));
        });
      } else {
        console.log('User cancelled login or did not fully authorize.');
      }
    }, {
      scope: "email"
    });
  }

}
