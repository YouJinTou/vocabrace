import { Injectable } from '@angular/core';
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
  userId: string
  poolId: string
}

@Injectable({
  providedIn: 'root'
})
export class ContextService {
  private userSource = new BehaviorSubject(new User());
  private wordlinesSource = new BehaviorSubject(new Wordlines());
  private isPlayingSource = new BehaviorSubject(false);
  user: User;
  wordlines: Wordlines;
  isPlaying: boolean;
  user$ = this.userSource.asObservable();
  wordlines$ = this.wordlinesSource.asObservable();
  isPlaying$ = this.isPlayingSource .asObservable();

  constructor() {
    this.user = { username: '', loggedIn: false, id: '', name: '' };
    this.user$.subscribe(u => this.user = u);
    this.wordlines = { players: 0, language: '', userId: '', poolId: ''};
    this.isPlaying = false;
    this.isPlaying$.subscribe(i => this.isPlaying = i);
    this.wordlines$.subscribe(w => this.wordlines = w);
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

  setIsPlaying(isPlaying: boolean) {
    this.isPlaying = isPlaying;
    this.isPlayingSource.next(this.isPlaying);
  }
}
