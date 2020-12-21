import { Injectable } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';
import { BehaviorSubject } from 'rxjs';

export class User {
  loggedIn: boolean
  username: string
  name: string
  id: string
}

export class Wordlines {
  players: number
  language: string
  poolId: string
}

export class IsPlaying {
  value: boolean
  pid: string
}

export class Cookies {
  accepted: boolean
}

@Injectable({
  providedIn: 'root'
})
export class ContextService {
  private userSource = new BehaviorSubject(new User());
  private wordlinesSource = new BehaviorSubject(new Wordlines());
  private isPlayingSource = new BehaviorSubject(new IsPlaying());
  private cookiesSource = new BehaviorSubject(new Cookies());
  user: User;
  wordlines: Wordlines;
  isPlaying: IsPlaying;
  cookies: Cookies;
  user$ = this.userSource.asObservable();
  wordlines$ = this.wordlinesSource.asObservable();
  isPlaying$ = this.isPlayingSource.asObservable();
  cookies$ = this.cookiesSource.asObservable();

  constructor(private cookieService: CookieService) {
    this.user = { username: '', loggedIn: false, id: '', name: '' };
    this.user$.subscribe(u => this.user = u);
    this.wordlines = { players: 0, language: '', poolId: '' };
    this.isPlaying = { value: false, pid: '' };
    this.isPlaying$.subscribe(i => this.isPlaying = i);
    this.wordlines$.subscribe(w => this.wordlines = w);
    this.cookies$.subscribe(c => this.cookies = c);
  }

  setUser(user: User) {
    this.userSource.next(user);
  }

  setWordlines(wordlines: Wordlines) {
    this.wordlinesSource.next(wordlines);
  }

  setLoginStatus(loggedIn: boolean) {
    this.user.loggedIn = loggedIn;
    this.userSource.next(this.user);
  }

  setIsPlaying(isPlaying: IsPlaying) {
    this.isPlayingSource.next(isPlaying);
    if (this.cookies.accepted) {
      const expires = new Date(new Date().getTime() + (1000 * 60 * 5));
      this.cookieService.set('pid', isPlaying.pid, { expires: expires, sameSite: 'Strict' });
    }
  }

  setCookies(cookies: Cookies) {
    this.cookiesSource.next(cookies);
  }
}
