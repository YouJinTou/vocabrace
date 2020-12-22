import { Injectable } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';
import { BehaviorSubject } from 'rxjs';

export class User {
  loggedIn: boolean
  username: string
  name: string
  id: string
}

export class Status {
  game: string
  value: boolean
  pid: string
  players: number
  language: string
}

export class Cookies {
  accepted: boolean
}

@Injectable({
  providedIn: 'root'
})
export class ContextService {
  private userSource = new BehaviorSubject(new User());
  private statusSource = new BehaviorSubject(new Status());
  private cookiesSource = new BehaviorSubject(new Cookies());
  user: User;
  status: Status;
  cookies: Cookies;
  user$ = this.userSource.asObservable();
  status$ = this.statusSource.asObservable();
  cookies$ = this.cookiesSource.asObservable();

  constructor(private cookieService: CookieService) {
    this.user = { username: '', loggedIn: false, id: '', name: '' };
    this.user$.subscribe(u => this.user = u);
    this.status = { game: '', value: false, pid: '', language: '', players: 0 };
    this.status$.subscribe(i => this.status = i);
    this.cookies$.subscribe(c => this.cookies = c);
  }

  setUser(user: User) {
    this.userSource.next(user);
  }

  setLoginStatus(loggedIn: boolean) {
    this.user.loggedIn = loggedIn;
    this.userSource.next(this.user);
  }

  setStatus(status: Status) {
    this.statusSource.next(status);
    if (this.cookies.accepted) {
      const expires = new Date(new Date().getTime() + (1000 * 60 * 5));
      this.cookieService.set('pid', status.pid, { expires: expires, sameSite: 'Strict' });
    }
  }

  setCookies(cookies: Cookies) {
    this.cookiesSource.next(cookies);
  }
}
