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
export class UserStatusService {
  private userSource = new BehaviorSubject(new User());
  private current: User;
  user$ = this.userSource.asObservable();

  constructor() {
    this.current = { username: '', loggedIn: false, id: '', name: '' };
    this.user$.subscribe(u => this.current = u);
  }

  setUser(user: User) {
    this.userSource.next(user);
  }

  setLoginStatus(loggedIn: boolean) {
    this.current.loggedIn = loggedIn;
    this.userSource.next(this.current);
  }
}
