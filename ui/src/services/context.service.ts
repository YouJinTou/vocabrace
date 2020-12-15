import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export class User {
  loggedIn: boolean
  username: string
  name: string
  id: string
}

@Injectable({
  providedIn: 'root'
})
export class ContextService {
  private userSource = new BehaviorSubject(new User());
  private isPlayingSource = new BehaviorSubject(false);
  user: User;
  user$ = this.userSource.asObservable();
  isPlaying: boolean;
  isPlaying$ = this.isPlayingSource .asObservable();

  constructor() {
    this.user = { username: '', loggedIn: false, id: '', name: '' };
    this.user$.subscribe(u => this.user = u);
    this.isPlaying = false;
    this.isPlaying$.subscribe(i => this.isPlaying = i);
  }

  setUser(user: User) {
    this.userSource.next(user);
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
