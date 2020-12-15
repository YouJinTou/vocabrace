import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class UsernameService {
  private adjectives = ['smart', 'quick', 'beautiful', 'sly', 'cunning', 'dapper', 'mighty', 'strong'];
  private nouns = ['ninja', 'plumber', 'fireman', 'accountant', 'pilot', 'hunter', 'assassin'];
  private store = {};

  constructor() { }

  generate(): string {
    let adjectiveIdx = Math.floor(Math.random() * Math.floor(this.adjectives.length));
    let nounIdx = Math.floor(Math.random() * Math.floor(this.nouns.length));
    let result = `${this.adjectives[adjectiveIdx]}_${this.nouns[nounIdx]}`;
    return result;
  }

  get(id: string): string {
    if (id in this.store) {
      return this.store[id];
    }
    this.store[id] = this.generate();
    return this.store[id];
  }
}
