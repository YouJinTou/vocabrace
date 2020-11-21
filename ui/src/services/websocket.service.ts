import { Injectable, OnDestroy } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay, retryWhen } from 'rxjs/operators';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  private connection$: WebSocketSubject<any>;

  connect(url: string, game: string): Observable<any> {
    let result = `${url}?game=${game}`
    return of(result).pipe(_ => {
      if (!this.connection$) {
        this.connection$ = webSocket(result);
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
