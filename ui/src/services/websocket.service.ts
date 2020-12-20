import { Injectable, OnDestroy } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay, retryWhen } from 'rxjs/operators';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  public connection$: WebSocketSubject<any>;
  public history = [];

  connect(url: string, params={}): Observable<any> {
    let queryString = '?';
    for (let key in params) {
      let val = params[key];
      if (val) {
        queryString += `${key}=${val}&`;
      }
    }
    queryString = queryString.slice(0, queryString.length - 1)
    let result = `${url}${queryString}`

    return of(result).pipe(_ => {
      if (!this.connection$) {
        this.connection$ = webSocket(result);
        this.connection$.subscribe(r => this.history.push(r));
      }
      return this.connection$;
    },
      retryWhen(errors => errors.pipe(delay(5)))
    );
  }

  send(data: any) {
    if (this.connection$) {
      let payload = {
        action: 'publish',
        d: JSON.stringify(data)
      }
      this.connection$.next(payload);
    } else {
      console.error('Did not send data. Open a connection first.');
    }
  }

  close() {
    if (this.connection$) {
      this.connection$.complete();
      this.connection$ = null;
    }
  }

  ngOnDestroy() {
    this.close();
  }
}
