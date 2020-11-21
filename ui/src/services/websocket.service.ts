import { Injectable, OnDestroy } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay, retryWhen } from 'rxjs/operators';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService implements OnDestroy {
  private connection$: WebSocketSubject<any>;

  connect(url: string): Observable<any> {
    return of(url).pipe(_ => {
      if (!this.connection$) {
        this.connection$ = webSocket(url);
      }
      return this.connection$;
    },
      retryWhen(errors => errors.pipe(delay(5)))
    );
  }

  send(data: any) {
    if (this.connection$) {
      this.connection$.next(data);
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
