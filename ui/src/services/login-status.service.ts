import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LoginStatusService {
  private loggedInSource = new BehaviorSubject(false);
  loggedIn$ = this.loggedInSource.asObservable();

  constructor() { }

  setStatus(loggedIn: boolean) {
    this.loggedInSource.next(loggedIn);
  }
}
