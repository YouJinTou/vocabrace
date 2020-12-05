import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export class User {
  loggedIn: boolean;
  username: string
}

@Injectable({
  providedIn: 'root'
})
export class UserStatusService {
  private userSource = new BehaviorSubject(new User());
  private current: User;
  user$ = this.userSource.asObservable();

  constructor() { 
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
